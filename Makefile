# Makefile

# Get ENV variables from .env file
include .env
export $(shell sed 's/=.*//' .env)

PATH_CURRENT := $(shell pwd)
GIT_COMMIT := $(shell git log --oneline -1 HEAD)

.PHONY: run swag db-up db-down migrate seed seed-test seed-prod

# Start PostgreSQL and Redis with Docker Compose
db-up:
	docker-compose up -d postgres redis

# Stop PostgreSQL and Redis
db-down:
	docker-compose down

# Migrate database
migrate:
	go run cmd/server/main.go migrate

# Seed database with all data
seed:
	go run cmd/server/main.go seed

# Seed database with test data
seed-test:
	go run cmd/server/main.go seed:test

# Seed database with production data
seed-prod:
	go run cmd/server/main.go seed:prod

# Run the main application
run:
	go run cmd/server/main.go

client: # run the client application: bun dev in next folder
	cd next && bun run dev

# Build the application
build:
	go build cmd/server/main.go

pre-build:
	go mod tidy
	go mod vendor

build-linux: pre-build
	env GOOS=linux GOARCH=amd64 go build -v -o ./build/server-linux -ldflags "-X 'main.GitCommitLog=${GIT_COMMIT}'" cmd/server/main.go

gcp: build-linux
	gcloud builds submit --config cloudbuild.yaml .

# Initialize swagger documentation
swag:
	swag init -g cmd/server/main.go --parseDependency --parseInternal --parseDepth 5

seeder:
	go run cmd/seeder/*.go
