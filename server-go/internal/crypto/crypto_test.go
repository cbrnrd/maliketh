package crypto

import (
	"bytes"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	if kp.PrivateKey == [KeySize]byte{} {
		t.Error("private key is empty")
	}
	if kp.PublicKey == [KeySize]byte{} {
		t.Error("public key is empty")
	}
}

func TestGenerateKeyPairB64(t *testing.T) {
	priv, pub, err := GenerateKeyPairB64()
	if err != nil {
		t.Fatalf("GenerateKeyPairB64 failed: %v", err)
	}

	if len(priv) == 0 {
		t.Error("private key is empty")
	}
	if len(pub) == 0 {
		t.Error("public key is empty")
	}

	// Verify we can decode them back
	_, err = DecodePrivateKey(priv)
	if err != nil {
		t.Errorf("failed to decode private key: %v", err)
	}
	_, err = DecodePublicKey(pub)
	if err != nil {
		t.Errorf("failed to decode public key: %v", err)
	}
}

func TestBoxEncryptDecrypt(t *testing.T) {
	// Generate sender and recipient key pairs
	sender, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate sender key pair: %v", err)
	}

	recipient, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate recipient key pair: %v", err)
	}

	plaintext := []byte("Hello, World!")

	// Encrypt with recipient's public key and sender's private key
	ciphertext, err := Encrypt(&recipient.PublicKey, &sender.PrivateKey, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Decrypt with sender's public key and recipient's private key
	decrypted, err := Decrypt(&sender.PublicKey, &recipient.PrivateKey, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("decrypted text doesn't match: got %q, want %q", decrypted, plaintext)
	}
}

func TestBoxEncryptDecryptB64(t *testing.T) {
	senderPriv, senderPub, _ := GenerateKeyPairB64()
	recipientPriv, recipientPub, _ := GenerateKeyPairB64()

	plaintext := []byte("Test message for B64 encryption")

	ciphertext, err := EncryptB64(recipientPub, senderPriv, plaintext)
	if err != nil {
		t.Fatalf("EncryptB64 failed: %v", err)
	}

	decrypted, err := DecryptB64(senderPub, recipientPriv, ciphertext)
	if err != nil {
		t.Fatalf("DecryptB64 failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("decrypted text doesn't match: got %q, want %q", decrypted, plaintext)
	}
}

func TestSigningKeyPair(t *testing.T) {
	kp, err := GenerateSigningKeyPair()
	if err != nil {
		t.Fatalf("GenerateSigningKeyPair failed: %v", err)
	}

	message := []byte("Message to sign")

	// Sign the message
	signed := kp.Sign(message)

	// Verify the signature
	verified, err := Verify(kp.VerifyKey, signed)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if !bytes.Equal(message, verified) {
		t.Errorf("verified message doesn't match: got %q, want %q", verified, message)
	}
}

func TestSigningKeyPairB64(t *testing.T) {
	signingKeyB64, verifyKeyB64, err := GenerateSigningKeyPairB64()
	if err != nil {
		t.Fatalf("GenerateSigningKeyPairB64 failed: %v", err)
	}

	// Load the signing key
	signingKey, err := DecodeSigningKey(signingKeyB64)
	if err != nil {
		t.Fatalf("DecodeSigningKey failed: %v", err)
	}

	kp := &SigningKeyPair{SigningKey: signingKey}
	message := []byte("B64 signing test")

	signedB64 := kp.SignB64(message)

	verified, err := VerifyB64(verifyKeyB64, signedB64)
	if err != nil {
		t.Fatalf("VerifyB64 failed: %v", err)
	}

	if !bytes.Equal(message, verified) {
		t.Errorf("verified message doesn't match: got %q, want %q", verified, message)
	}
}

func TestSecretBox(t *testing.T) {
	key, err := GenerateSecretKey()
	if err != nil {
		t.Fatalf("GenerateSecretKey failed: %v", err)
	}

	plaintext := []byte("Secret message")

	ciphertext, err := SecretBoxEncrypt(key, plaintext)
	if err != nil {
		t.Fatalf("SecretBoxEncrypt failed: %v", err)
	}

	decrypted, err := SecretBoxDecrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("SecretBoxDecrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("decrypted text doesn't match: got %q, want %q", decrypted, plaintext)
	}
}

func TestDecryptAndVerify(t *testing.T) {
	// Generate server key pair (for encryption)
	serverPriv, serverPub, _ := GenerateKeyPairB64()

	// Generate operator key pairs (for encryption and signing)
	operatorPriv, operatorPub, _ := GenerateKeyPairB64()
	signingKeyB64, verifyKeyB64, _ := GenerateSigningKeyPairB64()

	message := []byte("authenticated message")

	// Operator: sign then encrypt
	encrypted, err := SignAndEncrypt(signingKeyB64, operatorPriv, serverPub, message)
	if err != nil {
		t.Fatalf("SignAndEncrypt failed: %v", err)
	}

	// Server: decrypt then verify
	decrypted, err := DecryptAndVerify(serverPriv, operatorPub, verifyKeyB64, encrypted)
	if err != nil {
		t.Fatalf("DecryptAndVerify failed: %v", err)
	}

	if !bytes.Equal(message, decrypted) {
		t.Errorf("decrypted message doesn't match: got %q, want %q", decrypted, message)
	}
}

func TestRandomID(t *testing.T) {
	id := RandomID(16)
	if len(id) != 16 {
		t.Errorf("RandomID length mismatch: got %d, want 16", len(id))
	}

	// Generate another and ensure they're different
	id2 := RandomID(16)
	if id == id2 {
		t.Error("RandomID generated identical IDs")
	}
}

func TestPadUnpad(t *testing.T) {
	testCases := [][]byte{
		[]byte("hello"),
		[]byte("exactly16bytes!!"),
		[]byte("a"),
		[]byte(""),
	}

	for _, original := range testCases {
		padded := Pad(original)
		if len(padded)%16 != 0 {
			t.Errorf("Pad didn't produce 16-byte aligned output: len=%d", len(padded))
		}

		unpadded, err := Unpad(padded)
		if err != nil {
			t.Errorf("Unpad failed for %q: %v", original, err)
			continue
		}

		if !bytes.Equal(original, unpadded) {
			t.Errorf("Unpad mismatch: got %q, want %q", unpadded, original)
		}
	}
}
