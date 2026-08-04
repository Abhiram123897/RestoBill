# RestoBill

A full-featured restaurant billing system built with Go, PostgreSQL, and React.

## Tech Stack

- **Backend:** Go + Gin
- **Database:** PostgreSQL
- **Frontend:** React (coming later)
- **Auth:** JWT

## Quick Start

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- golang-migrate CLI (for migrations)

### Run

```bash
# Start PostgreSQL
make docker-up

# Run migrations
make migrate-up

# Start the server
make run
```

Server runs at `http://localhost:8080`

### API

- Health check: `GET /api/v1/health`

## Project Structure

```
cmd/server/          - Application entry point
internal/
  config/            - Configuration loading
  domain/            - Domain models
  dto/               - Request/response objects
  handler/           - HTTP handlers
  middleware/        - Auth, logging middleware
  repository/        - Database access layer
  service/           - Business logic
  router/            - Route definitions
migrations/          - SQL migration files
```
