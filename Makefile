.PHONY: help build run dev clean test test-coverage lint tidy deps mock \
        graceful-shutdown graceful-upgrade fast-shutdown \
        docker-build docker-up docker-down docker-logs docker-reset \
        migrate-up migrate-down migrate-version \
        quick-start

# Build variables
APP_NAME := agent-core-service
BUILD_DIR := ./bin
MAIN_FILE := ./cmd/api/main.go
CONFIG_PATH := config/config.yaml
DOCKER_IMAGE := $(APP_NAME):latest
PID_FILE := /tmp/$(APP_NAME).pid

# Go variables
GO := go
GOFLAGS := -v
LDFLAGS := -s -w

# DB variables
DB_URL := postgres://postgres:postgres@localhost:5434/agent_core_db?sslmode=disable

# Default target
help:
	@echo "Agent Core Service - Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make build          - Build the application binary"
	@echo "  make run            - Build and run the application"
	@echo "  make dev            - Run with hot reload (requires air)"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make lint           - Run linter (requires golangci-lint)"
	@echo "  make tidy           - Tidy go modules"
	@echo "  make deps           - Download dependencies"
	@echo "  make mock           - Generate mocks (requires mockgen)"
	@echo ""
	@echo "Graceful Operations:"
	@echo "  make graceful-shutdown  - Send SIGTERM for graceful shutdown"
	@echo "  make graceful-upgrade   - Send SIGQUIT for graceful upgrade"
	@echo "  make fast-shutdown      - Send SIGINT for fast shutdown"
	@echo ""
	@echo "Docker (Local Dev):"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start all services"
	@echo "  make docker-down    - Stop all services"
	@echo "  make docker-logs    - View application logs"
	@echo "  make docker-reset   - Reset all (remove volumes, rebuild)"
	@echo ""
	@echo "Database Migrations:"
	@echo "  make migrate-up      - Run all pending migrations"
	@echo "  make migrate-down    - Rollback last migration"
	@echo "  make migrate-version - Show current migration version"
	@echo ""
	@echo "Quick Start:"
	@echo "  make quick-start    - Reset and start everything fresh"

# Development
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

run: build
	@echo "Running $(APP_NAME)..."
	$(BUILD_DIR)/$(APP_NAME) -config=$(CONFIG_PATH)

dev:
	air

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

test:
	$(GO) test ./... -v -cover

test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	golangci-lint run ./...

tidy:
	$(GO) mod tidy

deps:
	$(GO) mod download

# Generate mocks (requires mockgen)
mock:
	@echo "Generating mocks..."
	mockgen -source=internal/domain/repository.go -destination=internal/mocks/repository_mock.go -package=mocks

# Graceful Operations
graceful-shutdown:
	@echo "Sending SIGTERM for graceful shutdown..."
	@if [ -f $(PID_FILE) ]; then \
		kill -TERM $$(cat $(PID_FILE)); \
		echo "Graceful shutdown signal sent"; \
	else \
		pkill -TERM $(APP_NAME) && echo "Graceful shutdown signal sent" || echo "No running instance found"; \
	fi

graceful-upgrade:
	@echo "Sending SIGQUIT for graceful upgrade..."
	@if [ -f $(PID_FILE) ]; then \
		kill -QUIT $$(cat $(PID_FILE)); \
		echo "Graceful upgrade signal sent"; \
	else \
		pkill -QUIT $(APP_NAME) && echo "Graceful upgrade signal sent" || echo "No running instance found"; \
	fi

fast-shutdown:
	@echo "Sending SIGINT for fast shutdown..."
	@if [ -f $(PID_FILE) ]; then \
		kill -INT $$(cat $(PID_FILE)); \
		echo "Fast shutdown signal sent"; \
	else \
		pkill -INT $(APP_NAME) && echo "Fast shutdown signal sent" || echo "No running instance found"; \
	fi

# Docker (Local Development)
docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app

docker-reset:
	docker compose down -v
	docker compose up -d --build

# Database Migrations
migrate-up:
	@echo "Running migrations..."
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Rolling back last migration..."
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

# Quick Start
quick-start: docker-reset
