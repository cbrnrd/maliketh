package admin

import (
	"encoding/json"
	"net/http"

	"github.com/cbrnrd/maliketh/server-go/internal/server"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type createAliasRequest struct {
	Alias string `json:"alias"`
}

// handleCreateAlias creates an alias for an implant.
func (h *Handler) handleCreateAlias() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		implantID := chi.URLParam(r, "implant_id")

		var req createAliasRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			server.BadRequest(w, "Invalid request body")
			return
		}

		if req.Alias == "" {
			server.BadRequest(w, "Invalid request, no alias field")
			return
		}

		// Check if implant exists
		implant, err := h.repos.Implant.GetByID(r.Context(), implantID)
		if err != nil {
			server.BadRequest(w, "Unknown implant")
			return
		}

		// Check if alias already exists
		existing, err := h.repos.Implant.GetAliasByName(r.Context(), req.Alias)
		if err == nil && existing != nil {
			server.BadRequest(w, "Alias already exists")
			return
		}

		// Create alias
		_, err = h.repos.Implant.CreateAlias(r.Context(), implant.ImplantID, req.Alias)
		if err != nil {
			h.logger.Error("failed to create alias",
				zap.String("implant_id", implant.ImplantID),
				zap.String("alias", req.Alias),
				zap.Error(err))
			server.InternalError(w, "Failed to create alias")
			return
		}

		h.logger.Info("created implant alias",
			zap.String("implant_id", implant.ImplantID),
			zap.String("alias", req.Alias),
			zap.String("operator", operator.Username))

		server.Success(w, "Alias created")
	}
}

// handleDeleteAlias deletes an alias for an implant.
func (h *Handler) handleDeleteAlias() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		implantID := chi.URLParam(r, "implant_id")
		alias := chi.URLParam(r, "alias")

		// Verify alias exists and belongs to this implant
		existing, err := h.repos.Implant.GetAliasByName(r.Context(), alias)
		if err != nil {
			server.BadRequest(w, "Alias does not exist")
			return
		}

		if existing.ImplantID != implantID {
			server.BadRequest(w, "Alias does not belong to this implant")
			return
		}

		if err := h.repos.Implant.DeleteAlias(r.Context(), alias); err != nil {
			h.logger.Error("failed to delete alias",
				zap.String("implant_id", implantID),
				zap.String("alias", alias),
				zap.Error(err))
			server.InternalError(w, "Failed to delete alias")
			return
		}

		h.logger.Info("deleted implant alias",
			zap.String("implant_id", implantID),
			zap.String("alias", alias),
			zap.String("operator", operator.Username))

		server.Success(w, "Alias deleted")
	}
}

// handleListAliases returns all aliases for an implant.
func (h *Handler) handleListAliases() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		implantID := chi.URLParam(r, "implant_id")

		aliases, err := h.repos.Implant.GetAliasesByImplantID(r.Context(), implantID)
		if err != nil {
			h.logger.Error("failed to list aliases",
				zap.String("implant_id", implantID),
				zap.Error(err))
			server.InternalError(w, "Failed to list aliases")
			return
		}

		aliasNames := make([]string, len(aliases))
		for i, a := range aliases {
			aliasNames[i] = a.Alias
		}

		server.SuccessData(w, map[string]any{
			"status":  true,
			"aliases": aliasNames,
		})
	}
}

// handleResolveAlias resolves an alias to an implant ID.
func (h *Handler) handleResolveAlias() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		alias := chi.URLParam(r, "alias")

		a, err := h.repos.Implant.GetAliasByName(r.Context(), alias)
		if err != nil {
			server.NotFound(w, "Unknown alias")
			return
		}

		server.SuccessData(w, map[string]any{
			"status":     true,
			"implant_id": a.ImplantID,
		})
	}
}
