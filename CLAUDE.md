# CLAUDE.md

## Project Overview

This is **firecrest**, a Go-based event management web application. The system handles organizations, events, races, and user registrations with authentication support.

## Technology Stack

- **Language**: Go 1.25.6
- **Web Framework**: Native `net/http` with middleware via `justinas/alice`
- **Database**: PostgreSQL 16 (via `pgx/v5` driver and connection pooling)
- **Templating**: `templ` v0.3.977 (type-safe Go templating)
- **Styling**: Tailwind CSS 4.1.18 with templUI-style theme system
- **Hot Reload**: Air (configured via `.air.toml`)
- **Session Management**: `alexedwards/scs` with PostgreSQL store
- **Password Hashing**: `golang.org/x/crypto/bcrypt`

## Quick Start Commands

### Development Server
```bash
# Start PostgreSQL database
docker-compose up -d

# Run development server with hot reload (watches .go, .templ, and .css files)
air

# The server will be available at http://localhost:8080
# Build errors are logged to build-errors.log
```

### Database Setup
```bash
# Start PostgreSQL
docker-compose up -d

# Stop database
docker-compose down

# Stop and remove volumes (full reset)
docker-compose down -v

# Check database health
docker-compose ps
```

### Frontend/CSS Development
```bash
# Build CSS once
npm run css:build

# Watch CSS for changes (use in separate terminal)
npm run css:watch
```

### Code Generation
```bash
# Generate Go code from SQL queries (after editing query.sql)
sqlc generate

# Generate templ templates (after editing .templ files)
~/go/bin/templ generate
# or if templ is in PATH:
templ generate
```

## Project Structure

```
/cmd/web/          - Web application entry point
  main.go          - Server setup and initialization
  handlers.go      - HTTP request handlers
  routes.go        - Route definitions
  middleware.go    - HTTP middleware
  helpers.go       - Helper functions
/internal/         - Internal packages (not importable by other projects)
/db/               - Database related files
/tutorial/         - Generated database query code (sqlc)
/ui/               - UI templates and assets
  /static/         - Static assets (CSS, images)
    input.css      - Tailwind CSS source file
    main.css       - Compiled Tailwind output (generated)
  /templates/      - Templ template files
/docs/             - Documentation
schema.sql         - Database schema definitions
query.sql          - SQL queries for sqlc generation
sqlc.yaml          - sqlc configuration
docker-compose.yml - Docker services configuration
.air.toml          - Air hot reload configuration
.golangci.yml      - Linting rules configuration
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

### Environment Variables

Create a `.env` file in the project root (see `.env.example`):

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=firecrest
DB_SSLMODE=disable

# Session Configuration
SESSION_SECRET=your-secret-key-change-this-in-production
SESSION_LIFETIME_HOURS=12
SESSION_LIFETIME_REMEMBER_HOURS=720  # 30 days

# Security Configuration
ACCOUNT_LOCKOUT_MINUTES=15
MAX_LOGIN_ATTEMPTS=5
```

### Database Management

- Schema is defined in `schema.sql`
- Queries are defined in `query.sql`
- Database code is generated using `sqlc` (see `sqlc.yaml`)
- Generated code goes into `/tutorial/` directory
- After editing queries: `sqlc generate`

**Database Migration Workflow:**
1. Edit `schema.sql` for schema changes
2. Edit `query.sql` for new queries
3. Run `sqlc generate` to update Go code
4. Apply schema changes to database manually or via migration tool

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

# With timeout for large codebases
golangci-lint run --timeout=5m ./...
```

**CI/CD:** Pull requests are automatically checked by GitHub Actions and will be blocked from merging if linting fails. Ensure your code passes all linting checks before pushing.

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests in a specific package
go test ./cmd/web

# Run a specific test
go test -run TestFunctionName ./...
```

## Frontend Development

### Tailwind CSS 4.1 with templUI Theme

The project uses Tailwind CSS 4.1 with a templUI-style semantic color system and theme configuration.

**Theme Configuration:** Located in `ui/static/input.css`

**Semantic Colors Available:**
- `background` / `foreground` - Base colors
- `card` / `card-foreground` - Card components
- `primary` / `primary-foreground` - Primary actions (green theme)
- `secondary` / `secondary-foreground` - Secondary elements
- `muted` / `muted-foreground` - Muted/disabled states
- `accent` / `accent-foreground` - Accent highlights
- `destructive` / `destructive-foreground` - Destructive actions
- `border` / `input` / `ring` - Form and border colors

**Border Radius Variables:**
- `--radius-sm` (0.25rem), `--radius-md` (0.375rem), `--radius-lg` (0.5rem), `--radius-xl` (0.75rem)

**Dark Mode:** Automatically supported via `prefers-color-scheme: dark`

**View Transitions:** Page transitions with 0.3s fade effect are configured and ready to use.

### CSS Development Workflow

1. Edit `ui/static/input.css` for global styles and theme changes
2. Use Tailwind utility classes in `.templ` files
3. Run `npm run css:build` or `npm run css:watch`
4. Output goes to `ui/static/main.css` (this file is generated, don't edit directly)

**CSS Content Paths:**
- Tailwind scans: `../templates/**/*.templ`
- Only classes used in templates will be included in the output

## Build Process

The Air hot reload system automatically runs this build chain on file changes:

```bash
# 1. Compile Tailwind CSS
npx @tailwindcss/cli -i ui/static/input.css -o ui/static/main.css

# 2. Generate templ templates
~/go/bin/templ generate

# 3. Build Go application
go build -o ./tmp/main ./cmd/web
```

**Watched File Extensions:** `.go`, `.templ`, `.css`
**Excluded Directories:** `assets`, `tmp`, `vendor`, `testdata`, `docs`, `firecrest`, `db`, `node_modules`
**Excluded Files:** `*_test.go`, `*_templ.go` (generated files)

## Coding Conventions

### Go Code Standards

1. **Error Handling**: Always check and handle errors explicitly
2. **Logging**: Use structured logging with `slog` package
3. **Database Queries**: Use sqlc-generated type-safe queries, never write raw SQL in handlers
4. **Context**: Pass `context.Context` for database operations and HTTP handlers
5. **Soft Deletes**: Use `deleted_at` fields, never hard delete records
6. **Validation**: Validate user input at handler level before database operations

### Templ Template Conventions

1. **Components**: Create reusable components for common UI patterns
2. **Type Safety**: Leverage templ's type-safe parameters
3. **Naming**: Use descriptive component names (e.g., `LoginForm`, `EventCard`)
4. **Styling**: Use Tailwind utility classes, prefer semantic color names

### Security Best Practices

1. **Authentication**: Use session-based auth with CSRF protection
2. **Passwords**: Always use bcrypt for password hashing
3. **Input Validation**: Sanitize all user input
4. **SQL Injection**: Use parameterized queries (sqlc handles this)
5. **XSS Prevention**: Templ auto-escapes output by default
6. **Session Security**: Configure secure cookies in production

### Database Conventions

1. **Naming**: Use `snake_case` for table and column names
2. **Primary Keys**: Use `SERIAL` or `BIGSERIAL` for `id` columns
3. **Timestamps**: Include `created_at`, `updated_at` (managed by triggers)
4. **Soft Deletes**: Include `deleted_at TIMESTAMPTZ` for all user-facing tables
5. **Foreign Keys**: Always define foreign key constraints
6. **Indexes**: Add indexes for frequently queried columns

## Key Features

- Multi-tenant organization support
- Event and race management
- User registration and authentication
- Role-based access control (entrant, organizer, admin)
- Soft deletes for data retention
- Social authentication support

## Key Dependencies

### Go Packages
- `github.com/a-h/templ` - Type-safe HTML templating
- `github.com/jackc/pgx/v5` - PostgreSQL driver and toolkit
- `github.com/alexedwards/scs/v2` - Session management
- `github.com/justinas/alice` - HTTP middleware chaining
- `golang.org/x/crypto` - Password hashing (bcrypt)
- `github.com/joho/godotenv` - Environment variable loading

### Frontend Dependencies
- `@tailwindcss/cli` v4.1.18 - Standalone Tailwind CSS compiler
- `tailwindcss` v4.1.18 - Tailwind CSS framework

### Development Tools
- `sqlc` - SQL to Go code generator
- `templ` - Template compiler (installed as Go tool)
- `air` - Hot reload for Go applications
- `golangci-lint` v2.8.0 - Comprehensive Go linter

## Architecture Patterns

### Request Flow
1. HTTP Request → Router (`routes.go`)
2. Middleware Chain (`middleware.go`) - auth, CSRF, logging, etc.
3. Handler (`handlers.go`) - business logic
4. Database Layer (`tutorial/*.go`) - sqlc-generated queries
5. Template Rendering (`ui/templates/*.templ`) - templ components
6. HTTP Response

### Middleware Stack (via alice)
Common middleware applied to routes:
- Session management
- CSRF protection
- Authentication checks
- Request logging
- Panic recovery

### Database Access Pattern
```go
// Always use context-aware queries
user, err := queries.GetUserByEmail(ctx, email)

// Use transactions for multi-step operations
tx, err := db.Begin(ctx)
defer tx.Rollback(ctx)
// ... perform operations
tx.Commit(ctx)
```

### Session Management
- Sessions stored in PostgreSQL via `pgxstore`
- Configurable lifetime (12 hours default, 30 days with "remember me")
- Automatic session renewal
- Secure cookie settings required in production

## Common Tasks

### Adding a New Route
1. Define handler in `cmd/web/handlers.go`
2. Add route in `cmd/web/routes.go`
3. Create templ template if needed in `ui/templates/`
4. Add middleware chain if authentication required

### Adding a New Database Query
1. Write SQL query in `query.sql`
2. Run `sqlc generate`
3. Use generated functions in `tutorial/` package
4. Call from handlers with proper error handling

### Creating a New Page
1. Create templ file: `ui/templates/mypage.templ`
2. Add handler function in `cmd/web/handlers.go`
3. Register route in `cmd/web/routes.go`
4. Add CSS styles using Tailwind utilities
5. Test with `air` running

### Adding a New Database Table
1. Add CREATE TABLE statement to `schema.sql`
2. Include required columns: `id`, `created_at`, `updated_at`, `deleted_at`
3. Add triggers for automatic timestamp management
4. Run migration (manual or via migration tool)
5. Add queries to `query.sql`
6. Run `sqlc generate`

## Troubleshooting

### Build Errors
- Check `build-errors.log` for Air build failures
- Ensure templ is installed: `go install github.com/a-h/templ/cmd/templ@latest`
- Verify Go version: `go version` (should be 1.25.6)

### Database Connection Issues
- Verify PostgreSQL is running: `docker-compose ps`
- Check connection string in `.env` file
- Ensure database exists: `docker-compose exec postgres psql -U postgres -l`

### CSS Not Updating
- Check if `npm run css:watch` is running
- Verify file paths in `ui/static/input.css` `@source` directive
- Clear browser cache
- Check if `ui/static/main.css` is being regenerated

### Templ Templates Not Updating
- Ensure Air is watching `.templ` files (check `.air.toml`)
- Verify templ is in PATH or using correct path: `~/go/bin/templ`
- Check for syntax errors in `.templ` files

### Linting Failures
- Run locally: `golangci-lint run ./...`
- Check `.golangci.yml` for enabled linters
- Common issues: unused variables, error handling, imports ordering
- Fix issues before pushing to avoid CI failures

## Development Notes

- The application uses structured logging via `slog`
- Database queries are type-safe via sqlc code generation
- Connection pooling is handled by `pgxpool`
- Templates are compiled with `templ` for type safety
- All timestamps are stored in UTC
- Cookie-based sessions with CSRF protection
- View transitions provide smooth page navigation
- Hot reload watches Go, templ, and CSS files simultaneously
