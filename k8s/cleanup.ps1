#!/usr/bin/env pwsh

# LUGX Gaming Kubernetes Cleanup Script
# This script removes all Kubernetes resources for the LUGX Gaming application

Write-Host "🧹 LUGX Gaming Kubernetes Cleanup" -ForegroundColor Red
Write-Host "==================================" -ForegroundColor Red

# Check if kubectl is available
if (!(Get-Command kubectl -ErrorAction SilentlyContinue)) {
    Write-Host "❌ kubectl not found. Please install kubectl and configure it to connect to your cluster." -ForegroundColor Red
    exit 1
}

# Confirm deletion
$confirmation = Read-Host "⚠️  This will delete ALL LUGX Gaming resources from Kubernetes. Are you sure? (y/N)"
if ($confirmation -ne 'y' -and $confirmation -ne 'Y') {
    Write-Host "❌ Cleanup cancelled." -ForegroundColor Yellow
    exit 0
}

Write-Host "🗑️ Deleting LUGX Gaming namespace and all resources..." -ForegroundColor Red

# Delete the entire namespace (this removes all resources)
kubectl delete namespace lugx-gaming

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Successfully deleted all LUGX Gaming resources." -ForegroundColor Green
} else {
    Write-Host "❌ Failed to delete resources. Please check manually." -ForegroundColor Red
    exit 1
}

Write-Host "🧹 Cleanup completed!" -ForegroundColor Green
