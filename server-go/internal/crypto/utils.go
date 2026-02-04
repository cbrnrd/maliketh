package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// GenerateNonce generates a random 24-byte nonce for NaCl.
func GenerateNonce() (*[NonceSize]byte, error) {
	var nonce [NonceSize]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	return &nonce, nil
}

// RandomBytes fills the provided slice with random bytes.
func RandomBytes(b []byte) error {
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return nil
}

// RandomHex generates a random hex string of the specified byte length.
// The resulting string will be 2*length characters.
func RandomHex(length int) (string, error) {
	b := make([]byte, length)
	if err := RandomBytes(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// RandomID generates a random ID of the specified length (in hex characters).
func RandomID(length int) string {
	// We need length/2 bytes to get length hex characters
	byteLen := (length + 1) / 2
	b := make([]byte, byteLen)
	rand.Read(b) // Ignore error for ID generation
	return hex.EncodeToString(b)[:length]
}

// Pad pads data to a multiple of 16 bytes using PKCS7.
func Pad(data []byte) []byte {
	padding := 16 - len(data)%16
	padBytes := make([]byte, padding)
	for i := range padBytes {
		padBytes[i] = byte(padding)
	}
	return append(data, padBytes...)
}

// Unpad removes PKCS7 padding.
func Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	padding := int(data[len(data)-1])
	if padding > 16 || padding > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}

	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padding], nil
}
