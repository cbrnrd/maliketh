.PHONY: all build test test-unit test-integration clean docker docker-up docker-down lint vet fmt sqlc migrate

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOFMT=gofmt
GOMOD=$(GOCMD) mod

# Binary names
ADMIN_BINARY=admin
C2_BINARY=c2
MIGRATE_BINARY=migrate

# Build directories
BUILD_DIR=build

all: clean build test-unit

build:
	@echo "Building binaries..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-s -w" -o $(BUILD_DIR)/$(ADMIN_BINARY) ./cmd/admin
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-s -w" -o $(BUILD_DIR)/$(C2_BINARY) ./cmd/c2
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-s -w" -o $(BUILD_DIR)/$(MIGRATE_BINARY) ./cmd/migrate
	@echo "Build complete: $(BUILD_DIR)/"

build-linux:
	@echo "Building Linux binaries..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o $(BUILD_DIR)/$(ADMIN_BINARY)-linux-amd64 ./cmd/admin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o $(BUILD_DIR)/$(C2_BINARY)-linux-amd64 ./cmd/c2
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o $(BUILD_DIR)/$(MIGRATE_BINARY)-linux-amd64 ./cmd/migrate

test: test-unit

test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v -race ./internal/...

test-integration:
	@echo "Running integration tests (requires Docker)..."
	$(GOTEST) -v -tags=integration ./tests/...

test-all: test-unit test-integration

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	$(GOCMD) clean -testcache

# Docker commands
docker:
	@echo "Building Docker images..."
	docker build -t maliketh-server:latest .
	docker build -t maliketh-admin:latest --target admin .
	docker build -t maliketh-c2:latest --target c2 .
	docker build -t maliketh-migrate:latest --target migrate .

docker-up:
	@echo "Starting services..."
	docker-compose up -d

docker-down:
	@echo "Stopping services..."
	docker-compose down

docker-logs:
	docker-compose logs -f

# Code quality
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed" && exit 1)
	golangci-lint run ./...

vet:
	@echo "Running go vet..."
	$(GOVET) ./...

fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

fmt-check:
	@echo "Checking code formatting..."
	@test -z "$$($(GOFMT) -l .)" || (echo "Code is not formatted. Run 'make fmt'" && exit 1)

# Dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Database
sqlc:
	@echo "Generating sqlc code..."
	sqlc generate

migrate-up:
	@echo "Running migrations..."
	$(GOCMD) run ./cmd/migrate -up -path ./migrations

migrate-down:
	@echo "Rolling back migrations..."
	$(GOCMD) run ./cmd/migrate -down -path ./migrations

# Development
dev-admin:
	@echo "Starting admin server in development mode..."
	$(GOCMD) run ./cmd/admin

dev-c2:
	@echo "Starting C2 server in development mode..."
	$(GOCMD) run ./cmd/c2 -profile ./configs/profile.example.yaml

# Operator management
create-operator:
	@echo "Creating operator..."
	@test -n "$(USERNAME)" || (echo "Usage: make create-operator USERNAME=<name> [ROLE=admin|operator]" && exit 1)
	$(GOCMD) run ./cmd/createop -username $(USERNAME) -role $(or $(ROLE),operator) -json

# Help
help:
	@echo "Maliketh Server (Go) - Available targets:"
	@echo ""
	@echo "  build            Build all binaries"
	@echo "  build-linux      Build Linux binaries (cross-compile)"
	@echo "  test             Run unit tests"
	@echo "  test-unit        Run unit tests"
	@echo "  test-integration Run integration tests (requires Docker)"
	@echo "  test-all         Run all tests"
	@echo "  clean            Clean build artifacts"
	@echo ""
	@echo "  docker           Build Docker images"
	@echo "  docker-up        Start services with docker-compose"
	@echo "  docker-down      Stop services"
	@echo "  docker-logs      Follow service logs"
	@echo ""
	@echo "  lint             Run golangci-lint"
	@echo "  vet              Run go vet"
	@echo "  fmt              Format code"
	@echo "  fmt-check        Check code formatting"
	@echo ""
	@echo "  deps             Download and tidy dependencies"
	@echo "  sqlc             Generate sqlc code"
	@echo "  migrate-up       Run database migrations"
	@echo "  migrate-down     Rollback database migrations"
	@echo ""
	@echo "  dev-admin        Start admin server (development)"
	@echo "  dev-c2           Start C2 server (development)"
	@echo ""
	@echo "  create-operator  Create a new operator (USERNAME=name [ROLE=admin|operator])"
