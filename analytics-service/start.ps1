# PowerShell script to start analytics service

# Create the lugx-network if it doesn't exist
Write-Host "Creating Docker network..."
docker network create lugx-network 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "Network lugx-network already exists or created successfully"
}

# Check if ClickHouse is running
Write-Host "Checking if ClickHouse is running..."
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8123/ping" -Method GET -TimeoutSec 5
    if ($response.StatusCode -eq 200) {
        Write-Host "ClickHouse is running!" -ForegroundColor Green
    }
} catch {
    Write-Host "WARNING: ClickHouse is not running!" -ForegroundColor Red
    Write-Host "Please start ClickHouse first by running the start script in the clickhouse folder" -ForegroundColor Yellow
    Write-Host "Path: ..\clickhouse\start.ps1" -ForegroundColor Yellow
    Write-Host ""
    $continue = Read-Host "Do you want to continue starting the analytics service anyway? (y/N)"
    if ($continue -ne "y" -and $continue -ne "Y") {
        Write-Host "Exiting..." -ForegroundColor Red
        exit 1
    }
}

# Start the analytics service
Write-Host "Starting Analytics Service..."
docker-compose up -d

# Wait for service to be ready
Write-Host "Waiting for Analytics Service to be ready..."
Start-Sleep -Seconds 5

# Check analytics service health
Write-Host "Checking Analytics Service health..."
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET -TimeoutSec 5
    if ($response.StatusCode -eq 200) {
        Write-Host "Analytics Service is ready!" -ForegroundColor Green
    }
} catch {
    Write-Host "Analytics Service not yet ready, please wait a moment" -ForegroundColor Yellow
}

Write-Host "`nAnalytics Service started!" -ForegroundColor Green
Write-Host "Analytics API: http://localhost:8080/api/analytics" -ForegroundColor Cyan
Write-Host "Note: ClickHouse Web UI: http://localhost:8123/play" -ForegroundColor Cyan
