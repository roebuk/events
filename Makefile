.PHONY: help test test-coverage test-race lint lint-fix fmt vet build run clean sqlc-generate templ-generate migrate-up migrate-down docker-up docker-down install-tools

# Default target
help:
	@echo "Available targets:"
	@echo "  make test              - Run all tests"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make test-race         - Run tests with race detector"
	@echo "  make lint              - Run golangci-lint"
	@echo "  make lint-fix          - Run golangci-lint with auto-fix"
	@echo "  make fmt               - Format code with gofumpt"
	@echo "  make vet               - Run go vet"
	@echo "  make build             - Build the application"
	@echo "  make run               - Run the application with hot reload (air)"
	@echo "  make clean             - Clean build artifacts and cache"
	@echo "  make sqlc-generate     - Generate sqlc code from query.sql"
	@echo "  make templ-generate    - Generate templ templates"
	@echo "  make migrate-up        - Run database migrations up"
	@echo "  make migrate-down      - Run database migrations down"
	@echo "  make docker-up         - Start PostgreSQL with docker-compose"
	@echo "  make docker-down       - Stop PostgreSQL docker-compose"
	@echo "  make install-tools     - Install development tools"

# Test targets
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	go test -v -race ./...

# Linting and formatting
lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

fmt:
	gofumpt -l -w .

vet:
	go vet ./...

# Build targets
build:
	go build -o bin/server ./cmd/web

run:
	air

# Clean targets
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache -testcache

# Code generation
sqlc-generate:
	sqlc generate

templ-generate:
	templ generate

# Database migrations
migrate-up:
	@echo "Running migrations up..."
	@echo "Note: Ensure PostgreSQL is running and DATABASE_URL is set"
	# Add migration tool command here when ready
	# migrate -path db/migrations -database "$(DATABASE_URL)" up

migrate-down:
	@echo "Running migrations down..."
	@echo "Note: Ensure PostgreSQL is running and DATABASE_URL is set"
	# Add migration tool command here when ready
	# migrate -path db/migrations -database "$(DATABASE_URL)" down

# Docker targets
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Tool installation
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/air-verse/air@latest
	@echo "Tools installed successfully!"

# CI target (runs in CI/CD pipeline)
ci: lint test-race
	@echo "CI checks passed!"
