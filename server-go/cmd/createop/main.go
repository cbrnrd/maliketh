package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/config"
	"github.com/cbrnrd/maliketh/server-go/internal/crypto"
	"github.com/cbrnrd/maliketh/server-go/internal/database"
	"github.com/cbrnrd/maliketh/server-go/internal/model"
	"github.com/cbrnrd/maliketh/server-go/internal/repository"
	"go.uber.org/zap"
)

type OperatorCredentials struct {
	Name        string `json:"name"`
	Role        string `json:"role"`
	Public      string `json:"public"`
	Secret      string `json:"secret"`
	VerifyKey   string `json:"verify_key"`
	SigningKey  string `json:"signing_key"`
	LoginSecret string `json:"login_secret"`
	ServerPub   string `json:"server_pub"`
	RMQQueue    string `json:"rmq_queue"`
	C2          string `json:"c2"`
	C2Port      int    `json:"c2_port"`
}

func main() {
	username := flag.String("username", "", "operator username (required)")
	role := flag.String("role", "operator", "operator role (admin or operator)")
	outputJSON := flag.Bool("json", false, "output credentials as JSON")
	flag.Parse()

	if *username == "" {
		fmt.Fprintln(os.Stderr, "error: -username is required")
		flag.Usage()
		os.Exit(1)
	}

	if *role != "admin" && *role != "operator" {
		fmt.Fprintln(os.Stderr, "error: -role must be 'admin' or 'operator'")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.ConnectWithRetry(ctx, cfg.Postgres.DSN(), logger, 3, 2*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Generate X25519 key pair for encryption
	encPrivKey, encPubKey, err := crypto.GenerateKeyPairB64()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate encryption keys: %v\n", err)
		os.Exit(1)
	}

	// Generate Ed25519 key pair for signing
	signKey, verifyKey, err := crypto.GenerateSigningKeyPairB64()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate signing keys: %v\n", err)
		os.Exit(1)
	}

	// Generate login secret
	loginSecret, err := crypto.RandomHex(32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate login secret: %v\n", err)
		os.Exit(1)
	}

	// Generate RMQ queue name
	rmqQueue := fmt.Sprintf("op_%s_%s", *username, crypto.RandomID(8))

	op := &model.Operator{
		Username:    *username,
		PublicKey:   encPubKey,
		VerifyKey:   verifyKey,
		LoginSecret: loginSecret,
		RMQQueue:    rmqQueue,
		Role:        model.OperatorRole(*role),
	}

	opRepo := repository.NewOperatorRepository(db.Pool)
	created, err := opRepo.Create(ctx, op)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create operator: %v\n", err)
		os.Exit(1)
	}

	creds := OperatorCredentials{
		Name:        created.Username,
		Role:        string(created.Role),
		Public:      encPubKey,
		Secret:      encPrivKey,
		VerifyKey:   verifyKey,
		SigningKey:  signKey,
		LoginSecret: loginSecret,
		ServerPub:   cfg.ServerPublicKey,
		RMQQueue:    rmqQueue,
		C2:          cfg.Admin.Host,
		C2Port:      cfg.Admin.Port,
	}

	if *outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(creds)
	} else {
		fmt.Printf("Operator created successfully!\n\n")
		fmt.Printf("Name:         %s\n", creds.Name)
		fmt.Printf("Role:         %s\n", creds.Role)
		fmt.Printf("RMQ Queue:    %s\n", creds.RMQQueue)
		fmt.Printf("C2:           %s\n", creds.C2)
		fmt.Printf("C2 Port:      %d\n", creds.C2Port)
		fmt.Printf("\n--- SAVE THESE KEYS SECURELY ---\n\n")
		fmt.Printf("Public (X25519):       %s\n", creds.Public)
		fmt.Printf("Secret (X25519):       %s\n", creds.Secret)
		fmt.Printf("Verify Key (Ed25519):  %s\n", creds.VerifyKey)
		fmt.Printf("Signing Key (Ed25519): %s\n", creds.SigningKey)
		fmt.Printf("Login Secret:          %s\n", creds.LoginSecret)
		fmt.Printf("Server Public Key:     %s\n", creds.ServerPub)
	}
}
