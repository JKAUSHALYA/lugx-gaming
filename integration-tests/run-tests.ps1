#!/usr/bin/env pwsh

<#
.SYNOPSIS
    Runs integration tests for all LUGX Gaming microservices
.DESCRIPTION
    This script runs integration tests for game-service, order-service, and analytics-service.
    It can run tests for individual services or all services sequentially.
.PARAMETER Service
    Specify which service to test: 'game', 'order', 'analytics', or 'all' (default)
.PARAMETER VerboseOutput
    Run tests with verbose output
.PARAMETER Parallel
    Run service tests in parallel (experimental)
.EXAMPLE
    .\run-tests.ps1
    Runs all integration tests
.EXAMPLE
    .\run-tests.ps1 -Service game -VerboseOutput
    Runs only game service tests with verbose output
#>

param(
    [Parameter(Mandatory = $false)]
    [ValidateSet('game', 'order', 'analytics', 'all')]
    [string]$Service = 'all',
    
    [Parameter(Mandatory = $false)]
    [switch]$VerboseOutput,
    
    [Parameter(Mandatory = $false)]
    [switch]$Parallel
)

# Color functions for better output
function Write-Success($message) {
    Write-Host "✅ $message" -ForegroundColor Green
}

function Write-Error($message) {
    Write-Host "❌ $message" -ForegroundColor Red
}

function Write-Info($message) {
    Write-Host "ℹ️  $message" -ForegroundColor Cyan
}

function Write-Warning($message) {
    Write-Host "⚠️  $message" -ForegroundColor Yellow
}

# Function to check if a service is running
function Test-ServiceHealth {
    param(
        [string]$ServiceName,
        [string]$HealthUrl
    )
    
    try {
        $response = Invoke-WebRequest -Uri $HealthUrl -Method Get -TimeoutSec 5 -UseBasicParsing
        if ($response.StatusCode -eq 200) {
            Write-Success "$ServiceName is running"
            return $true
        }
        else {
            Write-Warning "$ServiceName returned status code $($response.StatusCode)"
            return $false
        }
    }
    catch {
        Write-Error "$ServiceName is not accessible at $HealthUrl"
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Function to run tests for a specific service
function Invoke-ServiceTests {
    param(
        [string]$ServiceName,
        [string]$ServicePath,
        [bool]$VerboseOutput = $false
    )
    
    Write-Info "Running $ServiceName integration tests..."
    
    if (-not (Test-Path $ServicePath)) {
        Write-Error "Service test directory not found: $ServicePath"
        return $false
    }
    
    Push-Location $ServicePath
    
    try {
        # Initialize go module if needed
        if (-not (Test-Path "go.sum")) {
            Write-Info "Initializing Go module for $ServiceName tests..."
            & go mod tidy
        }
        
        # Run tests
        $testArgs = @('test')
        if ($VerboseOutput) {
            $testArgs += '-v'
        }
        $testArgs += '-timeout=30s'
        
        Write-Info "Executing: go $($testArgs -join ' ')"
        $result = & go @testArgs
        $exitCode = $LASTEXITCODE
        
        if ($exitCode -eq 0) {
            Write-Success "$ServiceName tests completed successfully"
            return $true
        }
        else {
            Write-Error "$ServiceName tests failed with exit code $exitCode"
            Write-Host $result -ForegroundColor Red
            return $false
        }
    }
    catch {
        Write-Error "Error running $ServiceName tests: $($_.Exception.Message)"
        return $false
    }
    finally {
        Pop-Location
    }
}

# Main execution
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$integrationTestsPath = $scriptPath

Write-Info "LUGX Gaming - Integration Test Runner"
Write-Info "======================================"

# Service configurations
$services = @{
    'game'      = @{
        Name      = 'Game Service'
        Path      = Join-Path $integrationTestsPath 'game-service'
        HealthUrl = 'http://localhost:30080/api/v1/health'
    }
    'order'     = @{
        Name      = 'Order Service'
        Path      = Join-Path $integrationTestsPath 'order-service'
        HealthUrl = 'http://localhost:30081/health'
    }
    'analytics' = @{
        Name      = 'Analytics Service'
        Path      = Join-Path $integrationTestsPath 'analytics-service'
        HealthUrl = 'http://localhost:30082/health'
    }
}

# Determine which services to test
$servicesToTest = @()
if ($Service -eq 'all') {
    $servicesToTest = @('game', 'order', 'analytics')
}
else {
    $servicesToTest = @($Service)
}

Write-Info "Testing services: $($servicesToTest -join ', ')"

# Check service health before running tests
Write-Info "Checking service health..."
$healthCheckPassed = $true

foreach ($serviceKey in $servicesToTest) {
    $serviceConfig = $services[$serviceKey]
    $isHealthy = Test-ServiceHealth -ServiceName $serviceConfig.Name -HealthUrl $serviceConfig.HealthUrl
    if (-not $isHealthy) {
        $healthCheckPassed = $false
    }
}

if (-not $healthCheckPassed) {
    Write-Warning "Some services are not healthy. Tests may fail."
    Write-Info "Make sure all required services are running:"
    foreach ($serviceKey in $servicesToTest) {
        $serviceConfig = $services[$serviceKey]
        Write-Host "  - $($serviceConfig.Name): $($serviceConfig.HealthUrl)" -ForegroundColor Yellow
    }
    Write-Host ""
    
    $continue = Read-Host "Continue with tests anyway? (y/N)"
    if ($continue -ne 'y' -and $continue -ne 'Y') {
        Write-Info "Tests cancelled by user"
        exit 1
    }
}

# Run tests
$testResults = @{}
$startTime = Get-Date

if ($Parallel -and $servicesToTest.Count -gt 1) {
    Write-Info "Running tests in parallel..."
    
    $jobs = @()
    foreach ($serviceKey in $servicesToTest) {
        $serviceConfig = $services[$serviceKey]
        $job = Start-Job -ScriptBlock {
            param($serviceName, $servicePath, $verbose)
            
            Push-Location $servicePath
            try {
                $testArgs = @('test')
                if ($verbose) { $testArgs += '-v' }
                $testArgs += '-timeout=30s'
                
                $result = & go @testArgs 2>&1
                $exitCode = $LASTEXITCODE
                
                return @{
                    Service  = $serviceName
                    Success  = ($exitCode -eq 0)
                    Output   = $result
                    ExitCode = $exitCode
                }
            }
            finally {
                Pop-Location
            }
        } -ArgumentList $serviceConfig.Name, $serviceConfig.Path, $VerboseOutput.IsPresent
        
        $jobs += $job
    }
    
    # Wait for all jobs to complete
    $jobs | Wait-Job | ForEach-Object {
        $result = Receive-Job $_
        $testResults[$result.Service] = $result.Success
        
        if ($result.Success) {
            Write-Success "$($result.Service) tests passed"
        }
        else {
            Write-Error "$($result.Service) tests failed"
            if ($VerboseOutput) {
                Write-Host $result.Output -ForegroundColor Red
            }
        }
        
        Remove-Job $_
    }
}
else {
    # Run tests sequentially
    Write-Info "Running tests sequentially..."
    
    foreach ($serviceKey in $servicesToTest) {
        $serviceConfig = $services[$serviceKey]
        $success = Invoke-ServiceTests -ServiceName $serviceConfig.Name -ServicePath $serviceConfig.Path -VerboseOutput $VerboseOutput.IsPresent
        $testResults[$serviceConfig.Name] = $success
        
        if (-not $success) {
            Write-Warning "Continuing with remaining tests..."
        }
        
        Write-Host "" # Add spacing between service tests
    }
}

# Summary
$endTime = Get-Date
$duration = $endTime - $startTime

Write-Info "Test Summary"
Write-Info "============"
Write-Info "Total duration: $($duration.TotalSeconds.ToString('F2')) seconds"

$passed = 0
$failed = 0

foreach ($serviceName in $testResults.Keys) {
    if ($testResults[$serviceName]) {
        Write-Success "${serviceName}: PASSED"
        $passed++
    }
    else {
        Write-Error "${serviceName}: FAILED"
        $failed++
    }
}

Write-Host ""
if ($failed -eq 0) {
    Write-Success "All tests passed! ($passed/$($testResults.Count))"
    exit 0
}
else {
    Write-Error "Some tests failed. Passed: $passed, Failed: $failed"
    exit 1
}
