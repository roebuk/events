.PHONY: help run-dev run-staging run-prod build test clean generate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

run-dev: ## Run application in development mode
	@echo "Starting in development mode..."
	@APP_ENV=development go run ./cmd/web

run-staging: ## Run application in staging mode
	@echo "Starting in staging mode..."
	@APP_ENV=staging go run ./cmd/web

run-prod: ## Run application in production mode
	@echo "Starting in production mode..."
	@APP_ENV=production go run ./cmd/web

build: ## Build the application binary
	@echo "Building application..."
	@go build -o bin/firecrest ./cmd/web

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

generate: ## Generate code (templ, sqlc, etc.)
	@echo "Generating code..."
	@templ generate
	@sqlc generate

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t firecrest:latest .

docker-run-dev: ## Run Docker container in development mode
	@docker run -p 8080:8080 --env-file .env.local -e APP_ENV=development firecrest:latest

docker-run-prod: ## Run Docker container in production mode
	@docker run -p 8080:8080 --env-file .env.production -e APP_ENV=production firecrest:latest

setup: ## Initial setup - copy example env files
	@echo "Setting up environment files..."
	@test -f .env.local || cp .env.example .env.local
	@test -f .env.staging || cp .env.example .env.staging
	@test -f .env.production || cp .env.example .env.production
	@echo "Environment files created. Please edit them with your configuration."

db-migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@migrate -path db/migrations -database "$$(cat .env.local | grep DB_HOST | cut -d'=' -f2- | xargs -I {} echo 'postgres://{}:5432/firecrest?sslmode=disable')" up

db-migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	@migrate -path db/migrations -database "$$(cat .env.local | grep DB_HOST | cut -d'=' -f2- | xargs -I {} echo 'postgres://{}:5432/firecrest?sslmode=disable')" down
