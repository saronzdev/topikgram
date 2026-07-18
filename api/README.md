# NotikGram

Social media REST API built with Go, PostgreSQL, and JWT authentication. Backend-only server supporting user profiles, posts with topic tags, follow mechanics, likes, saves, comments, and an interest-weighting system for feed personalization.

## Prerequisites

- Go 1.26+
- PostgreSQL 14+

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure:
   ```
   DATABASE_URL=postgres://user:password@localhost:5432/notikgram?sslmode=disable
   JWT_SECRET=your-secret-key-here
   PORT=3000
   APP_ENV=development
   ```
3. Generate a JWT secret key:
   ```
   go run cmd/keygen/main.go
   ```
4. Run the server (migrations run automatically on startup):
   ```
   go run main.go
   ```

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes | - | PostgreSQL connection string |
| `JWT_SECRET` | Yes | - | HMAC signing key for JWT tokens (base64-encoded 32-byte key) |
| `PORT` | No | `3000` | HTTP server listen port |
| `API_PREFIX` | No | `api/v1` | URL prefix for all API routes |
| `APP_ENV` | No | (empty) | Set to `production` to enable Secure flag on cookies |

## API Endpoints

### Auth

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/{prefix}/auth/register` | Register new user | No |
| POST | `/{prefix}/auth/login` | Login (email or username) | No |
| GET | `/{prefix}/auth/me` | Get current user | Yes |

### Users

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/{prefix}/users` | List all users (paginated) | No |
| GET | `/{prefix}/users/{id}` | Get user by ID | No |
| POST | `/{prefix}/users/{id}/follow` | Follow user | Yes |
| DELETE | `/{prefix}/users/{id}/follow` | Unfollow user | Yes |
| GET | `/{prefix}/users/{id}/followers` | Get user's followers | No |
| GET | `/{prefix}/users/{id}/following` | Get user's following | No |

### Posts

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/{prefix}/posts` | List posts (cursor pagination) | No |
| GET | `/{prefix}/posts/{id}` | Get post by ID | No |
| POST | `/{prefix}/posts` | Create post | Yes |
| PUT | `/{prefix}/posts/{id}` | Update post (owner only) | Yes |
| DELETE | `/{prefix}/posts/{id}` | Delete post (owner only) | Yes |
| POST | `/{prefix}/posts/{id}/like` | Like post | Yes |
| DELETE | `/{prefix}/posts/{id}/like` | Unlike post | Yes |
| POST | `/{prefix}/posts/{id}/save` | Save/bookmark post | Yes |
| DELETE | `/{prefix}/posts/{id}/save` | Unsave post | Yes |
| GET | `/{prefix}/posts/{id}/likes` | Get users who liked (paginated) | Yes |
| GET | `/{prefix}/posts/{id}/saves` | Get users who saved (paginated) | Yes |

### Comments

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/{prefix}/comments/{postid}` | Get comments for post | No |
| POST | `/{prefix}/comments` | Create comment | Yes |

### Health

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/health` | Health check | No |

## Request/Response Examples

### Register

```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "password": "secret123",
    "birthday": "1990-05-15"
  }'
```

Response `201 Created`:
```json
{
  "user": {
    "id": 1,
    "name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "birthday": "1990-05-15T00:00:00Z",
    "created_at": "2026-07-13T10:00:00Z"
  }
}
```

### Login

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "johndoe",
    "password": "secret123"
  }'
```

### Create Post

```bash
curl -X POST http://localhost:3000/api/v1/posts \
  -H "Content-Type: application/json" \
  -b "token=<jwt-cookie>" \
  -d '{
    "body": "Hello world! My first post.",
    "topics": [0, 4]
  }'
```

Response `201 Created`:
```json
{
  "id": 1,
  "user_id": 1,
  "body": "Hello world! My first post.",
  "topics_id": [0, 4],
  "created_at": "2026-07-13T10:05:00Z",
  "user": {"id": 1, "name": "John Doe", "username": "johndoe"},
  "likes": 0,
  "saved": false,
  "liked": false
}
```

### List Posts (Cursor Pagination)

```bash
curl http://localhost:3000/api/v1/posts?limit=10 \
  -b "token=<jwt-cookie>"
```

```json
{
  "posts": [...],
  "next_cursor": "2026-07-13T10:05:00.000Z",
  "has_more": true
}
```

Use `next_cursor` as the `cursor` query parameter for the next page. Omit it on the first request.

## Topics

Each post must include 1-3 topic tags. Topics are defined as integer IDs:

| ID | Topic | ID | Topic |
|----|-------|----|-------|
| 0 | General | 11 | Games |
| 1 | Programming | 12 | Literature |
| 2 | Cybersecurity | 13 | Travel |
| 3 | Entertainment | 14 | Cuisine |
| 4 | Funny | 15 | Tech |
| 5 | Art | 16 | Economy |
| 6 | Sports | 17 | Health |
| 7 | Politics | 18 | Philosophy |
| 8 | Science | 19 | Opinion |
| 9 | News | 20 | Ad |
| 10 | Cinema | | |

## Interest Weighting System

NotikGram tracks user interests across all 21 topic categories. When a user interacts with a post, the weights for that post's topics are adjusted:

| Action | Weight Delta |
|--------|-------------|
| Like post | +/- 0.1 |
| Save post | +/- 0.2 |
| Comment on post | +0.15 |

Weights are clamped to the range `[0.0, 1.0]`. This data can be used for feed personalization.

## Pagination

- **Posts feed**: Cursor-based. Use `?cursor=<timestamp>&limit=<n>` (default limit: 50, max: 100)
- **Likes/Saves lists**: Page-based. Use `?page=<n>&limit=<n>` (default limit: 20, max: 100)

## Database Schema

```
users          posts           comments
------         ------          --------
id (PK)        id (PK)         id (PK)
name           body            content
username (UQ)  topics_id[]     created_at
email (UQ)     created_at      user_id (FK -> users)
password       user_id (FK)    post_id (FK -> posts)
birthday       modified_at
created_at

follows        likes           saves
-------        -----           -----
follower_id    user_id         user_id
followee_id    post_id         post_id

interests
---------
user_id
topic_id
weight (0.0 - 1.0)
```

All foreign keys cascade on delete. See `internal/migrations/005_seed_schema.sql` for the full DDL.

## Security

- JWT tokens stored in HTTP-only, SameSite=Strict cookies (Secure flag in production)
- Token expiry: 7 days
- bcrypt password hashing
- Rate limiting: 5 req/min for auth endpoints, 60 req/min for API
- 1MB request body limit
- 30s request timeout
- Ownership enforcement on post update/delete

## Project Structure

```
notikgram/
├── cmd/
│   ├── keygen/              # JWT secret key generator
│   └── migrate/             # Standalone migration runner
├── internal/
│   ├── domain/              # Models, enums, sentinel errors
│   ├── migrations/          # SQL migration files + embed engine
│   ├── ratelimit/           # Token bucket rate limiter (per IP)
│   ├── store/               # Database access layer (PostgreSQL)
│   ├── transport/           # HTTP handlers + middleware
│   └── validator/           # Input validation functions
├── .env.example
├── go.mod
└── main.go
```

## Running Tests

```bash
go test ./internal/transport/ -v
```

## Building

```bash
go build -o notikgram .
```

## License

See [LICENSE](LICENSE) for details.
