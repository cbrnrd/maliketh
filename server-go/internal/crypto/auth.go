package crypto

import (
	"fmt"
)

// OperatorKeys holds the cryptographic keys for an operator.
type OperatorKeys struct {
	PublicKey string // Base64 X25519 public key
	VerifyKey string // Base64 Ed25519 verify key
}

// DecryptAndVerify decrypts a message from an operator and verifies their signature.
// The message flow is: operator signs message, then encrypts with server's public key.
// Server decrypts with its private key, then verifies the operator's signature.
func DecryptAndVerify(
	serverPrivKeyB64 string,
	operatorPubKeyB64 string,
	operatorVerifyKeyB64 string,
	ciphertextB64 string,
) ([]byte, error) {
	// First decrypt using box (server priv + operator pub)
	decrypted, err := DecryptB64(operatorPubKeyB64, serverPrivKeyB64, ciphertextB64)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	// Then verify the signature
	message, err := VerifyRaw(operatorVerifyKeyB64, decrypted)
	if err != nil {
		return nil, fmt.Errorf("signature verification failed: %w", err)
	}

	return message, nil
}

// DecryptAndVerifyRaw is like DecryptAndVerify but takes raw ciphertext bytes.
func DecryptAndVerifyRaw(
	serverPrivKeyB64 string,
	operatorPubKeyB64 string,
	operatorVerifyKeyB64 string,
	ciphertext []byte,
) ([]byte, error) {
	// First decrypt using box (server priv + operator pub)
	decrypted, err := DecryptB64Raw(operatorPubKeyB64, serverPrivKeyB64, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	// Then verify the signature
	message, err := VerifyRaw(operatorVerifyKeyB64, decrypted)
	if err != nil {
		return nil, fmt.Errorf("signature verification failed: %w", err)
	}

	return message, nil
}

// SignAndEncrypt signs a message and encrypts it for the recipient.
// Used by operator to send authenticated messages to server.
func SignAndEncrypt(
	signingKeyB64 string,
	senderPrivKeyB64 string,
	recipientPubKeyB64 string,
	message []byte,
) (string, error) {
	// Load the signing key
	signingKey, err := DecodeSigningKey(signingKeyB64)
	if err != nil {
		return "", fmt.Errorf("invalid signing key: %w", err)
	}

	// Create signing key pair for signing
	kp := &SigningKeyPair{SigningKey: signingKey}

	// Sign the message (prepends signature)
	signed := kp.Sign(message)

	// Encrypt the signed message
	encrypted, err := EncryptB64(recipientPubKeyB64, senderPrivKeyB64, signed)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	return encrypted, nil
}
