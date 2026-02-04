package admin

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/cbrnrd/maliketh/server-go/internal/server"
	"go.uber.org/zap"
)

type contextKey string

const operatorContextKey contextKey = "operator"

// GetOperator extracts the operator from the request context.
func GetOperator(ctx context.Context) *model.Operator {
	if v := ctx.Value(operatorContextKey); v != nil {
		return v.(*model.Operator)
	}
	return nil
}

// operatorAuthMiddleware verifies the Bearer token and checks operator status.
func (h *Handler) operatorAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		// Check Bearer token format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		// Look up operator by token
		operator, err := h.repos.Operator.GetByAuthToken(r.Context(), token)
		if err != nil {
			h.logger.Debug("operator not found for token", zap.Error(err))
			server.Unauthorized(w, "Not authenticated")
			return
		}

		// Check if token is expired
		if operator.AuthTokenExpiry == nil || operator.AuthTokenExpiry.Before(time.Now()) {
			h.logger.Debug("auth token expired",
				zap.String("operator", operator.Username),
				zap.Timep("expiry", operator.AuthTokenExpiry))
			server.Unauthorized(w, "Not authenticated")
			return
		}

		// Check if operator is revoked
		if operator.Revoked {
			h.logger.Warn("revoked operator attempted access",
				zap.String("operator", operator.Username))
			server.Error(w, http.StatusForbidden, "Your access has been revoked. Please contact the admin to resolve.")
			return
		}

		// Add operator to context
		ctx := context.WithValue(r.Context(), operatorContextKey, operator)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
