#!/usr/bin/env pwsh

<#
.SYNOPSIS
    Quick setup script for LUGX Gaming AWS EKS CI/CD pipeline
.DESCRIPTION
    This script helps set up the initial configuration for deploying
    LUGX Gaming platform to AWS EKS.
.PARAMETER AWSRegion
    AWS region to deploy to
.PARAMETER Environment
    Environment to set up (development, staging, production)
.EXAMPLE
    .\setup-aws-cicd.ps1 -Environment development -AWSRegion us-east-1
#>

param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('development', 'staging', 'production')]
    [string]$Environment,
    
    [Parameter(Mandatory = $false)]
    [string]$AWSRegion = "us-east-1"
)

function Write-Success($message) { Write-Host "✅ $message" -ForegroundColor Green }
function Write-Error($message) { Write-Host "❌ $message" -ForegroundColor Red }
function Write-Info($message) { Write-Host "ℹ️  $message" -ForegroundColor Cyan }
function Write-Warning($message) { Write-Host "⚠️  $message" -ForegroundColor Yellow }
function Write-Header($message) {
    Write-Host ""
    Write-Host "🚀 $message" -ForegroundColor Magenta
    Write-Host "=" * 50 -ForegroundColor Magenta
}

Write-Header "LUGX Gaming AWS EKS Setup"
Write-Info "Environment: $Environment"
Write-Info "AWS Region: $AWSRegion"

# Check prerequisites
Write-Header "Checking Prerequisites"

$requiredTools = @{
    'aws'       = 'AWS CLI'
    'kubectl'   = 'Kubernetes CLI'
    'terraform' = 'Terraform'
    'docker'    = 'Docker'
}

$allToolsAvailable = $true
foreach ($tool in $requiredTools.Keys) {
    if (Get-Command $tool -ErrorAction SilentlyContinue) {
        Write-Success "$($requiredTools[$tool]) found"
    }
    else {
        Write-Error "$($requiredTools[$tool]) not found"
        $allToolsAvailable = $false
    }
}

if (-not $allToolsAvailable) {
    Write-Error "Please install missing tools before continuing"
    exit 1
}

# Check AWS credentials
Write-Header "Verifying AWS Credentials"
try {
    $awsIdentity = aws sts get-caller-identity --output json | ConvertFrom-Json
    Write-Success "AWS credentials configured for: $($awsIdentity.Arn)"
    Write-Info "Account ID: $($awsIdentity.Account)"
}
catch {
    Write-Error "AWS credentials not configured. Run 'aws configure' first."
    exit 1
}

# Create directories if they don't exist
Write-Header "Setting Up Directory Structure"
$directories = @(
    "scripts",
    "infrastructure",
    ".github/workflows"
)

foreach ($dir in $directories) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
        Write-Success "Created directory: $dir"
    }
    else {
        Write-Info "Directory exists: $dir"
    }
}

# Create environment-specific configuration
Write-Header "Creating Environment Configuration"

$namespace = switch ($Environment) {
    'development' { 'lugx-gaming-dev' }
    'staging' { 'lugx-gaming-staging' }
    'production' { 'lugx-gaming-prod' }
}

# Create Terraform variables file
$tfVarsContent = @"
aws_region  = "$AWSRegion"
environment = "$Environment"
cluster_name = "lugx-gaming-cluster"
"@

$tfVarsPath = "infrastructure/terraform.tfvars"
$tfVarsContent | Out-File -FilePath $tfVarsPath -Encoding UTF8
Write-Success "Created Terraform variables file: $tfVarsPath"

# Create environment-specific namespace manifest
$namespaceManifest = @"
apiVersion: v1
kind: Namespace
metadata:
  name: $namespace
  labels:
    environment: $Environment
    project: lugx-gaming
---
apiVersion: v1
kind: ResourceQuota
metadata:
  name: lugx-gaming-quota
  namespace: $namespace
spec:
  hard:
    requests.cpu: "$(if ($Environment -eq 'production') { '4' } elseif ($Environment -eq 'staging') { '2' } else { '1' })"
    requests.memory: "$(if ($Environment -eq 'production') { '8Gi' } elseif ($Environment -eq 'staging') { '4Gi' } else { '2Gi' })"
    limits.cpu: "$(if ($Environment -eq 'production') { '8' } elseif ($Environment -eq 'staging') { '4' } else { '2' })"
    limits.memory: "$(if ($Environment -eq 'production') { '16Gi' } elseif ($Environment -eq 'staging') { '8Gi' } else { '4Gi' })"
    pods: "$(if ($Environment -eq 'production') { '20' } elseif ($Environment -eq 'staging') { '15' } else { '10' })"
    services: "$(if ($Environment -eq 'production') { '10' } elseif ($Environment -eq 'staging') { '8' } else { '5' })"
"@

$namespaceFile = "k8s/namespace-$Environment.yaml"
$namespaceManifest | Out-File -FilePath $namespaceFile -Encoding UTF8
Write-Success "Created namespace manifest: $namespaceFile"

# Check if infrastructure already exists
Write-Header "Checking Existing Infrastructure"

try {
    aws eks describe-cluster --name "lugx-gaming-cluster-$Environment" --region $AWSRegion 2>$null | Out-Null
    if ($LASTEXITCODE -eq 0) {
        Write-Success "EKS cluster already exists"
    }
    else {
        Write-Warning "EKS cluster does not exist - will be created by Terraform"
    }
}
catch {
    Write-Warning "EKS cluster does not exist - will be created by Terraform"
}

# Create ECR repositories
Write-Header "Setting Up ECR Repositories"

$repositories = @(
    "lugx-gaming-frontend",
    "lugx-gaming-game-service", 
    "lugx-gaming-order-service",
    "lugx-gaming-analytics-service"
)

foreach ($repo in $repositories) {
    try {
        aws ecr describe-repositories --repository-names $repo --region $AWSRegion | Out-Null
        Write-Success "ECR repository exists: $repo"
    }
    catch {
        Write-Info "Creating ECR repository: $repo"
        aws ecr create-repository --repository-name $repo --region $AWSRegion | Out-Null
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Created ECR repository: $repo"
        }
        else {
            Write-Error "Failed to create ECR repository: $repo"
        }
    }
}

# Generate GitHub Actions secrets template
Write-Header "Generating GitHub Secrets Template"

$secretsTemplate = @"
# GitHub Repository Secrets Configuration
# Add these secrets to your GitHub repository settings

## AWS Credentials
AWS_ACCESS_KEY_ID=<your-aws-access-key>
AWS_SECRET_ACCESS_KEY=<your-aws-secret-key>

## Will be populated after Terraform deployment
AWS_RDS_ENDPOINT=<rds-endpoint-from-terraform>
AWS_RDS_USERNAME=<rds-username-from-terraform>
AWS_RDS_PASSWORD=<rds-password-from-terraform>
AWS_RDS_DATABASE=lugx_gaming

## Optional: Slack Integration
SLACK_WEBHOOK_URL=<your-slack-webhook-url>

## EKS Configuration
EKS_CLUSTER_NAME=lugx-gaming-cluster-$Environment
AWS_REGION=$AWSRegion
"@

$secretsFile = "github-secrets-template.txt"
$secretsTemplate | Out-File -FilePath $secretsFile -Encoding UTF8
Write-Success "Created GitHub secrets template: $secretsFile"

# Create a simple deployment validation script
Write-Header "Creating Deployment Validation Script"

$validationScript = @"
#!/usr/bin/env pwsh
# Quick deployment validation for $Environment environment

Write-Host "🔍 Validating $Environment deployment..." -ForegroundColor Cyan

# Check namespace
kubectl get namespace $namespace
if (`$LASTEXITCODE -eq 0) {
    Write-Host "✅ Namespace $namespace exists" -ForegroundColor Green
} else {
    Write-Host "❌ Namespace $namespace not found" -ForegroundColor Red
    exit 1
}

# Check deployments
`$deployments = @('frontend', 'game-service', 'order-service', 'analytics-service')
foreach (`$deployment in `$deployments) {
    `$status = kubectl get deployment `$deployment -n $namespace -o jsonpath='{.status.readyReplicas}/{.status.replicas}' 2>`$null
    if (`$LASTEXITCODE -eq 0) {
        Write-Host "✅ `$deployment`: `$status replicas ready" -ForegroundColor Green
    } else {
        Write-Host "❌ `$deployment: not found or not ready" -ForegroundColor Red
    }
}

# Check services
kubectl get services -n $namespace
Write-Host "✅ Deployment validation completed" -ForegroundColor Green
"@

$validationFile = "scripts/validate-$Environment.ps1"
$validationScript | Out-File -FilePath $validationFile -Encoding UTF8
Write-Success "Created validation script: $validationFile"

# Next steps instructions
Write-Header "Next Steps"

Write-Info "1. Infrastructure Setup:"
Write-Host "   terraform init" -ForegroundColor White
Write-Host "   terraform plan" -ForegroundColor White
Write-Host "   terraform apply" -ForegroundColor White

Write-Info "2. Configure GitHub Secrets:"
Write-Host "   - Go to your GitHub repository settings" -ForegroundColor White
Write-Host "   - Add secrets from: $secretsFile" -ForegroundColor White

Write-Info "3. Deploy Application:"
Write-Host "   # Via GitHub Actions (recommended)" -ForegroundColor White
Write-Host "   # Or manually:" -ForegroundColor White
Write-Host "   .\scripts\deploy-aws-eks.ps1 -Environment $Environment -ImageTag latest" -ForegroundColor White

Write-Info "4. Validate Deployment:"
Write-Host "   .\scripts\validate-$Environment.ps1" -ForegroundColor White

Write-Info "5. Access Application:"
Write-Host "   kubectl port-forward service/frontend 3000:80 -n $namespace" -ForegroundColor White
Write-Host "   # Then visit: http://localhost:3000" -ForegroundColor White

Write-Success "🎉 Setup completed for $Environment environment!"
Write-Info "For detailed documentation, see: CICD-README.md"
