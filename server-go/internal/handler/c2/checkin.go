package c2

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cbrnrd/maliketh/server-go/internal/crypto"
	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (h *Handler) handleCheckin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		implant := GetImplant(r.Context())
		if implant == nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Update last seen
		if err := h.repos.Implant.UpdateLastSeen(r.Context(), implant.ImplantID); err != nil {
			h.logger.Error("failed to update last seen",
				zap.String("implant_id", implant.ImplantID),
				zap.Error(err))
		}

		h.logger.Info("implant checked in", zap.String("implant_id", implant.ImplantID))

		// Get oldest pending task
		task, err := h.repos.Task.GetOldestPendingForImplant(r.Context(), implant.ImplantID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			h.logger.Error("failed to get pending task",
				zap.String("implant_id", implant.ImplantID),
				zap.Error(err))
		}

		var responseBody []byte

		if task != nil {
			// Mark task as read/tasked
			if err := h.repos.Task.MarkAsTasked(r.Context(), task.TaskID); err != nil {
				h.logger.Error("failed to mark task as tasked",
					zap.String("task_id", task.TaskID),
					zap.Error(err))
			}

			// Handle self-destruct
			if task.Opcode == model.OpSelfDestruct {
				h.handleSelfDestruct(r, implant)
			}

			// Return task to implant
			responseBody, _ = json.Marshal(task.ToFilteredJSON())
		} else {
			// No pending tasks - return empty object
			responseBody = []byte("{}")
		}

		// Encrypt response
		encrypted, err := crypto.EncryptB64(implant.ImplantPK, implant.ServerSK, responseBody)
		if err != nil {
			h.logger.Error("failed to encrypt response",
				zap.String("implant_id", implant.ImplantID),
				zap.Error(err))
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte(encrypted))
	}
}

func (h *Handler) handleSelfDestruct(r *http.Request, implant *model.Implant) {
	if err := h.repos.Implant.Delete(r.Context(), implant.ImplantID); err != nil {
		h.logger.Error("failed to delete implant on self-destruct",
			zap.String("implant_id", implant.ImplantID),
			zap.Error(err))
		return
	}

	h.logger.Warn("implant self-destructed", zap.String("implant_id", implant.ImplantID))

	if h.publisher != nil {
		h.publisher.PublishAnnouncement(r.Context(), "Implant "+implant.ImplantID+" self-destructed")
	}
}
