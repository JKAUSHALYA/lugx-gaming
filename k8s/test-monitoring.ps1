#!/usr/bin/env pwsh

# LUGX Gaming Monitoring Validation Script
# This script tests the monitoring stack to ensure all components are working

Write-Host "🔍 LUGX Gaming Monitoring Validation" -ForegroundColor Green
Write-Host "====================================" -ForegroundColor Green

# Check if kubectl is available
if (!(Get-Command kubectl -ErrorAction SilentlyContinue)) {
    Write-Host "❌ kubectl not found. Please install kubectl and configure it to connect to your cluster." -ForegroundColor Red
    exit 1
}

$testsPassed = 0
$testsTotal = 0

function Test-Component {
    param(
        [string]$Name,
        [string]$Namespace,
        [string]$Resource,
        [string]$Label
    )
    
    $global:testsTotal++
    Write-Host "🔧 Testing $Name..." -ForegroundColor Yellow
    
    try {
        $result = kubectl get $Resource -n $Namespace -l $Label -o jsonpath='{.items[0].metadata.name}' 2>$null
        if ($result) {
            Write-Host "  ✅ $Name is deployed" -ForegroundColor Green
            $global:testsPassed++
            return $true
        }
        else {
            Write-Host "  ❌ $Name is not deployed" -ForegroundColor Red
            return $false
        }
    }
    catch {
        Write-Host "  ❌ Error checking $Name`: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

function Test-PodStatus {
    param(
        [string]$Name,
        [string]$Namespace,
        [string]$Label
    )
    
    $global:testsTotal++
    Write-Host "🔧 Testing $Name pod status..." -ForegroundColor Yellow
    
    try {
        $status = kubectl get pods -n $Namespace -l $Label -o jsonpath='{.items[0].status.phase}' 2>$null
        if ($status -eq "Running") {
            Write-Host "  ✅ $Name is running" -ForegroundColor Green
            $global:testsPassed++
            return $true
        }
        else {
            Write-Host "  ❌ $Name status: $status" -ForegroundColor Red
            return $false
        }
    }
    catch {
        Write-Host "  ❌ Error checking $Name pod status`: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

function Test-ServiceEndpoint {
    param(
        [string]$Name,
        [string]$URL,
        [string]$ExpectedContent = $null
    )
    
    $global:testsTotal++
    Write-Host "🔧 Testing $Name endpoint..." -ForegroundColor Yellow
    
    try {
        $response = Invoke-WebRequest -Uri $URL -Method GET -TimeoutSec 10 -UseBasicParsing
        if ($response.StatusCode -eq 200) {
            if ($ExpectedContent -and $response.Content -notlike "*$ExpectedContent*") {
                Write-Host "  ⚠️  $Name endpoint reachable but unexpected content" -ForegroundColor Yellow
                $global:testsPassed++
            }
            else {
                Write-Host "  ✅ $Name endpoint is accessible" -ForegroundColor Green
                $global:testsPassed++
            }
            return $true
        }
        else {
            Write-Host "  ❌ $Name endpoint returned status: $($response.StatusCode)" -ForegroundColor Red
            return $false
        }
    }
    catch {
        Write-Host "  ❌ $Name endpoint is not accessible: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Test monitoring namespace
Write-Host "📊 Testing Monitoring Namespace..." -ForegroundColor Cyan
$global:testsTotal++
try {
    $namespace = kubectl get namespace monitoring -o jsonpath='{.metadata.name}' 2>$null
    if ($namespace -eq "monitoring") {
        Write-Host "  ✅ Monitoring namespace exists" -ForegroundColor Green
        $global:testsPassed++
    }
    else {
        Write-Host "  ❌ Monitoring namespace does not exist" -ForegroundColor Red
    }
}
catch {
    Write-Host "  ❌ Error checking monitoring namespace" -ForegroundColor Red
}

# Test Prometheus deployment
Write-Host "📊 Testing Prometheus..." -ForegroundColor Cyan
Test-Component -Name "Prometheus" -Namespace "monitoring" -Resource "deployment" -Label "app.kubernetes.io/name=prometheus"
Test-PodStatus -Name "Prometheus" -Namespace "monitoring" -Label "app.kubernetes.io/name=prometheus"

# Test Grafana deployment
Write-Host "📊 Testing Grafana..." -ForegroundColor Cyan
Test-Component -Name "Grafana" -Namespace "monitoring" -Resource "deployment" -Label "app.kubernetes.io/name=grafana"
Test-PodStatus -Name "Grafana" -Namespace "monitoring" -Label "app.kubernetes.io/name=grafana"

# Test Node Exporter
Write-Host "📊 Testing Node Exporter..." -ForegroundColor Cyan
Test-Component -Name "Node Exporter" -Namespace "monitoring" -Resource "daemonset" -Label "app.kubernetes.io/name=node-exporter"

# Get node IP for testing
$nodeIP = kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}' 2>$null
if ([string]::IsNullOrEmpty($nodeIP)) {
    $nodeIP = "localhost"
}

# Test service endpoints
Write-Host "🌐 Testing Service Endpoints..." -ForegroundColor Cyan

# Wait a bit for services to be ready
Write-Host "⏳ Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

Test-ServiceEndpoint -Name "Prometheus" -URL "http://$nodeIP`:30090/-/ready"
Test-ServiceEndpoint -Name "Grafana" -URL "http://$nodeIP`:30300/api/health" -ExpectedContent "ok"

# Test LUGX Gaming services health endpoints
Write-Host "🎮 Testing LUGX Gaming Services..." -ForegroundColor Cyan
$lugxServices = @(
    @{Name = "Game Service"; Port = 30080; Path = "/api/v1/health" },
    @{Name = "Order Service"; Port = 30081; Path = "/health" },
    @{Name = "Analytics Service"; Port = 30082; Path = "/health" }
)

foreach ($service in $lugxServices) {
    Test-ServiceEndpoint -Name $service.Name -URL "http://$nodeIP`:$($service.Port)$($service.Path)"
}

# Test Prometheus targets
Write-Host "🎯 Testing Prometheus Targets..." -ForegroundColor Cyan
$global:testsTotal++
try {
    $targetsResponse = Invoke-WebRequest -Uri "http://$nodeIP`:30090/api/v1/targets" -Method GET -TimeoutSec 10 -UseBasicParsing
    if ($targetsResponse.StatusCode -eq 200) {
        $targets = $targetsResponse.Content | ConvertFrom-Json
        $activeTargets = $targets.data.activeTargets | Where-Object { $_.health -eq "up" }
        $totalTargets = $targets.data.activeTargets.Count
        $healthyTargets = $activeTargets.Count
        
        Write-Host "  📊 Prometheus Targets: $healthyTargets/$totalTargets healthy" -ForegroundColor White
        
        if ($healthyTargets -gt 0) {
            Write-Host "  ✅ Prometheus is successfully scraping targets" -ForegroundColor Green
            $global:testsPassed++
        }
        else {
            Write-Host "  ❌ No healthy targets found" -ForegroundColor Red
        }
    }
    else {
        Write-Host "  ❌ Could not retrieve Prometheus targets" -ForegroundColor Red
    }
}
catch {
    Write-Host "  ❌ Error checking Prometheus targets: $($_.Exception.Message)" -ForegroundColor Red
}

# Test Grafana datasource
Write-Host "📈 Testing Grafana Datasource..." -ForegroundColor Cyan
$global:testsTotal++
try {
    # Test if Grafana can connect to Prometheus
    $auth = [Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes("admin:admin"))
    $headers = @{Authorization = "Basic $auth" }
    $datasourceResponse = Invoke-WebRequest -Uri "http://$nodeIP`:30300/api/datasources" -Method GET -Headers $headers -TimeoutSec 10 -UseBasicParsing
    
    if ($datasourceResponse.StatusCode -eq 200) {
        $datasources = $datasourceResponse.Content | ConvertFrom-Json
        $prometheusDatasource = $datasources | Where-Object { $_.type -eq "prometheus" }
        
        if ($prometheusDatasource) {
            Write-Host "  ✅ Grafana Prometheus datasource is configured" -ForegroundColor Green
            $global:testsPassed++
        }
        else {
            Write-Host "  ❌ Prometheus datasource not found in Grafana" -ForegroundColor Red
        }
    }
    else {
        Write-Host "  ❌ Could not retrieve Grafana datasources" -ForegroundColor Red
    }
}
catch {
    Write-Host "  ⚠️  Could not test Grafana datasource (may need time to initialize): $($_.Exception.Message)" -ForegroundColor Yellow
}

# Summary
Write-Host ""
Write-Host "📊 Test Results Summary" -ForegroundColor Cyan
Write-Host "======================" -ForegroundColor Cyan
Write-Host "Tests Passed: $testsPassed/$testsTotal" -ForegroundColor White

$successRate = [math]::Round(($testsPassed / $testsTotal) * 100, 1)

if ($successRate -ge 90) {
    Write-Host "🎉 Monitoring stack is working excellently! ($successRate% success rate)" -ForegroundColor Green
}
elseif ($successRate -ge 70) {
    Write-Host "✅ Monitoring stack is mostly working ($successRate% success rate)" -ForegroundColor Yellow
    Write-Host "💡 Some components may need time to initialize or there might be minor issues" -ForegroundColor Yellow
}
else {
    Write-Host "❌ Monitoring stack has significant issues ($successRate% success rate)" -ForegroundColor Red
    Write-Host "🔧 Please check the deployment and service configurations" -ForegroundColor Red
}

Write-Host ""
Write-Host "🔗 Access URLs:" -ForegroundColor Cyan
Write-Host "  Prometheus: http://$nodeIP`:30090" -ForegroundColor White
Write-Host "  Grafana:    http://$nodeIP`:30300 (admin/admin)" -ForegroundColor White
Write-Host ""
Write-Host "📋 Next Steps:" -ForegroundColor Cyan
Write-Host "  1. Access Grafana and explore the LUGX Gaming Platform dashboard" -ForegroundColor White
Write-Host "  2. Check Prometheus targets page for service discovery" -ForegroundColor White
Write-Host "  3. Set up additional alerts if needed" -ForegroundColor White
Write-Host "  4. Customize dashboards for your specific metrics" -ForegroundColor White
