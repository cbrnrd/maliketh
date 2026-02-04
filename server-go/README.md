# Maliketh Server (Go)

Go rewrite of the Maliketh C2 server.

## Project Structure

```
server-go/
├── cmd/
│   ├── admin/          # Operator/Admin API server
│   ├── c2/             # C2/Implant API server
│   └── migrate/        # Database migration tool
├── internal/
│   ├── config/         # Configuration loading (env + YAML profiles)
│   ├── database/       # PostgreSQL connection and migrations
│   │   ├── queries/    # SQL queries for sqlc
│   │   └── repo/       # Generated repository code (sqlc)
│   ├── logging/        # Zap logger setup
│   ├── model/          # Domain models and opcodes
│   └── server/         # HTTP server (chi) and response helpers
├── migrations/         # PostgreSQL migrations
└── configs/            # Example configuration files
```

## Requirements

- Go 1.21+
- PostgreSQL 14+
- RabbitMQ 3.x
- sqlc (for code generation)

## Setup

1. Copy environment config:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

2. Run migrations:
   ```bash
   go run ./cmd/migrate -up -path ./migrations
   ```

3. Start the admin server:
   ```bash
   go run ./cmd/admin
   ```

4. Start the C2 server:
   ```bash
   go run ./cmd/c2 -profile ./configs/profile.example.yaml
   ```

## Development

### Regenerate sqlc code

After modifying SQL queries in `internal/database/queries/`:

```bash
sqlc generate
```

### Build all binaries

```bash
go build ./...
```

## Phase Status

- [x] Phase 1: Core Infrastructure
  - [x] Project scaffolding
  - [x] Config loading (env + YAML profiles)
  - [x] Database connection pool (pgx)
  - [x] SQL migrations + sqlc
  - [x] Logging (zap)
  - [x] HTTP server setup (chi)
- [x] Phase 2: Crypto Layer
  - [x] NaCl box encrypt/decrypt (X25519 + XSalsa20-Poly1305)
  - [x] Ed25519 signing/verification
  - [x] SecretBox (symmetric encryption)
  - [x] Key generation and base64 encoding utilities
  - [x] DecryptAndVerify for operator authentication
  - [x] Unit tests for all crypto operations
- [x] Phase 3: Database Repos
  - [x] Operator repository (CRUD, auth token management)
  - [x] Implant repository (CRUD, alias resolution)
  - [x] ImplantConfig repository (CRUD, profile-based creation)
  - [x] Task repository (CRUD, status management)
- [x] Phase 4: C2 Listener
  - [x] `/register` - implant registration with encrypted key exchange
  - [x] `/checkin` - get pending tasks (encrypted response)
  - [x] `/results` - submit task output (decrypted request)
  - [x] Implant authentication middleware
  - [x] Request decryption / response encryption
  - [x] RabbitMQ messaging integration
  - [x] Self-destruct handling
- [x] Phase 5: Admin Listener
  - [x] Token-based authentication (decrypt/verify challenge)
  - [x] Auth token generation with 6-hour expiry
  - [x] Verified middleware for protected endpoints
  - [x] Stats endpoint (implant count, task counts, uptime)
  - [x] Task CRUD (list, add, get result, delete)
  - [x] Implant management (list, get config, update config, kill)
  - [x] Alias management (create, delete, list, resolve)
  - [x] Admin revoke operator endpoint
- [x] Phase 6: Implant Builder
  - [x] BuildOptions with validation and defaults
  - [x] C++ implant build (Windows/mingw)
  - [x] Go implant build (cross-platform)
  - [x] Build endpoint with base64-encoded artifact response
  - [x] Builder info endpoint
- [x] Phase 7: Testing & Deployment
  - [x] Integration tests with testcontainers
  - [x] Dockerfile (multi-stage build with separate targets)
  - [x] docker-compose.yml with health checks
  - [x] Makefile with common commands
