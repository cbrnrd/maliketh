package crypto

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
)

const SecretKeySize = 32

// SecretBoxEncrypt encrypts data using NaCl secretbox (XSalsa20-Poly1305).
// Returns base64-encoded ciphertext (nonce prepended).
func SecretBoxEncrypt(key *[SecretKeySize]byte, plaintext []byte) (string, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	encrypted := secretbox.Seal(nonce[:], plaintext, nonce, key)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// SecretBoxDecrypt decrypts base64-encoded NaCl secretbox ciphertext.
func SecretBoxDecrypt(key *[SecretKeySize]byte, ciphertextB64 string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return SecretBoxDecryptRaw(key, ciphertext)
}

// SecretBoxDecryptRaw decrypts raw NaCl secretbox ciphertext (nonce prepended).
func SecretBoxDecryptRaw(key *[SecretKeySize]byte, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < NonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	var nonce [NonceSize]byte
	copy(nonce[:], ciphertext[:NonceSize])

	plaintext, ok := secretbox.Open(nil, ciphertext[NonceSize:], &nonce, key)
	if !ok {
		return nil, fmt.Errorf("decryption failed")
	}

	return plaintext, nil
}

// SecretBoxDecryptB64Key decrypts using a base64-encoded key.
func SecretBoxDecryptB64Key(keyB64 string, ciphertextB64 string) ([]byte, error) {
	key, err := DecodeSecretKey(keyB64)
	if err != nil {
		return nil, fmt.Errorf("invalid secret key: %w", err)
	}

	return SecretBoxDecrypt(key, ciphertextB64)
}

// SecretBoxDecryptB64KeyRaw decrypts raw ciphertext using a base64-encoded key.
func SecretBoxDecryptB64KeyRaw(keyB64 string, ciphertext []byte) ([]byte, error) {
	key, err := DecodeSecretKey(keyB64)
	if err != nil {
		return nil, fmt.Errorf("invalid secret key: %w", err)
	}

	return SecretBoxDecryptRaw(key, ciphertext)
}

// DecodeSecretKey decodes a base64-encoded secret key.
func DecodeSecretKey(b64 string) (*[SecretKeySize]byte, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	if len(data) != SecretKeySize {
		return nil, fmt.Errorf("invalid key size: got %d, want %d", len(data), SecretKeySize)
	}

	var key [SecretKeySize]byte
	copy(key[:], data)
	return &key, nil
}

// GenerateSecretKey generates a random 32-byte secret key.
func GenerateSecretKey() (*[SecretKeySize]byte, error) {
	var key [SecretKeySize]byte
	if err := RandomBytes(key[:]); err != nil {
		return nil, err
	}
	return &key, nil
}

// GenerateSecretKeyB64 generates a secret key and returns it as base64.
func GenerateSecretKeyB64() (string, error) {
	key, err := GenerateSecretKey()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key[:]), nil
}
