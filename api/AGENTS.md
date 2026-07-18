# AGENTS.md

Instructions for AI agents working on this codebase.

## Project Overview

NotikGram is a Go REST API (social media backend) using PostgreSQL, JWT auth, and the standard library `net/http.ServeMux` router. No web framework or ORM.

## Build & Run

```bash
# Run the server (includes migrations)
go run main.go

# Build binary
go build -o notikgram .

# Run migrations standalone
go run cmd/migrate/main.go -dir=up

# Generate JWT secret
go run cmd/keygen/main.go
```

## Test Commands

```bash
# Run all tests
go test ./internal/transport/ -v

# Run specific test
go test ./internal/transport/ -run TestGenerateToken -v

# Vet (no linter configured)
go vet ./...
```

No `Makefile`, `golangci-lint`, or CI config exists. Use `go vet` for static analysis.

## Project Structure

```
notikgram/
├── main.go                           # Entry point: wires everything together
├── cmd/
│   ├── keygen/main.go                # JWT key generator utility
│   └── migrate/main.go               # Standalone migration runner
├── internal/
│   ├── domain/
│   │   ├── models.go                 # All structs + TopicID enum (0-20)
│   │   └── errors.go                 # Sentinel errors (ErrNotFound, etc.)
│   ├── migrations/
│   │   ├── migrations.go             # Embed-based migration engine
│   │   └── *.sql                     # SQL migration files (sequential)
│   ├── ratelimit/
│   │   └── ratelimit.go              # Token bucket per IP
│   ├── store/
│   │   ├── db.go                     # pgxpool connection setup
│   │   ├── auth.go                   # Register, login (bcrypt)
│   │   ├── users.go                  # User CRUD, follow/unfollow
│   │   ├── posts.go                  # Post CRUD, like/unlike, save/unsave
│   │   ├── comments.go               # Comment create/list
│   │   └── interests.go              # Topic weight tracking
│   ├── transport/
│   │   ├── auth.go                   # Auth handlers (register, login, me)
│   │   ├── users.go                  # User handlers
│   │   ├── posts.go                  # Post handlers
│   │   ├── comments.go               # Comment handlers
│   │   ├── middlewares.go            # JWT, CORS, auth, body limit, timeout, logging
│   │   ├── helpers.go                # JSON response helpers, error mapping, health
│   │   ├── auth_test.go              # Auth validation tests
│   │   └── middlewares_test.go       # JWT/middleware tests
│   └── validator/
│       └── validator.go              # Email, username, password, max/min length
├── .env.example                      # Template environment config
└── go.mod                            # Module: notikgram
```

## Key Files

- **`main.go`** — Start here to understand bootstrap, route registration, middleware chain
- **`internal/domain/models.go`** — All data structures, TopicID constants
- **`internal/transport/helpers.go`** — JSON response helpers, error mapping
- **`internal/store/db.go`** — Database pool creation
- **`internal/migrations/005_seed_schema.sql`** — Full database schema

## Code Conventions

- **No comments** in code unless explicitly requested
- **Raw SQL** — no ORM, all queries in `internal/store/`
- **Error pattern**: check `domain.ErrNotFound` → 404, `domain.ErrIncorrectPassword` → 401, default → 500
- **Handler pattern**: decode JSON → validate → call store → return JSON
- **Auth**: wrap handler with `AuthMiddleware()` to require JWT cookie
- **JSON tags**: `snake_case`
- **Imports**: stdlib, external, internal (grouped with blank lines)

## Environment Variables

```
DATABASE_URL=postgresql://...   # Required
JWT_SECRET=...                  # Required (base64 32-byte key)
PORT=3000                       # Default 3000
API_PREFIX=api/v1               # Default api/v1
APP_ENV=                        # Set "production" for secure cookies
```

## Database

- PostgreSQL 14+ via `pgx/v5` with `pgxpool` connection pooling
- Pool: 2-10 connections, 1h max lifetime, 30min idle timeout
- Migrations auto-run on startup via `migrations.New(db).Up(ctx)`
- Schema tracked in `schema_migrations` table

## What NOT To Do

- Do not add comments or docstrings unless asked
- Do not introduce new dependencies without checking if the functionality exists in stdlib
- Do not modify `.env` or `.env.example` without explicit request
- Do not run `git commit` unless explicitly asked
- Do not add Docker/Makefile/CI configs unless asked
- Do not change the error response format (keep `{"error": "message"}`)
- Do not add ORM or query builder — keep raw SQL
- Do not skip the `go vet` step before suggesting code is done
