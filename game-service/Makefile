# Makefile for Game Service

.PHONY: help build run test clean docker-build docker-up docker-down deps

# Default target
help:
	@echo "Available commands:"
	@echo "  make deps        - Install Go dependencies"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run API tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-up   - Start services with Docker Compose"
	@echo "  make docker-down - Stop Docker Compose services"

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build the application
build:
	go build -o bin/game-service main.go

# Run the application
run:
	go run main.go

# Run API tests
test:
	@echo "Starting API tests..."
	go run test_api.go

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Docker commands
docker-build:
	docker build -t game-service .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Development setup
dev-setup: deps
	@echo "Development environment setup complete"
	@echo "Make sure PostgreSQL is running or use 'make docker-up'"
