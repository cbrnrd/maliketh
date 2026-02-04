package crypto

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/nacl/box"
)

const (
	KeySize   = 32
	NonceSize = 24
)

// Encrypt encrypts data using NaCl box (X25519 + XSalsa20-Poly1305).
// Returns base64-encoded ciphertext (nonce prepended).
func Encrypt(recipientPubKey, senderPrivKey *[KeySize]byte, plaintext []byte) (string, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	encrypted := box.Seal(nonce[:], plaintext, nonce, recipientPubKey, senderPrivKey)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// EncryptRaw encrypts and returns raw bytes (nonce prepended).
func EncryptRaw(recipientPubKey, senderPrivKey *[KeySize]byte, plaintext []byte) ([]byte, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	return box.Seal(nonce[:], plaintext, nonce, recipientPubKey, senderPrivKey), nil
}

// Decrypt decrypts base64-encoded NaCl box ciphertext.
func Decrypt(senderPubKey, recipientPrivKey *[KeySize]byte, ciphertextB64 string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return DecryptRaw(senderPubKey, recipientPrivKey, ciphertext)
}

// DecryptRaw decrypts raw NaCl box ciphertext (nonce prepended).
func DecryptRaw(senderPubKey, recipientPrivKey *[KeySize]byte, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < NonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	var nonce [NonceSize]byte
	copy(nonce[:], ciphertext[:NonceSize])

	plaintext, ok := box.Open(nil, ciphertext[NonceSize:], &nonce, senderPubKey, recipientPrivKey)
	if !ok {
		return nil, fmt.Errorf("decryption failed")
	}

	return plaintext, nil
}

// EncryptB64 encrypts using base64-encoded keys.
func EncryptB64(recipientPubKeyB64, senderPrivKeyB64 string, plaintext []byte) (string, error) {
	recipientPub, err := DecodePublicKey(recipientPubKeyB64)
	if err != nil {
		return "", fmt.Errorf("invalid recipient public key: %w", err)
	}

	senderPriv, err := DecodePrivateKey(senderPrivKeyB64)
	if err != nil {
		return "", fmt.Errorf("invalid sender private key: %w", err)
	}

	return Encrypt(recipientPub, senderPriv, plaintext)
}

// DecryptB64 decrypts using base64-encoded keys.
func DecryptB64(senderPubKeyB64, recipientPrivKeyB64 string, ciphertextB64 string) ([]byte, error) {
	senderPub, err := DecodePublicKey(senderPubKeyB64)
	if err != nil {
		return nil, fmt.Errorf("invalid sender public key: %w", err)
	}

	recipientPriv, err := DecodePrivateKey(recipientPrivKeyB64)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient private key: %w", err)
	}

	return Decrypt(senderPub, recipientPriv, ciphertextB64)
}

// DecryptB64Raw decrypts raw bytes using base64-encoded keys.
func DecryptB64Raw(senderPubKeyB64, recipientPrivKeyB64 string, ciphertext []byte) ([]byte, error) {
	senderPub, err := DecodePublicKey(senderPubKeyB64)
	if err != nil {
		return nil, fmt.Errorf("invalid sender public key: %w", err)
	}

	recipientPriv, err := DecodePrivateKey(recipientPrivKeyB64)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient private key: %w", err)
	}

	return DecryptRaw(senderPub, recipientPriv, ciphertext)
}
