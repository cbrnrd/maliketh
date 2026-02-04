package main

import (
	"encoding/json"
	"fmt"

	"github.com/cbrnrd/maliketh/server-go/internal/crypto"
)

func main() {
	encPriv, encPub, _ := crypto.GenerateKeyPairB64()
	signKey, verifyKey, _ := crypto.GenerateSigningKeyPairB64()
	loginSecret, _ := crypto.RandomHex(32)

	out := map[string]string{
		"public_key":   encPub,
		"private_key":  encPriv,
		"verify_key":   verifyKey,
		"signing_key":  signKey,
		"login_secret": loginSecret,
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}
