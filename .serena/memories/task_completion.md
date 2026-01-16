# Task Completion Checklist

When completing a task, ensure the following:

## Before Committing

1. **Generate code** (if applicable):
   ```bash
   templ generate    # If .templ files were modified
   sqlc generate     # If query.sql was modified
   ```

2. **Run linter**:
   ```bash
   golangci-lint run ./...
   ```
   Fix any issues before proceeding.

3. **Run tests**:
   ```bash
   go test ./...
   ```

4. **Manual verification** (if applicable):
   - Start the server with `air`
   - Test the affected endpoints in browser or with curl

## Code Changes Checklist
- [ ] Handlers follow `*application` method pattern
- [ ] Errors handled appropriately (serverError, clientError, notFound)
- [ ] Context passed to database calls
- [ ] Routes registered in routes.go
- [ ] Templates use templ (not html/template)

## Database Changes
- [ ] Schema changes added to schema.sql
- [ ] Queries added to query.sql
- [ ] Run `sqlc generate` to regenerate db package
