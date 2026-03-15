.PHONY: build run dev test test-coverage lint clean docker-build docker-up docker-down migrate-up migrate-down migrate-version

# Variables
APP_NAME    := agent-core-service
BUILD_DIR   := ./build
MAIN_PATH   := ./cmd/api
CONFIG_PATH := config/config.yaml
DOCKER_COMPOSE := docker-compose
DB_URL      := postgres://postgres:postgres@localhost:5434/agent_core_db?sslmode=disable

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@CONFIG_PATH=$(CONFIG_PATH) $(BUILD_DIR)/$(APP_NAME)

# Run with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	@echo "Running $(APP_NAME) in development mode..."
	@CONFIG_PATH=$(CONFIG_PATH) air -c .air.toml || CONFIG_PATH=$(CONFIG_PATH) go run $(MAIN_PATH)

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -count=1 ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

docker-up:
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE) up -d

docker-down:
	@echo "Stopping Docker containers..."
	@$(DOCKER_COMPOSE) down

docker-logs:
	@$(DOCKER_COMPOSE) logs -f app

docker-reset:
	@echo "Resetting Docker environment (volumes will be wiped)..."
	@$(DOCKER_COMPOSE) down -v
	@$(DOCKER_COMPOSE) up -d

# Database migrations
migrate-up:
	@echo "Running migrations..."
	@migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Rolling back last migration..."
	@migrate -path migrations -database "$(DB_URL)" down 1

migrate-version:
	@migrate -path migrations -database "$(DB_URL)" version

migrate-create:
	@echo "Usage: make migrate-create NAME=<migration_name>"
	@migrate create -ext sql -dir migrations -seq $(NAME)

# Go module management
mod-tidy:
	@go mod tidy

mod-download:
	@go mod download

# Generate mocks (requires mockgen)
mock:
	@echo "Generating mocks..."
	@mockgen -source=internal/domain/repository.go -destination=internal/mocks/repository_mock.go -package=mocks

# Help
help:
	@echo "Available targets:"
	@echo "  build            - Build the application"
	@echo "  run              - Build and run the application"
	@echo "  dev              - Run with hot reload"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  lint             - Run linter"
	@echo "  clean            - Clean build artifacts"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-up        - Start Docker containers"
	@echo "  docker-down      - Stop Docker containers"
	@echo "  docker-logs      - View Docker logs"
	@echo "  docker-reset     - Wipe volumes and restart"
	@echo "  migrate-up       - Run database migrations"
	@echo "  migrate-down     - Rollback last migration"
	@echo "  migrate-version  - Show migration version"
	@echo "  migrate-create   - Create new migration (NAME=<name>)"
	@echo "  mod-tidy         - Tidy go modules"
	@echo "  mod-download     - Download go modules"
	@echo "  mock             - Generate mocks"
