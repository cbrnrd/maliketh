package admin

import (
	"net/http"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/crypto"
	"github.com/cbrnrd/maliketh/server-go/internal/server"
	"go.uber.org/zap"
)

const tokenDuration = 6 * time.Hour

// handleRequestToken handles operator authentication via signed challenge.
// X-ID: operator username
// X-Signature: base64(encrypt(sign(login_secret, operator_signing_key), server_pub_key))
func (h *Handler) handleRequestToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check required headers
		username := r.Header.Get("X-ID")
		signature := r.Header.Get("X-Signature")

		if username == "" || signature == "" {
			server.BadRequest(w, "Unknown operator key")
			return
		}

		h.logger.Info("received auth token request", zap.String("operator", username))

		// Get operator from database
		operator, err := h.repos.Operator.GetByUsername(r.Context(), username)
		if err != nil {
			h.logger.Debug("operator not found", zap.String("username", username), zap.Error(err))
			server.BadRequest(w, "Unknown operator")
			return
		}

		// Decrypt and verify the signature
		// The operator signs their login_secret, then encrypts with server's public key
		// We decrypt with server's private key, then verify with operator's verify key
		originalMessage, err := crypto.DecryptAndVerify(
			h.serverPrivateKey,
			operator.PublicKey,
			operator.VerifyKey,
			signature,
		)
		if err != nil {
			h.logger.Error("failed to decrypt/verify signature",
				zap.String("operator", username),
				zap.Error(err))
			server.BadRequest(w, "Couldn't verify signature")
			return
		}

		// Verify the message matches the login secret
		if string(originalMessage) != operator.LoginSecret {
			h.logger.Error("invalid signature - login secret mismatch",
				zap.String("operator", username))
			server.BadRequest(w, "Unable to verify signature")
			return
		}

		// Check if existing token is still valid
		var token string
		if operator.AuthToken != nil && operator.AuthTokenExpiry != nil &&
			operator.AuthTokenExpiry.After(time.Now()) {
			token = *operator.AuthToken
		} else {
			// Generate new token
			var err error
			token, err = h.repos.Operator.GenerateAuthToken(r.Context(), username, tokenDuration)
			if err != nil {
				h.logger.Error("failed to generate auth token",
					zap.String("operator", username),
					zap.Error(err))
				server.InternalError(w, "Failed to generate token")
				return
			}
		}

		// Update last login
		if err := h.repos.Operator.UpdateLastLogin(r.Context(), username); err != nil {
			h.logger.Warn("failed to update last login",
				zap.String("operator", username),
				zap.Error(err))
		}

		h.logger.Info("operator authenticated successfully", zap.String("operator", username))

		server.SuccessData(w, map[string]any{
			"status":   true,
			"token":    token,
			"rmq_queue": operator.RMQQueue,
			"rmq_host": h.cfg.RabbitMQ.Host,
			"rmq_port": h.cfg.RabbitMQ.Port,
		})
	}
}

// handleRevokeToken revokes the current operator's auth token.
func (h *Handler) handleRevokeToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operator := GetOperator(r.Context())
		if operator == nil {
			server.Unauthorized(w, "Not authenticated")
			return
		}

		h.logger.Info("revoking auth token", zap.String("operator", operator.Username))

		if err := h.repos.Operator.ClearAuthToken(r.Context(), operator.Username); err != nil {
			h.logger.Error("failed to revoke token",
				zap.String("operator", operator.Username),
				zap.Error(err))
			server.InternalError(w, "Failed to revoke token")
			return
		}

		server.SuccessData(w, map[string]any{"status": true})
	}
}

// handleTokenStatus checks if the current token is valid.
func (h *Handler) handleTokenStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 {
			server.Error(w, http.StatusOK, "Not authenticated")
			return
		}

		token := authHeader[7:] // Skip "Bearer "

		operator, err := h.repos.Operator.GetByAuthToken(r.Context(), token)
		if err != nil {
			server.Error(w, http.StatusOK, "Not authenticated")
			return
		}

		if operator.AuthTokenExpiry == nil || operator.AuthTokenExpiry.Before(time.Now()) {
			server.Error(w, http.StatusOK, "Not authenticated")
			return
		}

		server.Success(w, "Authenticated")
	}
}
