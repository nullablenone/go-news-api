# Go News API

REST API service for managing news articles, user authentication, and role-based access control.

## Tech Stack

- Go 1.25
- Gin HTTP framework
- GORM with PostgreSQL
- Redis for optional article caching
- JWT for access and refresh token authentication
- Bcrypt for password hashing

## Project Structure

```text
.
├── cmd/
│   └── api/
│       └── main.go
├── docs/
│   ├── DEVELOPMENT.md
│   ├── DEVELOPMENT.id.md
│   └── openapi.yaml
├── internal/
│   ├── domain/
│   │   ├── article/
│   │   └── user/
│   ├── routes/
│   │   └── routes.go
│   └── shared/
│       ├── config/
│       ├── middleware/
│       └── utils/
├── .env.example
└── go.mod
```

## Quick Start

1. Copy `.env.example` to `.env` and fill in your database and JWT values.
2. Run the application:

```bash
go run ./cmd/api
```

## Documentation

- Developer guide in English: [docs/DEVELOPMENT.md][dev-en]
- Developer guide in Bahasa Indonesia: [docs/DEVELOPMENT.id.md][dev-id]
- OpenAPI specification: [docs/openapi.yaml][openapi]
- Swagger UI endpoint: `/docs`

[dev-en]: docs/DEVELOPMENT.md
[dev-id]: docs/DEVELOPMENT.id.md
[openapi]: docs/openapi.yaml
