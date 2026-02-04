package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const (
	SigningKeySize    = ed25519.SeedSize      // 32 bytes
	VerifyKeySize     = ed25519.PublicKeySize // 32 bytes
	SignatureSize     = ed25519.SignatureSize // 64 bytes
	SignedMessageSize = SignatureSize         // Signature is prepended
)

// SigningKeyPair represents an Ed25519 signing key pair.
type SigningKeyPair struct {
	SigningKey ed25519.PrivateKey
	VerifyKey  ed25519.PublicKey
}

// GenerateSigningKeyPair generates a new Ed25519 signing key pair.
func GenerateSigningKeyPair() (*SigningKeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signing key pair: %w", err)
	}

	return &SigningKeyPair{
		SigningKey: priv,
		VerifyKey:  pub,
	}, nil
}

// GenerateSigningKeyPairB64 generates a signing key pair and returns base64-encoded keys.
// Returns (signingKey, verifyKey, error).
// Note: The signing key returned is the seed (32 bytes) to match PyNaCl behavior.
func GenerateSigningKeyPairB64() (string, string, error) {
	kp, err := GenerateSigningKeyPair()
	if err != nil {
		return "", "", err
	}

	// PyNaCl stores the seed (first 32 bytes) as the signing key
	seed := kp.SigningKey.Seed()
	signingKeyB64 := base64.StdEncoding.EncodeToString(seed)
	verifyKeyB64 := base64.StdEncoding.EncodeToString(kp.VerifyKey)

	return signingKeyB64, verifyKeyB64, nil
}

// Sign signs a message and prepends the signature (NaCl style).
func (kp *SigningKeyPair) Sign(message []byte) []byte {
	signature := ed25519.Sign(kp.SigningKey, message)
	return append(signature, message...)
}

// SignB64 signs a message and returns base64-encoded signed message.
func (kp *SigningKeyPair) SignB64(message []byte) string {
	signed := kp.Sign(message)
	return base64.StdEncoding.EncodeToString(signed)
}

// Verify verifies a signed message and returns the original message.
// The signed message format is: signature (64 bytes) + original message.
func Verify(verifyKey ed25519.PublicKey, signedMessage []byte) ([]byte, error) {
	if len(signedMessage) < SignatureSize {
		return nil, fmt.Errorf("signed message too short")
	}

	signature := signedMessage[:SignatureSize]
	message := signedMessage[SignatureSize:]

	if !ed25519.Verify(verifyKey, message, signature) {
		return nil, fmt.Errorf("signature verification failed")
	}

	return message, nil
}

// VerifyB64 verifies using base64-encoded verify key and signed message.
func VerifyB64(verifyKeyB64 string, signedMessageB64 string) ([]byte, error) {
	verifyKey, err := DecodeVerifyKey(verifyKeyB64)
	if err != nil {
		return nil, fmt.Errorf("invalid verify key: %w", err)
	}

	signedMessage, err := base64.StdEncoding.DecodeString(signedMessageB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signed message: %w", err)
	}

	return Verify(verifyKey, signedMessage)
}

// VerifyRaw verifies using base64-encoded verify key and raw signed message bytes.
func VerifyRaw(verifyKeyB64 string, signedMessage []byte) ([]byte, error) {
	verifyKey, err := DecodeVerifyKey(verifyKeyB64)
	if err != nil {
		return nil, fmt.Errorf("invalid verify key: %w", err)
	}

	return Verify(verifyKey, signedMessage)
}

// DecodeVerifyKey decodes a base64-encoded Ed25519 verify (public) key.
func DecodeVerifyKey(b64 string) (ed25519.PublicKey, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	if len(data) != VerifyKeySize {
		return nil, fmt.Errorf("invalid verify key size: got %d, want %d", len(data), VerifyKeySize)
	}

	return ed25519.PublicKey(data), nil
}

// DecodeSigningKey decodes a base64-encoded Ed25519 signing key (seed).
func DecodeSigningKey(b64 string) (ed25519.PrivateKey, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	if len(data) != SigningKeySize {
		return nil, fmt.Errorf("invalid signing key size: got %d, want %d", len(data), SigningKeySize)
	}

	// Reconstruct the full private key from the seed
	return ed25519.NewKeyFromSeed(data), nil
}

// SigningKeyB64 returns the base64-encoded signing key seed.
func (kp *SigningKeyPair) SigningKeyB64() string {
	return base64.StdEncoding.EncodeToString(kp.SigningKey.Seed())
}

// VerifyKeyB64 returns the base64-encoded verify key.
func (kp *SigningKeyPair) VerifyKeyB64() string {
	return base64.StdEncoding.EncodeToString(kp.VerifyKey)
}
