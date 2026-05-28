# Go News API

Read this in [Bahasa Indonesia](README.id.md) 🇮🇩

Go News API is a high-performance RESTful backend designed for news portals, blogs. Built with **Golang**, this project adheres to **Clean Architecture** principles to ensure codebase scalability and maintainability.

## Key Features
* **Clean Architecture Layering**: Clear separation of concerns between *Domain*, *Repository*, *Service*, and *Handler*.
* **Dynamic Content Storage**: Utilizes `[]map[string]interface{}` mapped to PostgreSQL's `jsonb` column, making it highly adaptable to any modern frontend block-editor payloads (e.g., Editor.js, Quill).
* **Blazing Fast Redis Caching**: Implements a *Cache-Aside* pattern on public endpoints to handle high-traffic spikes seamlessly, complete with automated cache invalidation upon admin CRUD operations.
* **Authentication & Authorization**: Secured via JWT Tokens and Role-Based Access Control (RBAC).
* **Separation of Routes**: Dedicated public routes for readers (Client UI) and protected private routes for the CMS dashboard (Admin UI).
* **Frontend Ready**: Fully configured CORS to safely accept requests from frameworks like React, Vue, or Angular.

## Tech Stack
* **Language**: Go (1.25.0)
* **Web Framework**: [Gin-Gonic](https://gin-gonic.com/)
* **Database**: PostgreSQL (via [GORM](https://gorm.io/))
* **Caching**: Redis (via `go-redis/v9`)
* **Security**: JWT (`golang-jwt/v5`) & Bcrypt

## Project Structure
```text
.
├── config/              # Env configs, Database Connections (PostgreSQL & Redis)
├── internal/
│   ├── domain/          # Core entities, DTOs, Repositories, Services, Handlers
│   │   ├── article/     # Article Module (includes Cache helpers)
│   │   └── user/        # User Module (Auth & RBAC)
│   ├── middleware/      # JWT & Role validation middlewares
│   └── utils/           # Shared utilities (JWT generator, etc.)
├── routes/              # Gin routing & CORS setup
├── main.go              # Application entry point
└── .env                 # Environment variables
````

## Getting Started

### 1. Prerequisites

- [Go](https://golang.org/dl/) installed
    
- [PostgreSQL](https://www.postgresql.org/) running
    
- [Docker](https://www.docker.com/) (to easily run Redis)
    

### 2. Environment Setup

Copy `.env.example` to `.env` and fill in your credentials:

Cuplikan kode

```
DB_HOST=localhost
DB_USER=postgres
DB_PASS=your_password
DB_NAME=news_db
DB_PORT=5432
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=your_super_secret_key
```

### 3. Spin Up Redis via Docker

Bash

```
docker run --name news-api-redis -p 6379:6379 -d redis
```

### 4. Install Dependencies & Run

Bash

```
go mod tidy
go run main.go
```

The server will start at `http://localhost:8888`. The database tables will be migrated automatically.

## API Endpoints (Brief Overview)

**Authentication (Public):**
* `POST /register` - Register a new user account
* `POST /login` - Authenticate user and receive a JWT

**Public / Reader Website (No Token Required - Cache-Aside Pattern):**
* `GET /public/articles` - Retrieve a list of all articles for readers
* `GET /public/articles/:slug` - Retrieve article details by its slug

**Admin / CMS Dashboard (Requires Bearer Token, Role: 'admin'):**
* `GET /admin/articles` - Retrieve all articles (for admin data tables)
* `GET /admin/article/:id` - Retrieve specific article details by ID
* `POST /admin/article` - Create a new article (Automatically invalidates public cache)
* `PUT /admin/article/:id` - Update an existing article (Automatically invalidates public cache)
* `DELETE /admin/article/:id` - Delete an article (Automatically invalidates public cache)