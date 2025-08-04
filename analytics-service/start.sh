#!/bin/bash

# Create the lugx-network if it doesn't exist
docker network create lugx-network 2>/dev/null || echo "Network lugx-network already exists"

# Start the analytics service and ClickHouse
echo "Starting Analytics Service and ClickHouse..."
docker-compose up -d

# Wait for ClickHouse to be ready
echo "Waiting for ClickHouse to be ready..."
sleep 10

# Check if ClickHouse is responding
echo "Checking ClickHouse health..."
curl -f http://localhost:8123/ping || echo "ClickHouse not yet ready, please wait a moment"

# Check analytics service health
echo "Checking Analytics Service health..."
curl -f http://localhost:8080/health || echo "Analytics Service not yet ready, please wait a moment"

echo "Analytics setup complete!"
echo "ClickHouse Web UI: http://localhost:8123/play"
echo "Analytics API: http://localhost:8080/api/analytics"
