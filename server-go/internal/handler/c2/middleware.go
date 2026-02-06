package c2

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/cbrnrd/maliketh/server-go/internal/crypto"
	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"go.uber.org/zap"
)

type contextKey string

const (
	implantContextKey       contextKey = "implant"
	decryptedBodyContextKey contextKey = "decrypted_body"
)

func GetImplant(ctx context.Context) *model.Implant {
	if v := ctx.Value(implantContextKey); v != nil {
		return v.(*model.Implant)
	}
	return nil
}

func GetDecryptedBody(ctx context.Context) map[string]any {
	if v := ctx.Value(decryptedBodyContextKey); v != nil {
		return v.(map[string]any)
	}
	return nil
}

func (h *Handler) implantAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(h.profile.Globals.ImplantIDCookie)
		if err != nil || cookie.Value == "" {
			h.logger.Debug("missing implant cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		implantID := cookie.Value

		implant, err := h.repos.Implant.GetByID(r.Context(), implantID)
		if err != nil {
			h.logger.Debug("implant not found", zap.String("implant_id", implantID))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), implantContextKey, implant)

		// If there's a request body, decode and decrypt it
		if r.ContentLength > 0 {
			// Decrypt the body using implant's public key and server's private key
			decodedBody, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, r.Body))
			decrypted, err := crypto.DecryptB64Raw(implant.ImplantPK, implant.ServerSK, decodedBody)
			if err != nil {
				h.logger.Error("failed to decrypt request body",
					zap.String("implant_id", implantID),
					zap.Error(err))
				http.Error(w, "Failed to decrypt", http.StatusUnauthorized)
				return
			}

			// Parse the decrypted JSON
			var decryptedBody map[string]any
			if err := json.Unmarshal(decrypted, &decryptedBody); err != nil {
				h.logger.Error("failed to parse decrypted body",
					zap.String("implant_id", implantID),
					zap.Error(err))
				http.Error(w, "Failed to parse body", http.StatusBadRequest)
				return
			}

			ctx = context.WithValue(ctx, decryptedBodyContextKey, decryptedBody)

			// Restore body for potential re-reading
			r.Body = io.NopCloser(bytes.NewReader(decrypted))
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
