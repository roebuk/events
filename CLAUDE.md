# CLAUDE.md

## Project Overview

This is **firecrest**, a Go-based event management web application. The system handles organizations, events, races, and user registrations with authentication support.

## Technology Stack

- **Language**: Go 1.25.4
- **Web Framework**: Native `net/http` with middleware via `justinas/alice`
- **Database**: PostgreSQL (via `pgx/v5` driver and connection pooling)
- **Templating**: `templ` (type-safe Go templating)
- **Hot Reload**: Air (configured via `.air.toml`)

## Project Structure

```
/cmd/web/          - Web application entry point
  main.go          - Server setup and initialization
  handlers.go      - HTTP request handlers
  routes.go        - Route definitions
  middleware.go    - HTTP middleware
  helpers.go       - Helper functions
/db/               - Database related files
/tutorial/         - Generated database query code (sqlc)
/ui/               - UI templates and assets
/docs/             - Documentation
schema.sql         - Database schema definitions
query.sql          - SQL queries for sqlc generation
sqlc.yaml          - sqlc configuration
docker-compose.yml - Docker services configuration
```

## Database Schema

The application uses PostgreSQL with the following main entities:

- **users**: User accounts with roles (entrant, organizer, admin)
- **organisations**: Event organizing bodies
- **events**: Events hosted by organizations
- **races**: Individual races within events
- **organisation_users**: Many-to-many relationship between orgs and users
- **auth_credentials**: Password-based authentication
- **social_accounts**: OAuth authentication (Google, Apple)

All tables include soft delete support (`deleted_at`) and automatic timestamp management.

## Development

### Running the Application

The server runs on port `8080` and connects to PostgreSQL at `127.0.0.1:5432/firecrest`.

Default credentials: `postgres:postgres`

### Database Management

- Schema is defined in `schema.sql`
- Queries are defined in `query.sql`
- Database code is generated using `sqlc` (see `sqlc.yaml`)
- Run `sqlc generate` to regenerate database code after query changes

### Environment

See `.env.example` for required environment variables.

### Code Quality and Linting

All code must adhere to the linting rules defined in `.golangci.yml`. The project uses `golangci-lint` with a comprehensive set of linters covering:

- Error handling and correctness
- Code simplification and best practices
- Security vulnerabilities
- Code formatting and style
- Performance optimizations

**Running linting locally:**
```bash
golangci-lint run ./...
```

**CI/CD:** Pull requests are automatically checked by GitHub Actions and will be blocked from merging if linting fails. Ensure your code passes all linting checks before pushing.

## Key Features

- Multi-tenant organization support
- Event and race management
- User registration and authentication
- Role-based access control (entrant, organizer, admin)
- Soft deletes for data retention
- Social authentication support

## Development Notes

- The application uses structured logging via `slog`
- Database queries are type-safe via sqlc code generation
- Connection pooling is handled by `pgxpool`
- Templates are compiled with `templ` for type safety
