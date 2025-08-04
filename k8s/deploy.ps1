#!/usr/bin/env pwsh

# LUGX Gaming Kubernetes Deployment Script
# This script builds Docker images and deploys the application to Kubernetes

Write-Host "🚀 LUGX Gaming Kubernetes Deployment" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green

# Check if kubectl is available
if (!(Get-Command kubectl -ErrorAction SilentlyContinue)) {
    Write-Host "❌ kubectl not found. Please install kubectl and configure it to connect to your cluster." -ForegroundColor Red
    exit 1
}

# Check if Docker is available
if (!(Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Docker not found. Please install Docker." -ForegroundColor Red
    exit 1
}

# Function to build Docker image
function Build-DockerImage {
    param(
        [string]$ServicePath,
        [string]$ImageName
    )
    
    Write-Host "🔨 Building $ImageName..." -ForegroundColor Yellow
    Push-Location $ServicePath
    docker build -t $ImageName .
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to build $ImageName" -ForegroundColor Red
        Pop-Location
        exit 1
    }
    Pop-Location
    Write-Host "✅ Successfully built $ImageName" -ForegroundColor Green
}

# Get the script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $ScriptDir

Write-Host "📦 Building Docker Images..." -ForegroundColor Cyan

# Build all Docker images
Build-DockerImage -ServicePath "$RootDir\front-end" -ImageName "lugx-gaming-frontend:latest"
Build-DockerImage -ServicePath "$RootDir\game-service" -ImageName "lugx-game-service:latest"  
Build-DockerImage -ServicePath "$RootDir\order-service" -ImageName "lugx-order-service:latest"
Build-DockerImage -ServicePath "$RootDir\analytics-service" -ImageName "lugx-analytics-service:latest"

Write-Host "🗄️ Starting external databases..." -ForegroundColor Cyan

# Start PostgreSQL
Write-Host "🐘 Starting PostgreSQL..." -ForegroundColor Yellow
Push-Location "$RootDir\postgres"
docker-compose up -d
if ($LASTEXITCODE -ne 0) {
    Write-Host "⚠️ PostgreSQL may already be running or failed to start" -ForegroundColor Yellow
}
Pop-Location

Write-Host "🎯 Deploying to Kubernetes..." -ForegroundColor Cyan

# Apply Kubernetes manifests
Write-Host "📝 Creating namespace..." -ForegroundColor Yellow
kubectl apply -f "$ScriptDir\namespace.yaml"

Write-Host "🔧 Configuring external services..." -ForegroundColor Yellow
kubectl apply -f "$ScriptDir\external-services.yaml"

Write-Host "🏠 Deploying ClickHouse..." -ForegroundColor Yellow
kubectl apply -f "$ScriptDir\clickhouse.yaml"

Write-Host "🎮 Deploying Game Service..." -ForegroundColor Yellow
kubectl apply -f "$ScriptDir\game-service.yaml"

Write-Host "🛒 Deploying Order Service..." -ForegroundColor Yellow
kubectl apply -f "$ScriptDir\order-service.yaml"

Write-Host "📈 Deploying Analytics Service..." -ForegroundColor Yellow
kubectl apply -f "$ScriptDir\analytics-service.yaml"

Write-Host "🌐 Deploying Frontend..." -ForegroundColor Yellow
kubectl apply -f "$ScriptDir\frontend.yaml"

# Wait for all deployments to be ready
Write-Host "⏳ Waiting for all services to be ready..." -ForegroundColor Yellow
kubectl wait --for=condition=available deployment --all -n lugx-gaming --timeout=300s

# Ask if user wants to deploy monitoring
Write-Host ""
$deployMonitoring = Read-Host "🔍 Would you like to deploy Prometheus and Grafana monitoring? (Y/n)"
if ($deployMonitoring -ne 'n' -and $deployMonitoring -ne 'N') {
    Write-Host "📊 Deploying monitoring stack..." -ForegroundColor Cyan
    & "$ScriptDir\deploy-monitoring.ps1"
}

Write-Host "🎉 Deployment completed successfully!" -ForegroundColor Green
Write-Host ""

# Display service information
Write-Host "📊 Service Status:" -ForegroundColor Cyan
kubectl get pods -n lugx-gaming
Write-Host ""

Write-Host "🌍 Service URLs:" -ForegroundColor Cyan
Write-Host "Frontend:          http://localhost:30000" -ForegroundColor White
Write-Host "Game Service:      http://localhost:30080" -ForegroundColor White  
Write-Host "Order Service:     http://localhost:30081" -ForegroundColor White
Write-Host "Analytics Service: http://localhost:30082" -ForegroundColor White
Write-Host "ClickHouse HTTP:   http://localhost:30123" -ForegroundColor White
Write-Host "ClickHouse Native: localhost:30900" -ForegroundColor White
Write-Host ""
Write-Host "📊 Monitoring URLs (if deployed):" -ForegroundColor Cyan
Write-Host "Prometheus:        http://localhost:30090" -ForegroundColor White
Write-Host "Grafana:           http://localhost:30300 (admin/admin)" -ForegroundColor White
Write-Host ""

Write-Host "📋 Useful Commands:" -ForegroundColor Cyan
Write-Host "View all resources: kubectl get all -n lugx-gaming" -ForegroundColor White
Write-Host "View logs:         kubectl logs -n lugx-gaming deployment/<service-name>" -ForegroundColor White
Write-Host "Scale service:     kubectl scale deployment <service-name> --replicas=<number> -n lugx-gaming" -ForegroundColor White
Write-Host "Delete all:        kubectl delete namespace lugx-gaming" -ForegroundColor White
