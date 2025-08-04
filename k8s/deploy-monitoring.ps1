#!/usr/bin/env pwsh

# LUGX Gaming Monitoring Stack Deployment Script
# This script deploys Prometheus and Grafana monitoring for the LUGX Gaming platform

Write-Host "🔍 LUGX Gaming Monitoring Stack Deployment" -ForegroundColor Green
Write-Host "===========================================" -ForegroundColor Green

# Check if kubectl is available
if (!(Get-Command kubectl -ErrorAction SilentlyContinue)) {
    Write-Host "❌ kubectl not found. Please install kubectl and configure it to connect to your cluster." -ForegroundColor Red
    exit 1
}

# Get the script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

Write-Host "📊 Deploying monitoring infrastructure..." -ForegroundColor Cyan

try {
    # Deploy monitoring namespace first
    Write-Host "🔧 Creating monitoring namespace..." -ForegroundColor Yellow
    kubectl apply -f "$ScriptDir\monitoring-namespace.yaml"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to create monitoring namespace" -ForegroundColor Red
        exit 1
    }
    
    # Deploy RBAC for Prometheus
    Write-Host "🔧 Setting up Prometheus RBAC..." -ForegroundColor Yellow
    kubectl apply -f "$ScriptDir\prometheus-rbac.yaml"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to setup Prometheus RBAC" -ForegroundColor Red
        exit 1
    }
    
    # Deploy Prometheus configuration
    Write-Host "🔧 Deploying Prometheus configuration..." -ForegroundColor Yellow
    kubectl apply -f "$ScriptDir\prometheus-config.yaml"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to deploy Prometheus configuration" -ForegroundColor Red
        exit 1
    }
    
    # Deploy Prometheus
    Write-Host "🔧 Deploying Prometheus..." -ForegroundColor Yellow
    kubectl apply -f "$ScriptDir\prometheus.yaml"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to deploy Prometheus" -ForegroundColor Red
        exit 1
    }
    
    # Deploy Node Exporter
    Write-Host "🔧 Deploying Node Exporter..." -ForegroundColor Yellow
    kubectl apply -f "$ScriptDir\node-exporter.yaml"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to deploy Node Exporter" -ForegroundColor Red
        exit 1
    }
    
    # Deploy Grafana configuration
    Write-Host "🔧 Deploying Grafana configuration..." -ForegroundColor Yellow
    kubectl apply -f "$ScriptDir\grafana-config.yaml"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to deploy Grafana configuration" -ForegroundColor Red
        exit 1
    }
    
    # Deploy Grafana
    Write-Host "🔧 Deploying Grafana..." -ForegroundColor Yellow
    kubectl apply -f "$ScriptDir\grafana.yaml"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to deploy Grafana" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "✅ Monitoring stack deployed successfully!" -ForegroundColor Green
    
    # Wait for pods to be ready
    Write-Host "⏳ Waiting for monitoring pods to be ready..." -ForegroundColor Yellow
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/component=monitoring -n monitoring --timeout=300s
    
    # Get the node IP for accessing services
    $NodeIP = kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}'
    if ([string]::IsNullOrEmpty($NodeIP)) {
        $NodeIP = "localhost"
    }
    
    Write-Host ""
    Write-Host "🎉 Monitoring Stack Deployment Complete!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "📊 Access your monitoring services:" -ForegroundColor Cyan
    Write-Host "  Prometheus: http://$NodeIP`:30090" -ForegroundColor White
    Write-Host "  Grafana:    http://$NodeIP`:30300" -ForegroundColor White
    Write-Host ""
    Write-Host "🔐 Grafana Credentials:" -ForegroundColor Cyan
    Write-Host "  Username: admin" -ForegroundColor White
    Write-Host "  Password: admin" -ForegroundColor White
    Write-Host ""
    Write-Host "📈 Pre-configured Dashboard:" -ForegroundColor Cyan
    Write-Host "  - LUGX Gaming Platform Monitoring" -ForegroundColor White
    Write-Host ""
    Write-Host "🔍 Monitoring Features:" -ForegroundColor Cyan
    Write-Host "  ✅ Service Health Monitoring" -ForegroundColor Green
    Write-Host "  ✅ Performance Metrics (CPU, Memory)" -ForegroundColor Green
    Write-Host "  ✅ Availability Tracking" -ForegroundColor Green
    Write-Host "  ✅ Node-level Metrics" -ForegroundColor Green
    Write-Host "  ✅ Kubernetes Cluster Metrics" -ForegroundColor Green
    Write-Host "  ✅ Alerting Rules" -ForegroundColor Green
    Write-Host ""
    Write-Host "💡 Tip: After accessing Grafana, the pre-configured dashboard" -ForegroundColor Yellow
    Write-Host "   should automatically be available showing your LUGX Gaming" -ForegroundColor Yellow
    Write-Host "   platform metrics!" -ForegroundColor Yellow
    
}
catch {
    Write-Host "❌ Deployment failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
