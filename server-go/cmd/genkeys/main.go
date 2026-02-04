package main

import (
	"fmt"

	"github.com/cbrnrd/maliketh/server-go/internal/crypto"
)

func main() {
	priv, pub, err := crypto.GenerateKeyPairB64()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SERVER_PRIVATE_KEY=%s\n", priv)
	fmt.Printf("SERVER_PUBLIC_KEY=%s\n", pub)
}
