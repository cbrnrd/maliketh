package admin

import (
	"encoding/json"
	"net/http"

	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/cbrnrd/maliketh/server-go/internal/server"
	"go.uber.org/zap"
)

type revokeOperatorRequest struct {
	Username string `json:"username"`
}

// handleAdminRevokeOperator allows admins to revoke another operator's access.
func (h *Handler) handleAdminRevokeOperator() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		// Check if operator is admin
		if operator.Role != model.RoleAdmin {
			h.logger.Warn("non-admin tried to revoke operator",
				zap.String("operator", operator.Username))
			server.Error(w, http.StatusForbidden, "Admin access required")
			return
		}

		var req revokeOperatorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			server.BadRequest(w, "Invalid request body")
			return
		}

		if req.Username == "" {
			server.BadRequest(w, "username is required")
			return
		}

		// Cannot revoke yourself
		if req.Username == operator.Username {
			server.BadRequest(w, "Cannot revoke your own access")
			return
		}

		// Check if target operator exists
		target, err := h.repos.Operator.GetByUsername(r.Context(), req.Username)
		if err != nil {
			server.BadRequest(w, "Unknown operator")
			return
		}

		// Revoke the operator
		if err := h.repos.Operator.Revoke(r.Context(), target.Username); err != nil {
			h.logger.Error("failed to revoke operator",
				zap.String("target", target.Username),
				zap.String("admin", operator.Username),
				zap.Error(err))
			server.InternalError(w, "Failed to revoke operator")
			return
		}

		h.logger.Info("revoked operator access",
			zap.String("target", target.Username),
			zap.String("admin", operator.Username))

		server.Success(w, "Operator access revoked")
	}
}
