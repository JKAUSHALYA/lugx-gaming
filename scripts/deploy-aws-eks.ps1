#!/usr/bin/env pwsh

<#
.SYNOPSIS
    AWS EKS deployment script for LUGX Gaming Platform
.DESCRIPTION
    This script deploys the LUGX Gaming platform to AWS EKS with rolling releases
    and integration testing. It supports multiple environments and uses AWS managed PostgreSQL.
.PARAMETER Environment
    Target environment: development, staging, or production
.PARAMETER ImageTag
    Docker image tag to deploy
.PARAMETER SkipTests
    Skip integration tests after deployment
.PARAMETER RollbackOnFailure
    Automatically rollback on deployment failure
.EXAMPLE
    .\deploy-aws-eks.ps1 -Environment production -ImageTag v1.2.3
.EXAMPLE
    .\deploy-aws-eks.ps1 -Environment staging -ImageTag main-abc123 -SkipTests
#>

param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('development', 'staging', 'production')]
    [string]$Environment,
    
    [Parameter(Mandatory = $true)]
    [string]$ImageTag,
    
    [Parameter(Mandatory = $false)]
    [switch]$SkipTests,
    
    [Parameter(Mandatory = $false)]
    [switch]$RollbackOnFailure,
    
    [Parameter(Mandatory = $false)]
    [string]$AWSRegion = "us-east-1",
    
    [Parameter(Mandatory = $false)]
    [string]$ClusterName = "lugx-gaming-cluster"
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

function Write-Header($message) {
    Write-Host ""
    Write-Host "🚀 $message" -ForegroundColor Magenta
    Write-Host "=" * 50 -ForegroundColor Magenta
}

# Determine namespace based on environment
$Namespace = switch ($Environment) {
    'development' { 'lugx-gaming-dev' }
    'staging' { 'lugx-gaming-staging' }
    'production' { 'lugx-gaming-prod' }
}

# Set rollback default based on environment
if (-not $PSBoundParameters.ContainsKey('RollbackOnFailure')) {
    $RollbackOnFailure = $Environment -eq 'production'
}

$startTime = Get-Date

Write-Header "LUGX Gaming AWS EKS Deployment"
Write-Info "Environment: $Environment"
Write-Info "Namespace: $Namespace"
Write-Info "Image Tag: $ImageTag"
Write-Info "AWS Region: $AWSRegion"
Write-Info "EKS Cluster: $ClusterName"

# Check prerequisites
Write-Header "Checking Prerequisites"

# Check if AWS CLI is available
if (!(Get-Command aws -ErrorAction SilentlyContinue)) {
    Write-Error "AWS CLI not found. Please install AWS CLI."
    exit 1
}
Write-Success "AWS CLI found"

# Check if kubectl is available
if (!(Get-Command kubectl -ErrorAction SilentlyContinue)) {
    Write-Error "kubectl not found. Please install kubectl."
    exit 1
}
Write-Success "kubectl found"

# Check AWS credentials
try {
    $awsIdentity = aws sts get-caller-identity --output json | ConvertFrom-Json
    Write-Success "AWS credentials configured for user: $($awsIdentity.Arn)"
}
catch {
    Write-Error "AWS credentials not configured or invalid"
    exit 1
}

# Update kubeconfig for EKS
Write-Info "Updating kubeconfig for EKS cluster..."
aws eks update-kubeconfig --region $AWSRegion --name $ClusterName
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to update kubeconfig for EKS cluster"
    exit 1
}
Write-Success "Kubeconfig updated"

# Test cluster connectivity
Write-Info "Testing cluster connectivity..."
$nodes = kubectl get nodes --output=json | ConvertFrom-Json
if ($LASTEXITCODE -ne 0) {
    Write-Error "Cannot connect to EKS cluster"
    exit 1
}
Write-Success "Connected to EKS cluster with $($nodes.items.Count) nodes"

# Get ECR registry URL
Write-Info "Getting ECR registry URL..."
aws ecr get-login-password --region $AWSRegion | docker login --username AWS --password-stdin $(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWSRegion.amazonaws.com | Out-Null
$registryUrl = "$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWSRegion.amazonaws.com"
Write-Success "ECR registry URL: $registryUrl"

# Create namespace if it doesn't exist
Write-Header "Setting up Namespace"
kubectl create namespace $Namespace --dry-run=client -o yaml | kubectl apply -f -
Write-Success "Namespace $Namespace is ready"

# Deploy environment-specific configuration
Write-Header "Deploying Environment Configuration"
$configFile = "k8s/aws-$Environment-config.yaml"
if (Test-Path $configFile) {
    # Replace placeholders with actual values from AWS Secrets Manager
    Write-Info "Retrieving database credentials from AWS Secrets Manager..."
    
    try {
        $secretName = "lugx-gaming-postgres-$Environment"
        $secret = aws secretsmanager get-secret-value --secret-id $secretName --query SecretString --output text | ConvertFrom-Json
        
        # Update config file with actual values
        $configContent = Get-Content $configFile -Raw
        $configContent = $configContent -replace "REPLACE_WITH_RDS_ENDPOINT", $secret.host
        $configContent = $configContent -replace "REPLACE_WITH_BASE64_USERNAME", [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes($secret.username))
        $configContent = $configContent -replace "REPLACE_WITH_BASE64_PASSWORD", [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes($secret.password))
        
        # Apply updated configuration
        $configContent | kubectl apply -f -
        Write-Success "Environment configuration deployed"
    }
    catch {
        Write-Warning "Could not retrieve AWS secrets, using placeholder configuration"
        kubectl apply -f $configFile
    }
}
else {
    Write-Warning "Environment config file not found: $configFile"
}

# Update image tags in Kubernetes manifests
Write-Header "Updating Image Tags"
$services = @('frontend', 'game-service', 'order-service', 'analytics-service')

foreach ($service in $services) {
    $manifestFile = "k8s/$service.yaml"
    if (Test-Path $manifestFile) {
        Write-Info "Updating $service image tag..."
        
        # Read manifest content
        $manifestContent = Get-Content $manifestFile -Raw
        
        # Update namespace
        $manifestContent = $manifestContent -replace "namespace: lugx-gaming", "namespace: $Namespace"
        
        # Update image tag
        $imageNameMap = @{
            'frontend'          = 'lugx-gaming-frontend'
            'game-service'      = 'lugx-gaming-game-service'
            'order-service'     = 'lugx-gaming-order-service'
            'analytics-service' = 'lugx-gaming-analytics-service'
        }
        
        $imageName = $imageNameMap[$service]
        $fullImageName = "$registryUrl/$imageName`:$ImageTag"
        
        # Replace image reference
        $manifestContent = $manifestContent -replace "$imageName`:latest", $fullImageName
        $manifestContent = $manifestContent -replace "lugx-$service`:latest", $fullImageName
        
        # Apply updated manifest
        $manifestContent | kubectl apply -f -
        Write-Success "$service manifest updated and applied"
    }
}

# Deploy ClickHouse if not exists
Write-Header "Deploying ClickHouse"
$clickhouseManifest = "k8s/clickhouse.yaml"
if (Test-Path $clickhouseManifest) {
    # Update namespace in ClickHouse manifest
    $clickhouseContent = Get-Content $clickhouseManifest -Raw
    $clickhouseContent = $clickhouseContent -replace "namespace: lugx-gaming", "namespace: $Namespace"
    $clickhouseContent | kubectl apply -f -
    Write-Success "ClickHouse deployed"
}

# Perform rolling deployment with health checks
Write-Header "Performing Rolling Deployment"

foreach ($service in $services) {
    Write-Info "Deploying $service with rolling update..."
    
    # Wait for deployment to be ready
    kubectl rollout status deployment/$service -n $Namespace --timeout=600s
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Deployment of $service failed"
        
        if ($RollbackOnFailure) {
            Write-Warning "Rolling back $service deployment..."
            kubectl rollout undo deployment/$service -n $Namespace
            kubectl rollout status deployment/$service -n $Namespace --timeout=300s
        }
        exit 1
    }
    
    Write-Success "$service deployed successfully"
}

# Wait for all pods to be ready
Write-Header "Waiting for All Pods to be Ready"
kubectl wait --for=condition=ready pod --all -n $Namespace --timeout=600s
if ($LASTEXITCODE -ne 0) {
    Write-Error "Some pods failed to become ready"
    
    Write-Info "Pod status:"
    kubectl get pods -n $Namespace
    
    if ($RollbackOnFailure) {
        Write-Warning "Rolling back all deployments..."
        foreach ($service in $services) {
            kubectl rollout undo deployment/$service -n $Namespace
        }
    }
    exit 1
}

Write-Success "All pods are ready"

# Verify deployment health
Write-Header "Verifying Deployment Health"
$deployments = kubectl get deployments -n $Namespace -o json | ConvertFrom-Json

foreach ($deployment in $deployments.items) {
    $name = $deployment.metadata.name
    $ready = $deployment.status.readyReplicas
    $desired = $deployment.status.replicas
    
    if ($ready -eq $desired) {
        Write-Success "$name`: $ready/$desired replicas ready"
    }
    else {
        Write-Error "$name`: $ready/$desired replicas ready"
        exit 1
    }
}

# Get service endpoints
Write-Header "Service Endpoints"
$services = kubectl get services -n $Namespace -o json | ConvertFrom-Json

foreach ($service in $services.items) {
    $name = $service.metadata.name
    $type = $service.spec.type
    
    if ($type -eq "LoadBalancer") {
        $ingress = $service.status.loadBalancer.ingress
        if ($ingress) {
            $endpoint = if ($ingress[0].hostname) { $ingress[0].hostname } else { $ingress[0].ip }
            Write-Info "$name (LoadBalancer): http://$endpoint"
        }
        else {
            Write-Info "$name (LoadBalancer): Pending external IP..."
        }
    }
    elseif ($type -eq "NodePort") {
        $nodePort = $service.spec.ports[0].nodePort
        Write-Info "$name (NodePort): http://<node-ip>:$nodePort"
    }
}

# Run integration tests if not skipped
if (-not $SkipTests) {
    Write-Header "Running Integration Tests"
    
    # Set up port forwarding for testing
    Write-Info "Setting up port forwarding for integration tests..."
    
    $portForwards = @()
    $portForwards += Start-Process kubectl -ArgumentList "port-forward", "service/game-service", "8080:8080", "-n", $Namespace -PassThru
    $portForwards += Start-Process kubectl -ArgumentList "port-forward", "service/order-service", "8081:8081", "-n", $Namespace -PassThru
    $portForwards += Start-Process kubectl -ArgumentList "port-forward", "service/analytics-service", "8082:8082", "-n", $Namespace -PassThru
    
    Start-Sleep -Seconds 10
    
    try {
        # Set environment variables for tests
        $env:GAME_SERVICE_URL = "http://localhost:8080"
        $env:ORDER_SERVICE_URL = "http://localhost:8081"
        $env:ANALYTICS_SERVICE_URL = "http://localhost:8082"
        
        # Run integration tests
        Write-Info "Running Game Service integration tests..."
        Push-Location "integration-tests/game-service"
        go test -v ./... -timeout=5m
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Game Service integration tests failed"
            throw "Integration tests failed"
        }
        Pop-Location
        
        Write-Info "Running Order Service integration tests..."
        Push-Location "integration-tests/order-service"
        go test -v ./... -timeout=5m
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Order Service integration tests failed"
            throw "Integration tests failed"
        }
        Pop-Location
        
        Write-Info "Running Analytics Service integration tests..."
        Push-Location "integration-tests/analytics-service"
        go test -v ./... -timeout=5m
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Analytics Service integration tests failed"
            throw "Integration tests failed"
        }
        Pop-Location
        
        Write-Success "All integration tests passed"
    }
    catch {
        Write-Error "Integration tests failed: $_"
        
        if ($RollbackOnFailure) {
            Write-Warning "Rolling back deployments due to test failures..."
            foreach ($service in @('frontend', 'game-service', 'order-service', 'analytics-service')) {
                kubectl rollout undo deployment/$service -n $Namespace
            }
        }
        
        exit 1
    }
    finally {
        # Clean up port forwards
        Write-Info "Cleaning up port forwards..."
        foreach ($process in $portForwards) {
            if (-not $process.HasExited) {
                $process.Kill()
            }
        }
    }
}

# Final deployment summary
Write-Header "Deployment Summary"
Write-Success "Deployment to $Environment completed successfully!"
Write-Info "Environment: $Environment"
Write-Info "Namespace: $Namespace"
Write-Info "Image Tag: $ImageTag"
Write-Info "Total time: $((Get-Date) - $startTime)"

Write-Info "To access the application:"
Write-Info "  kubectl get services -n $Namespace"
Write-Info "  kubectl port-forward service/frontend 3000:80 -n $Namespace"

Write-Success "🎉 LUGX Gaming deployment completed successfully!"
