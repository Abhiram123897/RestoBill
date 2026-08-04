.PHONY: run build test clean migrate-up migrate-down docker-up docker-down

# Run the server
run:
	go run cmd/server/main.go

# Build binary
build:
	go build -o restobill cmd/server/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f restobill

# Start Docker services (PostgreSQL)
docker-up:
	docker compose up -d

# Stop Docker services
docker-down:
	docker compose down

# Run database migrations
migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/restobill?sslmode=disable" up

# Rollback database migrations
migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/restobill?sslmode=disable" down
