# Project Tools & Commands

This document describes the tools used in this project and how to run them.

## Go

The project is built with Go 1.25.4.

### Build

Compile the application:
```bash
go build ./cmd/web
```

Build for a specific OS/architecture:
```bash
GOOS=linux GOARCH=amd64 go build ./cmd/web
```

### Run

Run the application directly:
```bash
go run ./cmd/web
```

### Test

Run all tests:
```bash
go test ./...
```

Run tests for a specific package:
```bash
go test ./cmd/web
```

Run tests with verbose output:
```bash
go test -v ./...
```

### Dependencies

Install dependencies:
```bash
go mod download
```

Update dependencies:
```bash
go mod tidy
```

## Linting

The project uses **golangci-lint** for comprehensive Go code linting and analysis.

### Install

If not already installed:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Run

Lint all Go files:
```bash
golangci-lint run ./...
```

Lint a specific directory:
```bash
golangci-lint run ./cmd/web/...
```

Auto-fix issues where possible:
```bash
golangci-lint run --fix ./...
```

### Configuration

Linting rules are configured in `.golangci.yml`. Key settings:
- **Enabled linters**: errcheck, govet, staticcheck, gosimple, gocritic, revive, gofmt, goimports, gosec, and more
- **Excluded paths**: Generated files (`*_templ.go`, `*.sql.go`) and vendor directories
- **Timeout**: 5 minutes

The configuration enables strict checks including error handling, context usage, security analysis, and code style.

## Templ

The project uses **templ** for type-safe HTML templating.

### Install

If not already installed:
```bash
go install github.com/a-h/templ/cmd/templ@latest
```

### Generate

Generate Go code from templ files:
```bash
templ generate
```

Watch for changes and auto-generate:
```bash
templ generate --watch
```

### File Format

Templ files are located in `ui/templates/` with `.templ` extension:
- `ui/templates/*.templ` - Page templates
- `ui/templates/components/*.templ` - Reusable components
- `ui/templates/auth/*.templ` - Authentication templates

Generated files (`*_templ.go`) are automatically created and should not be edited manually.

## sqlc

The project uses **sqlc** for type-safe SQL code generation.

### Install

If not already installed:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Generate

Generate Go code from SQL queries:
```bash
sqlc generate
```

Watch for changes and auto-generate:
```bash
sqlc watch
```

### File Format

SQL queries are located in `db/` directory with `.sql` extension:
- `db/*.sql` - Query definitions
- `db/schema.sql` - Database schema

Generated files (`*.sql.go`) are automatically created and should not be edited manually.

Configuration is in `sqlc.yaml`.

## Database (PostgreSQL)

The project uses PostgreSQL for data storage.

### Connection

Default connection string:
```
postgres://postgres:postgres@127.0.0.1:5432/firecrest
```

To change the connection string, edit `cmd/web/main.go` line 29.

### Setup

Refer to `docs/postgres-docker.md` for Docker setup instructions.

### Running Migrations

Apply schema to the database:
```bash
psql -h 127.0.0.1 -U postgres -d firecrest -f db/schema.sql
```

## Development Workflow

### Before Committing

1. Run linter to catch issues:
   ```bash
   golangci-lint run ./...
   ```

2. Build to verify code compiles:
   ```bash
   go build ./cmd/web
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

### Adding New Features

1. **If adding database features:**
   - Write SQL queries in `db/*.sql`
   - Run `sqlc generate` to create Go types and functions

2. **If adding HTML templates:**
   - Create/edit `.templ` files in `ui/templates/`
   - Run `templ generate` to create Go code

3. **After changes:**
   - Run linter: `golangci-lint run ./...`
   - Build and test: `go build ./cmd/web && go test ./...`

## Useful Flags

### golangci-lint

- `--fix` - Auto-fix issues where possible
- `--out-format json` - Output as JSON for programmatic parsing
- `--timeout 10m` - Increase timeout for large projects
- `--verbose` - Show all checks being run

### go build

- `-o <name>` - Output filename
- `-ldflags "-X main.Version=1.0"` - Set build-time variables
- `-race` - Enable race detector (slower, detects concurrency issues)

### go test

- `-v` - Verbose output
- `-cover` - Show code coverage
- `-race` - Enable race detector
- `-timeout 10m` - Set test timeout
