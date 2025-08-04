.PHONY: build run stop clean logs network

# Default target
all: network build run

# Create Docker network
network:
	@echo "Creating Docker network..."
	@docker network create lugx-network 2>/dev/null || echo "Network already exists"

# Build the analytics service
build:
	@echo "Building analytics service..."
	@docker-compose build

# Run the services
run:
	@echo "Starting analytics service and ClickHouse..."
	@docker-compose up -d

# Stop the services
stop:
	@echo "Stopping analytics service and ClickHouse..."
	@docker-compose down

# Clean up everything
clean: stop
	@echo "Cleaning up volumes and images..."
	@docker-compose down -v
	@docker image prune -f

# View logs
logs:
	@docker-compose logs -f

# Check service health
health:
	@echo "Checking ClickHouse health..."
	@curl -f http://localhost:8123/ping || echo "ClickHouse not ready"
	@echo "Checking Analytics Service health..."
	@curl -f http://localhost:8080/health || echo "Analytics Service not ready"

# Development mode (rebuild and restart)
dev: stop build run logs

# Initialize Go modules
init:
	@echo "Initializing Go modules..."
	@go mod tidy

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Show analytics URLs
info:
	@echo "Analytics Service URLs:"
	@echo "  API: http://localhost:8080/api/analytics"
	@echo "  Health: http://localhost:8080/health"
	@echo "ClickHouse URLs:"
	@echo "  HTTP: http://localhost:8123"
	@echo "  Web UI: http://localhost:8123/play"
	@echo "  TCP: localhost:9000"
