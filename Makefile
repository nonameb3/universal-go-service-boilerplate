# Universal Go Service Boilerplate - Makefile
# Production-ready commands for development and deployment

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=universal-service
BINARY_UNIX=$(BINARY_NAME)_unix

# Default environment
GO_ENV?=development

# Database defaults for development
DB_HOST?=localhost
DB_PORT?=5432
DB_USERNAME?=postgres
DB_PASSWORD?=password
DB_DATABASE?=universal_service_dev

# Service defaults
SERVICE_NAME?=universal-service
LOG_LEVEL?=debug

.PHONY: help run build test clean lint docker dev prod local install deps tidy

# Default target
all: build

## Help - Show available commands
help:
	@echo "🚀 Universal Go Service Boilerplate - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  make run          - Run development server with hot reload"
	@echo "  make dev          - Run in development mode"
	@echo "  make local        - Run in local mode (minimal setup)"
	@echo "  make prod         - Run in production mode"
	@echo ""
	@echo "Building:"
	@echo "  make build        - Build binary for current platform"
	@echo "  make build-linux  - Build binary for Linux"
	@echo "  make build-all    - Build binaries for all platforms"
	@echo ""
	@echo "Testing & Quality:"
	@echo "  make test              - Run all tests"
	@echo "  make test-unit         - Run unit tests only"
	@echo "  make test-integration  - Run integration tests only"
	@echo "  make test-cover        - Run tests with coverage report"
	@echo "  make test-cover-threshold - Run tests with 70% coverage requirement"
	@echo "  make test-domain       - Run domain layer tests"
	@echo "  make test-repository   - Run repository layer tests"
	@echo "  make test-usecase      - Run usecase layer tests"
	@echo "  make test-handler      - Run handler layer tests"
	@echo "  make test-race         - Run race condition tests"
	@echo "  make test-bench        - Run benchmark tests"
	@echo "  make test-clean        - Clean test artifacts"
	@echo "  make lint              - Run linter (requires golangci-lint)"
	@echo "  make fmt               - Format code"
	@echo "  make vet               - Run go vet"
	@echo ""
	@echo "Dependencies:"
	@echo "  make install      - Install development tools"
	@echo "  make deps         - Download dependencies"
	@echo "  make tidy         - Clean up dependencies"
	@echo ""
	@echo "Docker:"
	@echo "  make docker       - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make db-up        - Start local PostgreSQL with Docker"
	@echo "  make db-down      - Stop local PostgreSQL"
	@echo ""

## Development Commands

# Run development server (default)
run: deps
	@echo "🚀 Running Universal Go Service in $(GO_ENV) mode..."
	@GO_ENV=$(GO_ENV) \
	 SERVICE_NAME=$(SERVICE_NAME) \
	 LOG_LEVEL=$(LOG_LEVEL) \
	 $(GOCMD) run cmd/server/main.go

# Development mode
dev:
	@$(MAKE) run GO_ENV=development LOG_LEVEL=debug

# Local mode (minimal setup)
local:
	@$(MAKE) run GO_ENV=local DB_HOST=localhost DB_DATABASE=local_test

# Production mode
prod:
	@$(MAKE) run GO_ENV=production LOG_LEVEL=info

## Building Commands

# Build for current platform
build: deps
	@echo "🔨 Building $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BINARY_NAME) -v cmd/server/main.go
	@echo "✅ Built $(BINARY_NAME)"

# Build for Linux
build-linux: deps
	@echo "🔨 Building $(BINARY_NAME) for Linux..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v cmd/server/main.go
	@echo "✅ Built $(BINARY_UNIX)"

# Build for all platforms
build-all: deps
	@echo "🔨 Building for all platforms..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64 cmd/server/main.go
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64 cmd/server/main.go
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)-darwin-arm64 cmd/server/main.go
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows-amd64.exe cmd/server/main.go
	@echo "✅ Built binaries for all platforms"

## Testing Commands

# Run all tests
test: deps
	@echo "🧪 Running tests..."
	@$(GOTEST) -v ./...

# Run tests with verbose output
test-verbose: deps
	@echo "🧪 Running tests (verbose)..."
	@$(GOTEST) -v -race -timeout 30s ./...

# Run tests with coverage
test-cover: deps
	@echo "🧪 Running tests with coverage..."
	@$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Run unit tests only (exclude integration)
test-unit: deps
	@echo "🧪 Running unit tests..."
	@$(GOTEST) -v -short ./internal/domain/... ./internal/repository/... ./internal/usecase/... ./internal/handler/...

# Run integration tests only
test-integration: deps
	@echo "🧪 Running integration tests..."
	@$(GOTEST) -v ./testing/integration/...

# Run tests by layer
test-domain: deps
	@echo "🧪 Running domain layer tests..."
	@$(GOTEST) -v ./internal/domain/...

test-repository: deps
	@echo "🧪 Running repository layer tests..."
	@$(GOTEST) -v ./internal/repository/...

test-usecase: deps
	@echo "🧪 Running usecase layer tests..."
	@$(GOTEST) -v ./internal/usecase/...

test-handler: deps
	@echo "🧪 Running handler layer tests..."
	@$(GOTEST) -v ./internal/handler/...

# Run race condition tests
test-race: deps
	@echo "🧪 Running race condition tests..."
	@if [ -f "./scripts/run_race_tests.sh" ]; then \
		chmod +x ./scripts/run_race_tests.sh && ./scripts/run_race_tests.sh; \
	else \
		echo "⚠️  Race condition test script not found"; \
	fi

# Generate test coverage by layer
test-cover-detailed: deps
	@echo "🧪 Running detailed coverage analysis..."
	@mkdir -p coverage
	@echo "📋 Domain Layer Coverage:"
	@$(GOTEST) -v -race -coverprofile=coverage/domain.out ./internal/domain/... || echo "⚠️  Domain tests failed"
	@echo "📋 Repository Layer Coverage:"
	@$(GOTEST) -v -race -coverprofile=coverage/repository.out ./internal/repository/... || echo "⚠️  Repository tests failed"
	@echo "📋 UseCase Layer Coverage:"
	@$(GOTEST) -v -race -coverprofile=coverage/usecase.out ./internal/usecase/... || echo "⚠️  UseCase tests failed"
	@echo "📋 Handler Layer Coverage:"
	@$(GOTEST) -v -race -coverprofile=coverage/handler.out ./internal/handler/... || echo "⚠️  Handler tests failed"
	@echo "📊 Layer-specific coverage reports generated in coverage/"
	@echo "📈 Coverage Summary:"
	@ls -la coverage/ 2>/dev/null || echo "No coverage files generated"

# Test with minimum coverage threshold
test-cover-threshold: deps
	@echo "🧪 Running tests with coverage threshold..."
	@$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -func=coverage.out | grep "total:" | awk '{print "Coverage: " $$3}' | tee coverage.txt
	@if [ -f coverage.txt ]; then \
		COVERAGE=$$(grep -o '[0-9.]*' coverage.txt | head -1); \
		if [ "$$(echo "$$COVERAGE < 70" | bc)" -eq 1 ]; then \
			echo "❌ Coverage $$COVERAGE% is below 70% threshold"; \
			exit 1; \
		else \
			echo "✅ Coverage $$COVERAGE% meets threshold"; \
		fi; \
	fi

# Run benchmark tests
test-bench: deps
	@echo "🏃 Running benchmark tests..."
	@$(GOTEST) -bench=. -benchmem ./...

# Clean test artifacts
test-clean:
	@echo "🧹 Cleaning test artifacts..."
	@rm -f coverage.out coverage.html coverage.txt
	@rm -rf coverage/
	@echo "✅ Test artifacts cleaned"

# Format code
fmt:
	@echo "📝 Formatting code..."
	@$(GOCMD) fmt ./...

# Run go vet
vet:
	@echo "🔍 Running go vet..."
	@$(GOCMD) vet ./...

# Run linter (requires golangci-lint)
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Run 'make install' to install it."; \
	fi

## Dependency Commands

# Install development tools
install:
	@echo "📦 Installing development tools..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2; \
	fi
	@echo "✅ Development tools installed"

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	@$(GOMOD) download

# Clean up dependencies
tidy:
	@echo "🧹 Tidying dependencies..."
	@$(GOMOD) tidy
	@$(GOMOD) verify

## Docker Commands

# Build Docker image
docker:
	@echo "🐳 Building Docker image..."
	@docker build -t $(SERVICE_NAME):latest .
	@echo "✅ Docker image built: $(SERVICE_NAME):latest"

# Run Docker container
docker-run: docker
	@echo "🐳 Running Docker container..."
	@docker run --rm -p 8080:8080 \
		-e GO_ENV=production \
		-e DB_HOST=$(DB_HOST) \
		-e DB_USERNAME=$(DB_USERNAME) \
		-e DB_PASSWORD=$(DB_PASSWORD) \
		-e DB_DATABASE=$(DB_DATABASE) \
		$(SERVICE_NAME):latest

## Database Commands (for local development)

# Start local PostgreSQL with Docker
db-up:
	@echo "🐘 Starting local PostgreSQL..."
	@docker run --name postgres-dev -d \
		-e POSTGRES_USER=$(DB_USERNAME) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_DATABASE) \
		-p $(DB_PORT):5432 \
		postgres:15-alpine || true
	@echo "✅ PostgreSQL started on port $(DB_PORT)"
	@echo "📊 Connection: postgresql://$(DB_USERNAME):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_DATABASE)"

# Stop local PostgreSQL
db-down:
	@echo "🐘 Stopping local PostgreSQL..."
	@docker stop postgres-dev || true
	@docker rm postgres-dev || true
	@echo "✅ PostgreSQL stopped"

## Utility Commands

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_UNIX)
	@rm -f $(BINARY_NAME)-*
	@rm -f coverage.out coverage.html
	@echo "✅ Cleaned"

# Quick development setup
setup: install deps db-up
	@echo "✅ Development environment ready!"
	@echo "🚀 Run 'make dev' to start the server"

# Production build and run
production-run: build
	@echo "🚀 Running production build..."
	@GO_ENV=production \
	 LOG_LEVEL=info \
	 ./$(BINARY_NAME)