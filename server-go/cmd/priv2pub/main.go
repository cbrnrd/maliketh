package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/curve25519"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: priv2pub <base64-private-key>")
		os.Exit(1)
	}

	privB64 := os.Args[1]
	privBytes, err := base64.StdEncoding.DecodeString(privB64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid base64:", err)
		os.Exit(1)
	}

	var priv, pub [32]byte
	copy(priv[:], privBytes)
	curve25519.ScalarBaseMult(&pub, &priv)

	fmt.Println(base64.StdEncoding.EncodeToString(pub[:]))
}
