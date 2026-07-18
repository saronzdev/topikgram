# Architecture

## Overview

NotikGram follows a clean layered architecture with three main layers:

```
HTTP Request
    │
    ▼
┌─────────────────────┐
│   Transport Layer    │  ← Handlers, middleware, JSON serialization
│   (internal/transport)│
├─────────────────────┤
│     Store Layer      │  ← Database queries, business logic
│     (internal/store) │
├─────────────────────┤
│    PostgreSQL DB     │  ← Raw SQL via pgx
└─────────────────────┘
```

Additional packages provide cross-cutting concerns: `domain` (models/errors), `validator` (input validation), `ratelimit` (token bucket), and `migrations` (schema management).

## Packages

| Package | Responsibility |
|---------|---------------|
| `main` | Application bootstrap: env loading, DB connection, migration execution, route registration, middleware composition, server lifecycle |
| `cmd/keygen` | Standalone utility to generate a base64-encoded 32-byte JWT secret |
| `cmd/migrate` | Standalone CLI to run database migrations outside the main server |
| `domain` | Core data models (`User`, `Post`, `Comment`, etc.), `TopicID` enum constants, sentinel errors (`ErrNotFound`, `ErrIncorrectPassword`) |
| `store` | Data access layer. Each sub-store wraps a `*store.Pool` and contains all PostgreSQL queries for its domain: `AuthStore`, `UserStore`, `PostStore`, `CommentStore`, `InterestStore` |
| `transport` | HTTP layer: request handlers (per domain), JSON response helpers, error-to-HTTP mapping, middleware (CORS, logging, JWT auth, body limit, timeout, rate limit) |
| `validator` | Stateless input validation functions: `Email()`, `Username()`, `Password()`, `MaxLength()`, `MinLength()` |
| `ratelimit` | Token bucket rate limiter keyed by client IP address |
| `migrations` | Custom migration engine using Go's `embed.FS` for SQL files and a `schema_migrations` table for version tracking |

## Request Flow

A typical authenticated request follows this path:

```
1. Client sends: POST /api/v1/posts
2. CORS middleware          → adds CORS headers
3. Logging middleware       → logs method, path, status
4. RequestTimeout (30s)     → sets deadline on context
5. BodyLimit (1MB)          → rejects oversized bodies
6. RateLimit                → checks token bucket for IP
7. AuthMiddleware           → reads JWT cookie, validates token,
                              injects user_id into context
8. Handler (PostHandler.Create)
   → decodes JSON body
   → validates input
   → calls PostStore.Create()
   → returns JSON response
9. PostStore.Create()       → executes INSERT SQL, returns Post
```

## Middleware Stack

Applied in `main.go` (outside-in order):

```
CorsMiddleware
  └─ LogMiddleware
       └─ RequestTimeoutMiddleware(30s)
            └─ BodyLimitMiddleware(1MB)
                 └─ RateLimit (auth: 5/min, api: 60/min)
                      └─ Mux (routes)
```

Each handler that requires authentication wraps itself with `AuthMiddleware`, which:
1. Reads the `token` cookie from the request
2. Validates the JWT signature and expiry
3. Extracts `user_id` from the token claims
4. Injects it into the request context via `UserIDFromCtx()`

## Authentication

- **Token format**: HMAC-signed JWT (`golang-jwt/jwt/v5`)
- **Storage**: HTTP-only cookie named `token`, SameSite=Strict, Secure in production
- **Expiry**: 7 days
- **Claims payload**: `{ "user_id": <int>, "exp": <unix_timestamp> }`
- **Password hashing**: bcrypt (default cost)

Registration and login both return a `Set-Cookie` header with the JWT. Subsequent requests must include the cookie.

## Database

- **Driver**: pgx/v5 with connection pooling via pgxpool
- **Pool config**: 2-10 connections, 1h max lifetime, 30min idle timeout
- **Query style**: Raw SQL (no ORM), parameterized queries via `pgx`

### Connection Pool

```go
store.NewPool(ctx, connStr)  // returns *pgxpool.Pool
```

Each store function receives `context.Context` and uses the shared pool. The pool is closed on server shutdown.

## Migrations

Custom engine using Go's `embed.FS`:

1. SQL files are embedded at compile time from `internal/migrations/*.sql`
2. Files are sorted by filename prefix (e.g., `005_seed_schema.sql`)
3. Applied versions are tracked in a `schema_migrations` table
4. `migrations.New(db).Up(ctx)` runs all pending migrations on startup

### Adding a Migration

Create a new file in `internal/migrations/` with the next sequential number:

```
008_your_migration_name.sql
```

The file is automatically picked up by the embed system. Use `go run cmd/migrate/main.go -dir=up` to test standalone.

## Interest Weighting System

Tracks user preferences across 21 topic categories (`TopicID` constants in `domain/models.go`).

### Weight Updates

| Action | Delta | SQL Mechanism |
|--------|-------|---------------|
| Like | +0.1 | UPSERT with `LEAST(1.0, GREATEST(0.0, weight + delta))` |
| Unlike | -0.1 | Same, negative delta |
| Save | +0.2 | Same |
| Unsave | -0.2 | Same |
| Comment | +0.15 | Same |

The `interests` table stores `(user_id, topic_id, weight)` tuples. Weights are clamped to `[0.0, 1.0]` in SQL using `LEAST/GREATEST`.

## Rate Limiting

Token bucket algorithm per IP address (via `X-Forwarded-For` or `X-Real-IP` or remote addr):

- **Auth endpoints**: 5 tokens/second capacity, refill 5/60 per second
- **API endpoints**: 60 tokens/second capacity, refill 60/60 per second

## Error Handling

Domain-level sentinel errors in `domain/errors.go`:
- `ErrNotFound` → HTTP 404
- `ErrIncorrectPassword` → HTTP 401
- `ErrDuplicateKey` (username/email conflict) → HTTP 409

The `mapStoreError()` function in `transport/helpers.go` maps store errors to human-readable messages. Handlers use `jsonResp()`, `badRequest()`, `notFound()`, `unauthorized()`, `conflict()`, and `internalError()` helpers for consistent JSON error responses.

## Graceful Shutdown

1. Server listens for `SIGINT` or `SIGTERM`
2. On signal, initiates `srv.Shutdown()` with a 10-second timeout
3. Existing requests complete within the timeout window
4. Database pool is closed via `defer db.Close()`
5. Process exits cleanly
