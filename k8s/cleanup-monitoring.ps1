#!/usr/bin/env pwsh

# LUGX Gaming Monitoring Stack Cleanup Script
# This script removes Prometheus and Grafana monitoring from the cluster

Write-Host "🗑️  LUGX Gaming Monitoring Stack Cleanup" -ForegroundColor Red
Write-Host "=========================================" -ForegroundColor Red

# Check if kubectl is available
if (!(Get-Command kubectl -ErrorAction SilentlyContinue)) {
    Write-Host "❌ kubectl not found. Please install kubectl and configure it to connect to your cluster." -ForegroundColor Red
    exit 1
}

# Get the script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Confirm cleanup
$confirmation = Read-Host "⚠️  This will remove ALL monitoring components. Are you sure? (y/N)"
if ($confirmation -ne 'y' -and $confirmation -ne 'Y') {
    Write-Host "❌ Cleanup cancelled" -ForegroundColor Yellow
    exit 0
}

Write-Host "🗑️  Removing monitoring infrastructure..." -ForegroundColor Yellow

try {
    # Delete Grafana
    Write-Host "🔧 Removing Grafana..." -ForegroundColor Yellow
    kubectl delete -f "$ScriptDir\grafana.yaml" --ignore-not-found=true
    kubectl delete -f "$ScriptDir\grafana-config.yaml" --ignore-not-found=true
    
    # Delete Node Exporter
    Write-Host "🔧 Removing Node Exporter..." -ForegroundColor Yellow
    kubectl delete -f "$ScriptDir\node-exporter.yaml" --ignore-not-found=true
    
    # Delete Prometheus
    Write-Host "🔧 Removing Prometheus..." -ForegroundColor Yellow
    kubectl delete -f "$ScriptDir\prometheus.yaml" --ignore-not-found=true
    kubectl delete -f "$ScriptDir\prometheus-config.yaml" --ignore-not-found=true
    
    # Delete RBAC
    Write-Host "🔧 Removing Prometheus RBAC..." -ForegroundColor Yellow
    kubectl delete -f "$ScriptDir\prometheus-rbac.yaml" --ignore-not-found=true
    
    # Delete persistent volume claims
    Write-Host "🔧 Removing persistent volumes..." -ForegroundColor Yellow
    kubectl delete pvc -n monitoring prometheus-pvc grafana-pvc --ignore-not-found=true
    
    # Delete monitoring namespace (this will remove any remaining resources)
    Write-Host "🔧 Removing monitoring namespace..." -ForegroundColor Yellow
    kubectl delete namespace monitoring --ignore-not-found=true
    
    Write-Host ""
    Write-Host "✅ Monitoring stack cleanup completed!" -ForegroundColor Green
    Write-Host "All monitoring components have been removed from the cluster." -ForegroundColor White
    
}
catch {
    Write-Host "❌ Cleanup failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "You may need to manually remove some resources." -ForegroundColor Yellow
    exit 1
}
