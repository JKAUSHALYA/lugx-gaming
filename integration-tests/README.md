# LUGX Gaming - Integration Tests

This directory contains integration tests for all three microservices in the LUGX Gaming platform:

- **game-service**: Tests for game CRUD operations
- **order-service**: Tests for order management functionality
- **analytics-service**: Tests for analytics data collection and retrieval

## Prerequisites

Before running the integration tests, ensure that:

1. All services are running and accessible:

   - Game Service: `http://localhost:8080`
   - Order Service: `http://localhost:8081`
   - Analytics Service: `http://localhost:8080` (note: may conflict with game service)

2. Required databases are running:

   - PostgreSQL for game-service and order-service
   - ClickHouse for analytics-service

3. Go is installed (version 1.21 or later)

## Running the Tests

### Individual Service Tests

Navigate to each service directory and run the tests:

```powershell
# Game Service Tests
cd integration-tests\game-service
go test -v

# Order Service Tests
cd integration-tests\order-service
go test -v

# Analytics Service Tests
cd integration-tests\analytics-service
go test -v
```

### Running All Tests

From the integration-tests root directory:

```powershell
# Run all tests in sequence
Get-ChildItem -Directory | ForEach-Object {
    Write-Host "Testing $($_.Name)..." -ForegroundColor Green
    cd $_.FullName
    go test -v
    cd ..
}
```

## Test Coverage

### Game Service Tests

- ✅ Health check endpoint
- ✅ Create game with valid data
- ✅ Get all games
- ✅ Get games by category filter
- ✅ Get specific game by ID
- ✅ Update game details
- ✅ Delete game
- ✅ Invalid data validation

### Order Service Tests

- ✅ Health check endpoint
- ✅ Create order with items
- ✅ Get all orders (with pagination)
- ✅ Get order statistics
- ✅ Get specific order by ID
- ✅ Update order status
- ✅ Get orders by customer ID
- ✅ Delete order
- ✅ Invalid data validation

### Analytics Service Tests

- ✅ Health check endpoint
- ✅ Track page views
- ✅ Track click events
- ✅ Track scroll depth
- ✅ Track page time
- ✅ Track session time
- ✅ Retrieve analytics data
- ✅ Date range filtering
- ✅ Multiple event types
- ✅ Concurrent request handling
- ✅ Invalid data validation

## Service URLs and Ports

| Service           | Default URL           | Default Port |
| ----------------- | --------------------- | ------------ |
| Game Service      | http://localhost:8080 | 8080         |
| Order Service     | http://localhost:8081 | 8081         |
| Analytics Service | http://localhost:8080 | 8080         |

**Note**: Analytics and Game services both default to port 8080. Make sure to configure different ports or run them separately when testing.

## Environment Setup

### Starting Services

Use the provided VS Code tasks or start services manually:

```powershell
# Start Game Service
cd game-service
go mod tidy
go run main.go

# Start Order Service
cd order-service
go mod tidy
go run main.go

# Start Analytics Service
cd analytics-service
go mod tidy
go run main.go
```

### Database Setup

Ensure PostgreSQL and ClickHouse databases are running with proper schemas. Check each service's README for specific database setup instructions.

## Test Data

The integration tests create and clean up their own test data. However, some tests may leave residual data in the databases. For a clean test environment, consider resetting the databases between test runs.

## Troubleshooting

### Common Issues

1. **Connection Refused**: Ensure services are running on expected ports
2. **Database Errors**: Verify database connections and schemas
3. **Port Conflicts**: Analytics and Game services both use port 8080 by default
4. **Test Timeouts**: Increase timeout values if services are slow to respond

### Test Isolation

Tests are designed to be independent but may interact with shared databases. For complete isolation:

1. Use separate test databases
2. Run services in Docker containers
3. Reset databases between test suites

## Contributing

When adding new tests:

1. Follow the existing test patterns
2. Include proper error handling
3. Add descriptive test names and comments
4. Test both success and failure scenarios
5. Clean up test data when possible

## CI/CD Integration

These tests can be integrated into CI/CD pipelines. Consider:

1. Running services in Docker containers
2. Using test databases
3. Parallel test execution where appropriate
4. Proper test reporting and artifacts
