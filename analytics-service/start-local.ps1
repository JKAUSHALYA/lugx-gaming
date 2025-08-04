# Set environment variables for local development
$env:CLICKHOUSE_HOST = "localhost"
$env:CLICKHOUSE_PORT = "9000"
$env:CLICKHOUSE_USER = "default"
$env:CLICKHOUSE_PASSWORD = "password"
$env:CLICKHOUSE_DB = "analytics"
$env:PORT = "8080"

Write-Host "Environment variables set:"
Write-Host "CLICKHOUSE_HOST: $env:CLICKHOUSE_HOST"
Write-Host "CLICKHOUSE_PORT: $env:CLICKHOUSE_PORT"
Write-Host "CLICKHOUSE_USER: $env:CLICKHOUSE_USER"
Write-Host "CLICKHOUSE_DB: $env:CLICKHOUSE_DB"
Write-Host ""

# Start the analytics service
Write-Host "Starting analytics service..."
go run main.go
