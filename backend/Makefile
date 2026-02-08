.PHONY: migrate-up migrate-down migrate-create migrate-version sqlc-generate

# Load environment variables
include .env

# Database connection string
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

# Run migrations up
migrate-up:
	@migrate -path database/migrations -database "$(DB_URL)" -verbose up $(filter-out $@,$(MAKECMDGOALS))

# Rollback migrations
migrate-down:
	@migrate -path database/migrations -database "$(DB_URL)" -verbose down $(filter-out $@,$(MAKECMDGOALS))

# Create a new migration file
migrate-create:
	@migrate create -ext sql -dir database/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

# Check migration version
migrate-version:
	@migrate -path database/migrations -database "$(DB_URL)" version

# Generate sqlc code
sqlc-generate:
	@sqlc generate

# Run the application
run: build
	@./bin/api

# Build the application
build:
	@go build -o bin/api cmd/api/main.go

# Run tests
test:
	@go test -v ./...

# Install dependencies
deps:
	@go mod download
	@go mod tidy
