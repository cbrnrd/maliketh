package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"emperror.dev/errors"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

const NONCE_SIZE = 24

// Create a new public/private key pair
func CreateKeyPair() (publicKey, privateKey *[32]byte, err error) {
	return box.GenerateKey(rand.Reader)
}

func CreateBase64KeyPair() (publicKey, privateKey string, err error) {
	pub, priv, err := CreateKeyPair()
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(pub[:]), base64.StdEncoding.EncodeToString(priv[:]), nil
}

// Encrypt a plaintext using a public key and a private key
func Encrypt(plaintext []byte, serverPublicKey *[32]byte, implantPrivateKey *[32]byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, errors.Wrap(err, "Could not read from random")
	}
	// This encrypts msg and appends the result to the nonce.
	encrypted := box.Seal(nonce[:], plaintext, &nonce, serverPublicKey, implantPrivateKey)
	return encrypted, nil
}

// Encrypt a plaintext using a base64 encoded public and private key
func EncryptStringToBase64(plaintext string, serverPublicKeyB64, implantPrivateKeyB64 string) (string, error) {
	fmt.Printf("Encrypting \"%s\" with public key \"%s\" and private key \"%s\"\n", plaintext, serverPublicKeyB64, implantPrivateKeyB64)
	pub, err := base64.StdEncoding.DecodeString(serverPublicKeyB64)
	if err != nil {
		return "", err
	}
	priv, err := base64.StdEncoding.DecodeString(implantPrivateKeyB64)
	if err != nil {
		return "", err
	}
	enc, err := Encrypt([]byte(plaintext), (*[32]byte)(pub), (*[32]byte)(priv))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(enc), nil
}

// Encrypt a plaintext using a base64 encoded shared key
func EncryptSecretBoxBase64(plaintext []byte, key string) (string, error) {
	k, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	enc, err := EncryptSecretBox(plaintext, (*[32]byte)(k))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(enc), nil
}

func EncryptSecretBox(plaintext []byte, key *[32]byte) ([]byte, error) {

	nonce := new([NONCE_SIZE]byte)
	// Read bytes from random and put them in nonce until it is full.
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, errors.Wrap(err, "Could not read from random")
	}

	// out will hold the nonce and the encrypted message (ciphertext)
	out := make([]byte, NONCE_SIZE)
	// Copy the nonce to the start of out
	copy(out, nonce[:])
	// Encrypt the message and append it to out, assign the result to out
	out = secretbox.Seal(out, plaintext, nonce, key)

	return out, nil
}

func DecryptB64SecretBox(ciphertextB64 string, keyB64 string) (string, error) {
	k, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return "", err
	}
	c, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}
	dec, ok := DecryptSecretBox(c, ([32]byte)(k))
	if !ok {
		return "", errors.New("failed to decrypt")
	}
	return string(dec), nil
}

func DecryptSecretBox(ciphertext []byte, key [32]byte) ([]byte, bool) {
	var nonce [NONCE_SIZE]byte
	copy(nonce[:], ciphertext[:NONCE_SIZE])
	return secretbox.Open(nil, ciphertext[NONCE_SIZE:], &nonce, &key)
}

/*************
* Decryption *
**************/

// Decrypt a ciphertext using a public key and a private key. The nonce is expected to be prepended to the ciphertext.
func Decrypt(ciphertext []byte, publicKey *[32]byte, privateKey *[32]byte) ([]byte, bool) {
	var decryptNonce [NONCE_SIZE]byte
	copy(decryptNonce[:], ciphertext[:NONCE_SIZE])
	decrypted, ok := box.Open(nil, ciphertext[NONCE_SIZE:], &decryptNonce, publicKey, privateKey)
	if !ok {
		panic("decryption error")
	}
	return decrypted, true
}

func DecryptB64String(ciphertextb64 string, serverPublicKeyB64, implantPrivateKeyB64 string) (string, bool) {
	pub, err := base64.StdEncoding.DecodeString(serverPublicKeyB64)
	if err != nil {
		fmt.Println("error decoding public key")
		return "", false
	}
	priv, err := base64.StdEncoding.DecodeString(implantPrivateKeyB64)
	if err != nil {
		fmt.Println("error decoding private key")
		return "", false
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertextb64)
	if err != nil {
		fmt.Println("error decoding ciphertext")
		return "", false
	}

	dec, ok := Decrypt([]byte(decodedCiphertext), (*[32]byte)(pub), (*[32]byte)(priv))
	if !ok {
		fmt.Println("error decrypting")
		return "", false
	}
	return string(dec), true
}
