# Project Overview: firecrest-go

## Purpose
A Go-based event management web application that handles organizations, events, races, and user registrations with authentication support.

## Tech Stack
- **Language**: Go 1.25.4
- **Web Framework**: Native `net/http` with middleware via `justinas/alice`
- **Database**: PostgreSQL (via `pgx/v5` driver with connection pooling)
- **Templating**: `templ` (type-safe Go templating)
- **Hot Reload**: Air (configured via `.air.toml`)
- **Linting**: golangci-lint with comprehensive ruleset

## Project Structure
```
/cmd/web/          - Web application entry point
  main.go          - Server setup and initialization
  handlers.go      - HTTP request handlers (methods on *application)
  routes.go        - Route definitions
  middleware.go    - HTTP middleware
  helpers.go       - Helper functions
/db/               - Database package (sqlc generated)
/ui/               - UI templates and static assets
  /templates/      - templ template files
  /static/         - Static files (CSS, JS, images)
/docs/             - Documentation
schema.sql         - Database schema definitions
query.sql          - SQL queries for sqlc generation
sqlc.yaml          - sqlc configuration
```

## Database
- PostgreSQL at `127.0.0.1:5432/firecrest`
- Default credentials: `postgres:postgres`
- Schema defined in `schema.sql`
- Queries in `query.sql`, code generated via sqlc

## Key Patterns
- Handlers are methods on `*application` struct
- Routes use `mux.Handle()` or `mux.HandleFunc()` with alice middleware chains
- Templates rendered via `app.render(ctx, w, status, template)`
- Error handling via `app.serverError()`, `app.clientError()`, `app.notFound()`
