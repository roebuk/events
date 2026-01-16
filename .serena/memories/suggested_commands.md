# Development Commands

## Running the Application
```bash
# Hot reload development server (recommended)
air

# Manual build and run
go build -o ./tmp/main ./cmd/web && ./tmp/main
```

## Code Generation
```bash
# Generate templ templates (required after .templ file changes)
templ generate

# Generate database code (required after query.sql changes)
sqlc generate
```

## Linting & Formatting
```bash
# Run linter
golangci-lint run ./...

# Format code
gofmt -w .
goimports -w .
```

## Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Database
```bash
# Start PostgreSQL via Docker
docker-compose up -d

# Connect to database
psql -h 127.0.0.1 -U postgres -d firecrest
```

## Windows-Specific Notes
- Use PowerShell or Git Bash for commands
- Path separator is `\` but Go tools accept `/`
- Use `go.exe` if `go` is not in PATH
