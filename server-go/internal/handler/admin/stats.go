package admin

import (
	"net/http"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/server"
	"go.uber.org/zap"
)

// handleStats returns basic server statistics.
func (h *Handler) handleStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		ctx := r.Context()

		// Get counts
		numImplants, err := h.repos.Implant.Count(ctx)
		if err != nil {
			h.logger.Error("failed to count implants", zap.Error(err))
			numImplants = 0
		}

		numActiveTasks, err := h.repos.Task.CountActive(ctx, operator.Username)
		if err != nil {
			h.logger.Error("failed to count active tasks", zap.Error(err))
			numActiveTasks = 0
		}

		numTotalTasks, err := h.repos.Task.CountTotal(ctx, operator.Username)
		if err != nil {
			h.logger.Error("failed to count total tasks", zap.Error(err))
			numTotalTasks = 0
		}

		numOperators, err := h.repos.Operator.Count(ctx)
		if err != nil {
			h.logger.Error("failed to count operators", zap.Error(err))
			numOperators = 0
		}

		uptime := time.Since(h.startTime).String()

		server.SuccessData(w, map[string]any{
			"status":       true,
			"implants":     numImplants,
			"active_tasks": numActiveTasks,
			"total_tasks":  numTotalTasks,
			"operators":    numOperators,
			"uptime":       uptime,
		})
	}
}
