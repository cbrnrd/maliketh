package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/curve25519"
)

type Config struct {
	Postgres        PostgresConfig
	RabbitMQ        RabbitMQConfig
	Admin           ServerConfig
	C2              ServerConfig
	LogLevel        string
	Env             string
	ServerPublicKey string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.Database)
}

type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func (r RabbitMQConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		r.User, r.Password, r.Host, r.Port)
}

type ServerConfig struct {
	Host string
	Port int
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func Load() (*Config, error) {
	pgPort, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_PORT: %w", err)
	}

	rmqPort, err := strconv.Atoi(getEnv("RABBITMQ_PORT", "5672"))
	if err != nil {
		return nil, fmt.Errorf("invalid RABBITMQ_PORT: %w", err)
	}

	adminPort, err := strconv.Atoi(getEnv("ADMIN_PORT", "5000"))
	if err != nil {
		return nil, fmt.Errorf("invalid ADMIN_PORT: %w", err)
	}

	c2Port, err := strconv.Atoi(getEnv("C2_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid C2_PORT: %w", err)
	}

	cfg := &Config{
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "postgres"),
			Port:     pgPort,
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: getEnv("POSTGRES_DB", "maliketh"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:     getEnv("RABBITMQ_HOST", "rabbitmq"),
			Port:     rmqPort,
			User:     getEnv("RABBITMQ_USER", "guest"),
			Password: getEnv("RABBITMQ_PASSWORD", "guest"),
		},
		Admin: ServerConfig{
			Host: getEnv("ADMIN_HOST", "0.0.0.0"),
			Port: adminPort,
		},
		C2: ServerConfig{
			Host: getEnv("C2_HOST", "0.0.0.0"),
			Port: c2Port,
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Env:      getEnv("ENV", "development"),
	}

	if cfg.Postgres.Password == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD must be set")
	}

	serverPrivKeyB64 := os.Getenv("SERVER_PRIVATE_KEY")
	if serverPrivKeyB64 != "" {
		pubKey, err := derivePublicKeyB64(serverPrivKeyB64)
		if err != nil {
			return nil, fmt.Errorf("failed to derive server public key: %w", err)
		}
		cfg.ServerPublicKey = pubKey
	}

	return cfg, nil
}

func derivePublicKeyB64(privKeyB64 string) (string, error) {
	privKey, err := base64.StdEncoding.DecodeString(privKeyB64)
	if err != nil {
		return "", err
	}
	if len(privKey) != 32 {
		return "", fmt.Errorf("invalid private key length: %d", len(privKey))
	}
	var priv, pub [32]byte
	copy(priv[:], privKey)
	curve25519.ScalarBaseMult(&pub, &priv)
	return base64.StdEncoding.EncodeToString(pub[:]), nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
