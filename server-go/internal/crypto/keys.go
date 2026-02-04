package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/curve25519"
)

// KeyPair represents an X25519 key pair for NaCl box encryption.
type KeyPair struct {
	PrivateKey [KeySize]byte
	PublicKey  [KeySize]byte
}

// GenerateKeyPair generates a new X25519 key pair.
func GenerateKeyPair() (*KeyPair, error) {
	var priv, pub [KeySize]byte

	if _, err := io.ReadFull(rand.Reader, priv[:]); err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	curve25519.ScalarBaseMult(&pub, &priv)

	return &KeyPair{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}

// GenerateKeyPairB64 generates a key pair and returns base64-encoded keys.
// Returns (privateKey, publicKey, error).
func GenerateKeyPairB64() (string, string, error) {
	kp, err := GenerateKeyPair()
	if err != nil {
		return "", "", err
	}

	privB64 := base64.StdEncoding.EncodeToString(kp.PrivateKey[:])
	pubB64 := base64.StdEncoding.EncodeToString(kp.PublicKey[:])

	return privB64, pubB64, nil
}

// DecodePublicKey decodes a base64-encoded public key.
func DecodePublicKey(b64 string) (*[KeySize]byte, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	if len(data) != KeySize {
		return nil, fmt.Errorf("invalid key size: got %d, want %d", len(data), KeySize)
	}

	var key [KeySize]byte
	copy(key[:], data)
	return &key, nil
}

// DecodePrivateKey decodes a base64-encoded private key.
func DecodePrivateKey(b64 string) (*[KeySize]byte, error) {
	return DecodePublicKey(b64) // Same format
}

// EncodePublicKey encodes a public key to base64.
func EncodePublicKey(key *[KeySize]byte) string {
	return base64.StdEncoding.EncodeToString(key[:])
}

// EncodePrivateKey encodes a private key to base64.
func EncodePrivateKey(key *[KeySize]byte) string {
	return base64.StdEncoding.EncodeToString(key[:])
}

// LoadKeyPairFromFiles loads a key pair from files containing base64-encoded keys.
func LoadKeyPairFromFiles(privPath, pubPath string) (*KeyPair, error) {
	privData, err := os.ReadFile(privPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	pubData, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	priv, err := DecodePrivateKey(string(privData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	pub, err := DecodePublicKey(string(pubData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	return &KeyPair{
		PrivateKey: *priv,
		PublicKey:  *pub,
	}, nil
}

// SaveKeyPairToFiles saves a key pair to files as base64-encoded strings.
func (kp *KeyPair) SaveToFiles(privPath, pubPath string) error {
	privB64 := EncodePrivateKey(&kp.PrivateKey)
	pubB64 := EncodePublicKey(&kp.PublicKey)

	if err := os.WriteFile(privPath, []byte(privB64), 0600); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	if err := os.WriteFile(pubPath, []byte(pubB64), 0644); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}

// PrivateKeyB64 returns the base64-encoded private key.
func (kp *KeyPair) PrivateKeyB64() string {
	return EncodePrivateKey(&kp.PrivateKey)
}

// PublicKeyB64 returns the base64-encoded public key.
func (kp *KeyPair) PublicKeyB64() string {
	return EncodePublicKey(&kp.PublicKey)
}
