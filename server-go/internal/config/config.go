package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Postgres PostgresConfig
	RabbitMQ RabbitMQConfig
	Admin    ServerConfig
	C2       ServerConfig
	LogLevel string
	Env      string
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

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
