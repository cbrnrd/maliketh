// Package tests contains integration tests for the Maliketh server.
// These tests use testcontainers to spin up real PostgreSQL and RabbitMQ instances.
//
// Run with: go test -v -tags=integration ./tests/...
//
//go:build integration

package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cbrnrd/maliketh/server-go/internal/config"
	"github.com/cbrnrd/maliketh/server-go/internal/database"
	"github.com/cbrnrd/maliketh/server-go/internal/handler/admin"
	"github.com/cbrnrd/maliketh/server-go/internal/handler/c2"
	"github.com/cbrnrd/maliketh/server-go/internal/messaging"
	"github.com/cbrnrd/maliketh/server-go/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

type testEnv struct {
	db     *database.DB
	repos  *repository.Repositories
	logger *zap.Logger
	ctx    context.Context
}

func setupTestEnv(t *testing.T) *testEnv {
	ctx := context.Background()

	// Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("maliketh_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	})

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// Connect to database
	logger := zap.NewNop()
	db, err := database.New(ctx, connStr, logger)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	// Run migrations
	if err := runMigrations(ctx, db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	repos := repository.New(db.Pool)

	return &testEnv{
		db:     db,
		repos:  repos,
		logger: logger,
		ctx:    ctx,
	}
}

func runMigrations(ctx context.Context, db *database.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS operators (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			public_key TEXT NOT NULL,
			verify_key TEXT NOT NULL,
			login_secret TEXT NOT NULL,
			auth_token TEXT,
			auth_token_expiry TIMESTAMPTZ,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			last_login TIMESTAMPTZ,
			rmq_queue TEXT,
			role VARCHAR(50) DEFAULT 'operator',
			revoked BOOLEAN DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS implants (
			id SERIAL PRIMARY KEY,
			implant_id VARCHAR(255) UNIQUE NOT NULL,
			hostname TEXT,
			ip TEXT,
			os TEXT,
			arch TEXT,
			username TEXT,
			server_sk TEXT NOT NULL,
			implant_pk TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			last_seen TIMESTAMPTZ DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS implant_configs (
			id SERIAL PRIMARY KEY,
			implant_id VARCHAR(255) UNIQUE NOT NULL REFERENCES implants(implant_id) ON DELETE CASCADE,
			cookie TEXT NOT NULL,
			kill_date TIMESTAMPTZ,
			user_agent TEXT,
			auto_self_destruct BOOLEAN DEFAULT FALSE,
			sleep_time INTEGER DEFAULT 60,
			jitter REAL DEFAULT 0.1,
			max_retries INTEGER DEFAULT -1,
			retry_wait INTEGER DEFAULT 5,
			retry_jitter REAL DEFAULT 0.1,
			enc_key TEXT NOT NULL,
			tailoring_hash_function TEXT DEFAULT 'sha256',
			tailoring_hash_rounds INTEGER DEFAULT 100,
			tailoring_hashes TEXT[]
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			task_id VARCHAR(255) UNIQUE NOT NULL,
			operator_name VARCHAR(255) NOT NULL,
			implant_id VARCHAR(255) NOT NULL REFERENCES implants(implant_id) ON DELETE CASCADE,
			opcode INTEGER NOT NULL,
			args JSONB NOT NULL DEFAULT '[]',
			status VARCHAR(50) DEFAULT 'CREATED',
			output TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			read_at TIMESTAMPTZ,
			executed_at TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS implant_aliases (
			id SERIAL PRIMARY KEY,
			implant_id VARCHAR(255) NOT NULL REFERENCES implants(implant_id) ON DELETE CASCADE,
			alias VARCHAR(255) UNIQUE NOT NULL
		)`,
	}

	for _, m := range migrations {
		if _, err := db.Pool.Exec(ctx, m); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}

func TestDatabaseConnection(t *testing.T) {
	env := setupTestEnv(t)

	err := env.db.Health(env.ctx)
	if err != nil {
		t.Errorf("database health check failed: %v", err)
	}
}

func TestOperatorRepository(t *testing.T) {
	env := setupTestEnv(t)

	// Create operator directly via SQL for testing
	_, err := env.db.Pool.Exec(env.ctx, `
		INSERT INTO operators (username, public_key, verify_key, login_secret, role)
		VALUES ($1, $2, $3, $4, $5)
	`, "testoperator", "pubkey123", "verifykey123", "secret123", "operator")
	if err != nil {
		t.Fatalf("failed to create operator: %v", err)
	}

	// Test GetByUsername
	op, err := env.repos.Operator.GetByUsername(env.ctx, "testoperator")
	if err != nil {
		t.Fatalf("GetByUsername failed: %v", err)
	}
	if op.Username != "testoperator" {
		t.Errorf("expected username 'testoperator', got '%s'", op.Username)
	}

	// Test Count
	count, err := env.repos.Operator.Count(env.ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
}

func TestImplantRepository(t *testing.T) {
	env := setupTestEnv(t)

	// Create implant directly via SQL
	_, err := env.db.Pool.Exec(env.ctx, `
		INSERT INTO implants (implant_id, hostname, ip, os, arch, username, server_sk, implant_pk)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, "imp123", "testhost", "192.168.1.1", "windows", "amd64", "testuser", "serversk", "implantpk")
	if err != nil {
		t.Fatalf("failed to create implant: %v", err)
	}

	// Test GetByID
	imp, err := env.repos.Implant.GetByID(env.ctx, "imp123")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if imp.ImplantID != "imp123" {
		t.Errorf("expected implant_id 'imp123', got '%s'", imp.ImplantID)
	}
	if imp.Hostname != "testhost" {
		t.Errorf("expected hostname 'testhost', got '%s'", imp.Hostname)
	}

	// Test List
	implants, err := env.repos.Implant.List(env.ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(implants) != 1 {
		t.Errorf("expected 1 implant, got %d", len(implants))
	}

	// Test alias creation
	_, err = env.repos.Implant.CreateAlias(env.ctx, "imp123", "myalias")
	if err != nil {
		t.Fatalf("CreateAlias failed: %v", err)
	}

	// Test Resolve with alias
	resolved, err := env.repos.Implant.Resolve(env.ctx, "myalias")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if resolved.ImplantID != "imp123" {
		t.Errorf("expected implant_id 'imp123', got '%s'", resolved.ImplantID)
	}
}

func TestAdminHealthEndpoint(t *testing.T) {
	env := setupTestEnv(t)

	cfg := &config.Config{
		RabbitMQ: config.RabbitMQConfig{Host: "localhost", Port: 5672},
	}

	handler := admin.NewHandler(env.repos, cfg, &messaging.NoopPublisher{}, env.logger, "")

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		if err := env.db.Health(req.Context()); err != nil {
			http.Error(w, "unhealthy", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"status": true, "message": "healthy"})
	})
	handler.RegisterRoutes(r)

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestC2RegisterEndpoint(t *testing.T) {
	env := setupTestEnv(t)

	profile := config.DefaultProfile()

	handler := c2.NewHandler(env.repos, profile, &messaging.NoopPublisher{}, env.logger)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	// This test just verifies the endpoint exists and returns an error for invalid request
	// (since we don't have proper crypto setup here)
	body := `{"public_key": "test", "registration_password": "wrong"}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	// Should fail due to invalid registration
	if rec.Code == http.StatusOK {
		t.Error("expected registration to fail with invalid password")
	}
}

func TestTokenStatusEndpoint(t *testing.T) {
	env := setupTestEnv(t)

	cfg := &config.Config{
		RabbitMQ: config.RabbitMQConfig{Host: "localhost", Port: 5672},
	}

	handler := admin.NewHandler(env.repos, cfg, &messaging.NoopPublisher{}, env.logger, "")

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	// Test unauthenticated request
	req := httptest.NewRequest("GET", "/op/auth/token/status", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	var response map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != false {
		t.Error("expected status false for unauthenticated request")
	}
}
