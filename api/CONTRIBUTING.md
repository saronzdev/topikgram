# Contributing to NotikGram

## Prerequisites

- Go 1.26+
- PostgreSQL 14+
- A running PostgreSQL instance with a database created

## Local Development

1. Fork and clone the repo
2. Copy `.env.example` to `.env` and configure your database connection
3. Generate a JWT secret: `go run cmd/keygen/main.go`
4. Run the server: `go run main.go` (migrations run automatically)

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use short variable names for loop indices and local scopes (`i`, `db`, `w`, `r`)
- Error handling: always check errors, return early, use sentinel errors from `domain/errors.go`
- No comments unless explicitly requested
- Import grouping: stdlib, external packages, then internal packages (separated by blank line)
- JSON tags use `snake_case`
- Handler methods are PascalCase, unexported helpers are camelCase

## Package Conventions

### Adding a New Endpoint

1. Define input/output models in `internal/domain/models.go`
2. Add SQL queries in the appropriate `internal/store/*.go` file
3. Create the handler method in `internal/transport/*.go`
4. Register the route in the handler's `RegisterRoutes()` method
5. Use `AuthMiddleware()` wrapper for authenticated endpoints

Example route registration:
```go
h.mux.HandleFunc("GET /"+prefix+"/your-resource", h.List)
h.mux.Handle("POST /"+prefix+"/your-resource", AuthMiddleware(http.HandlerFunc(h.Create)))
```

### Adding a New Migration

1. Create a new file in `internal/migrations/` with the next sequential number:
   ```
   008_your_description.sql
   ```
2. Write the SQL statement (use `CREATE TABLE IF NOT EXISTS` or `ALTER TABLE` as needed)
3. The file is automatically embedded and applied on next startup

### Adding a New Store

1. Create a new file in `internal/store/`
2. Define a struct with a `*Pool` field
3. Add a constructor: `func NewXxxStore(db *Pool) *XxxStore`
4. Methods should accept `context.Context` as the first parameter
5. Use `pgx` query methods (`QueryRow`, `Query`, `Exec`) with parameterized queries

### Input Validation

Use functions from `internal/validator/`:
```go
validator.Email(email)       // returns bool
validator.Username(username) // returns bool
validator.Password(password) // returns bool
validator.MaxLength(s, n)    // returns bool
validator.MinLength(s, n)    // returns bool
```

### Error Mapping

Use `mapStoreError()` to convert domain errors to user-facing messages, then call the appropriate HTTP helper:
```go
if err == domain.ErrNotFound {
    notFound(w, "resource not found")
    return
}
internalError(w, mapStoreError(err))
```

## Testing

Run the existing tests:
```bash
go test ./internal/transport/ -v
```

### Writing Tests

- Tests are in `_test.go` files alongside the code they test
- Use table-driven tests for input validation
- Tests in `transport/` are unit-level and do not require a database
- Use `httptest.NewRecorder()` and `httptest.NewRequest()` for handler tests
- Set `JWT_SECRET` env var in `TestMain` for JWT-related tests

Example:
```go
func TestMyHandler(t *testing.T) {
    tests := []struct {
        name       string
        body       string
        wantStatus int
    }{
        {"valid input", `{"name":"test"}`, http.StatusOK},
        {"empty body", `{}`, http.StatusBadRequest},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("POST", "/path", strings.NewReader(tt.body))
            req.Header.Set("Content-Type", "application/json")
            rec := httptest.NewRecorder()

            handler(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("got %d, want %d", rec.Code, tt.wantStatus)
            }
        })
    }
}
```

## Commit Messages

Use clear, concise commit messages:
- `add user follow endpoint`
- `fix cursor pagination offset`
- `add migration for posts.modified_at`

## Pull Requests

1. Create a feature branch from `main`
2. Make your changes following the conventions above
3. Run `go vet ./...` and `go test ./internal/transport/ -v`
4. Open a PR with a clear description of what changed and why
