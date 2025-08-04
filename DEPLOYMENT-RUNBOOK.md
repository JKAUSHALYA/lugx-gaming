# LUGX Gaming Platform - Deployment Runbook

**Version**: 1.0.0  
**Last Updated**: $(Get-Date -Format "yyyy-MM-dd")  
**Author**: DevOps Team

---

## 📋 Overview

This runbook provides step-by-step instructions for deploying and testing the LUGX Gaming platform across different environments (local, Kubernetes, and AWS EKS).

## 🎯 Deployment Options

- **Local Development**: Docker Compose for rapid development
- **Kubernetes**: Local or cloud Kubernetes cluster
- **AWS EKS**: Production-ready AWS deployment with managed services

---

## 🔧 Prerequisites

### Required Tools

```powershell
# Check prerequisites
Write-Host "Checking prerequisites..." -ForegroundColor Yellow

# Docker
docker --version
if ($LASTEXITCODE -ne 0) { Write-Host "❌ Docker not installed" -ForegroundColor Red; exit 1 }

# kubectl (for Kubernetes deployments)
kubectl version --client
if ($LASTEXITCODE -ne 0) { Write-Host "⚠️ kubectl not installed (required for K8s)" -ForegroundColor Yellow }

# AWS CLI (for EKS deployments)
aws --version
if ($LASTEXITCODE -ne 0) { Write-Host "⚠️ AWS CLI not installed (required for EKS)" -ForegroundColor Yellow }

# Go (for local service development)
go version
if ($LASTEXITCODE -ne 0) { Write-Host "⚠️ Go not installed (required for local dev)" -ForegroundColor Yellow }
```

### Environment Variables

```powershell
# For AWS deployments
$env:AWS_ACCESS_KEY_ID = "your-access-key"
$env:AWS_SECRET_ACCESS_KEY = "your-secret-key"
$env:AWS_DEFAULT_REGION = "us-east-1"

# Database configuration (for local/K8s)
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_NAME = "lugx_gaming"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "password"
```

---

## 🚀 Deployment Procedures

## Option 1: Local Development Deployment

### Step 1: Start PostgreSQL Database

```powershell
# Navigate to project root
cd "e:\Learning\MSc\CMM 707 - Cloud Computing\CW\lugx-gaming"

# Start PostgreSQL
cd postgres
.\start.ps1

# Verify database is running
docker ps | findstr postgres
```

**Expected Output:**

```
✓ Network 'lugx-network' created successfully
✓ PostgreSQL started successfully
✓ PostgreSQL is ready and accepting connections
```

### Step 2: Start Backend Services

```powershell
# Option 2a: Start services individually
# Terminal 1 - Game Service
cd game-service
go mod tidy
go run main.go

# Terminal 2 - Order Service
cd order-service
go mod tidy
go run main.go

# Terminal 3 - Analytics Service
cd analytics-service
go mod tidy
go run main.go
```

**OR use VS Code tasks:**

```powershell
# Using VS Code tasks (recommended)
# Press Ctrl+Shift+P, type "Tasks: Run Task"
# Select "Start Order Service"
# Select "Start Analytics Service"
```

### Step 3: Start Frontend

```powershell
cd front-end

# Option 3a: Docker Compose (recommended)
docker-compose up -d --build

# Option 3b: Docker manually
docker build -t lugx-gaming-frontend .
docker run -d --name lugx-gaming-frontend -p 3000:80 lugx-gaming-frontend
```

### Step 4: Verify Local Deployment

```powershell
# Check service health
curl http://localhost:8080/health  # Game Service
curl http://localhost:8081/health  # Order Service
curl http://localhost:8082/health  # Analytics Service

# Check frontend
Start-Process "http://localhost:3000"
```

---

## Option 2: Kubernetes Deployment

### Step 1: Prepare Kubernetes Cluster

```powershell
# For local Kubernetes (Docker Desktop/minikube)
kubectl cluster-info

# Create namespace
kubectl create namespace lugx-gaming
kubectl config set-context --current --namespace=lugx-gaming
```

### Step 2: Deploy External Services

```powershell
cd k8s

# Deploy PostgreSQL and other external services
kubectl apply -f external-services.yaml

# Wait for PostgreSQL to be ready
kubectl wait --for=condition=ready pod -l app=postgres --timeout=300s
```

### Step 3: Build and Deploy Application

```powershell
# Build Docker images and deploy
.\deploy.ps1

# Monitor deployment
kubectl get pods -w
```

**Expected Output:**

```
🚀 LUGX Gaming Kubernetes Deployment
=====================================
🔨 Building lugx-gaming-frontend...
✅ Successfully built lugx-gaming-frontend
🔨 Building lugx-game-service...
✅ Successfully built lugx-game-service
...
✅ All services deployed successfully!
```

### Step 4: Verify Kubernetes Deployment

```powershell
# Check all pods are running
kubectl get pods

# Check services
kubectl get services

# Port forward to access services locally
kubectl port-forward service/frontend 3000:80
kubectl port-forward service/game-service 8080:8080
kubectl port-forward service/order-service 8081:8080
kubectl port-forward service/analytics-service 8082:8080
```

---

## Option 3: AWS EKS Deployment

### Step 1: Configure AWS Access

```powershell
# Configure AWS CLI
aws configure

# Verify access
aws sts get-caller-identity

# Update kubeconfig for EKS
aws eks update-kubeconfig --region us-east-1 --name lugx-gaming-cluster
```

### Step 2: Deploy Infrastructure (if not exists)

```powershell
# Deploy AWS infrastructure using Terraform
cd infrastructure
terraform init
terraform plan
terraform apply

# Note: This creates EKS cluster, RDS, ECR repos, etc.
```

### Step 3: Deploy Application to EKS

```powershell
cd scripts

# Deploy to development environment
.\deploy-aws-eks.ps1 -Environment development -ImageTag latest

# Deploy to staging environment
.\deploy-aws-eks.ps1 -Environment staging -ImageTag v1.0.0

# Deploy to production environment
.\deploy-aws-eks.ps1 -Environment production -ImageTag v1.0.0
```

**Expected Output:**

```
🚀 Starting AWS EKS deployment...
✅ AWS credentials verified
✅ EKS cluster connection established
🔨 Building and pushing images to ECR...
✅ Images pushed successfully
🚀 Deploying to EKS...
✅ Deployment completed successfully
🧪 Running integration tests...
✅ All tests passed
```

### Step 4: Verify EKS Deployment

```powershell
# Check deployment status
kubectl get pods -n lugx-gaming-prod
kubectl get services -n lugx-gaming-prod

# Get load balancer URL
kubectl get service frontend -n lugx-gaming-prod -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'

# Check RDS connection
kubectl exec -it deployment/game-service -n lugx-gaming-prod -- nc -zv your-rds-endpoint 5432
```

---

## 🧪 Testing Procedures

### Step 1: Health Checks

```powershell
# Local/Kubernetes health checks
$services = @(
    @{Name="Game Service"; URL="http://localhost:8080/health"},
    @{Name="Order Service"; URL="http://localhost:8081/health"},
    @{Name="Analytics Service"; URL="http://localhost:8082/health"}
)

foreach ($service in $services) {
    try {
        $response = Invoke-RestMethod -Uri $service.URL -Method GET -TimeoutSec 10
        Write-Host "✅ $($service.Name): $($response.status)" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ $($service.Name): Failed" -ForegroundColor Red
    }
}
```

### Step 2: Integration Tests

```powershell
cd integration-tests

# Run all integration tests
.\run-tests.ps1

# Run specific service tests
.\run-tests.ps1 -Service game -VerboseOutput
.\run-tests.ps1 -Service order -VerboseOutput
.\run-tests.ps1 -Service analytics -VerboseOutput
```

**Expected Output:**

```
🧪 Running Integration Tests for LUGX Gaming Platform
=====================================================
✅ Game Service: All 12 tests passed
✅ Order Service: All 8 tests passed
✅ Analytics Service: All 6 tests passed
🎉 All integration tests completed successfully!
```

### Step 3: End-to-End Testing

```powershell
# Test complete user workflow
$baseUrl = "http://localhost:3000"  # or load balancer URL for EKS

# 1. Frontend accessibility
$response = Invoke-WebRequest -Uri $baseUrl
if ($response.StatusCode -eq 200) {
    Write-Host "✅ Frontend accessible" -ForegroundColor Green
} else {
    Write-Host "❌ Frontend not accessible" -ForegroundColor Red
}

# 2. Game catalog loading
$gameService = "http://localhost:8080/api/games"
$games = Invoke-RestMethod -Uri $gameService
if ($games.Count -gt 0) {
    Write-Host "✅ Game catalog loaded: $($games.Count) games" -ForegroundColor Green
} else {
    Write-Host "❌ Game catalog empty" -ForegroundColor Red
}

# 3. Order creation test
$orderService = "http://localhost:8081/api/orders"
$orderData = @{
    user_id = "test-user"
    game_id = $games[0].id
    quantity = 1
} | ConvertTo-Json

try {
    $order = Invoke-RestMethod -Uri $orderService -Method POST -Body $orderData -ContentType "application/json"
    Write-Host "✅ Order created: ID $($order.id)" -ForegroundColor Green
} catch {
    Write-Host "❌ Order creation failed: $($_.Exception.Message)" -ForegroundColor Red
}
```

### Step 4: Performance Testing

```powershell
# Simple load test using curl
$endpoints = @(
    "http://localhost:8080/api/games",
    "http://localhost:8081/health",
    "http://localhost:8082/health"
)

foreach ($endpoint in $endpoints) {
    Write-Host "Testing $endpoint..." -ForegroundColor Yellow

    $startTime = Get-Date
    for ($i = 1; $i -le 10; $i++) {
        try {
            Invoke-RestMethod -Uri $endpoint -TimeoutSec 5 | Out-Null
        } catch {
            Write-Host "Request $i failed" -ForegroundColor Red
        }
    }
    $endTime = Get-Date
    $duration = ($endTime - $startTime).TotalMilliseconds

    Write-Host "✅ 10 requests completed in $([math]::Round($duration, 2))ms" -ForegroundColor Green
    Write-Host "   Average: $([math]::Round($duration/10, 2))ms per request" -ForegroundColor Cyan
}
```

---

## 🔍 Monitoring and Verification

### Step 1: Deploy Monitoring Stack (Optional)

```powershell
cd k8s

# Deploy Prometheus and Grafana
.\deploy-monitoring.ps1

# Test monitoring stack
.\test-monitoring.ps1
```

### Step 2: Check Application Metrics

```powershell
# Kubernetes metrics
kubectl top pods -n lugx-gaming
kubectl top nodes

# Application logs
kubectl logs -f deployment/game-service -n lugx-gaming
kubectl logs -f deployment/order-service -n lugx-gaming
kubectl logs -f deployment/analytics-service -n lugx-gaming
```

### Step 3: Database Verification

```powershell
# For local PostgreSQL
docker exec -it postgres-container psql -U postgres -d lugx_gaming -c "\dt"

# For Kubernetes PostgreSQL
kubectl exec -it deployment/postgres -- psql -U postgres -d lugx_gaming -c "\dt"

# For AWS RDS (via port forward)
kubectl port-forward service/game-service 5432:5432 -n lugx-gaming-prod
psql -h localhost -U postgres -d lugx_gaming -c "\dt"
```

---

## 🚨 Troubleshooting

### Common Issues and Solutions

#### 1. Pod Not Starting

```powershell
# Check pod status
kubectl describe pod <pod-name>

# Check logs
kubectl logs <pod-name> --previous

# Common fixes:
# - Image pull errors: Check ECR permissions
# - Resource limits: Increase CPU/memory
# - Config errors: Verify ConfigMaps/Secrets
```

#### 2. Database Connection Issues

```powershell
# Test database connectivity
kubectl run -it --rm debug --image=postgres:15 --restart=Never -- \
  psql -h postgres-service -U postgres -d lugx_gaming

# Check service discovery
kubectl get endpoints
nslookup postgres-service
```

#### 3. Service Unavailable

```powershell
# Check service and endpoints
kubectl get service <service-name>
kubectl get endpoints <service-name>

# Port forward for debugging
kubectl port-forward service/<service-name> 8080:8080

# Test service directly
curl http://localhost:8080/health
```

### Rollback Procedures

#### Quick Rollback

```powershell
# Rollback specific deployment
kubectl rollout undo deployment/<deployment-name>

# Rollback all services
$deployments = @("frontend", "game-service", "order-service", "analytics-service")
foreach ($deployment in $deployments) {
    kubectl rollout undo deployment/$deployment
    kubectl rollout status deployment/$deployment
}
```

#### Emergency Scale Down

```powershell
# Scale down all services
kubectl scale deployment/frontend --replicas=0
kubectl scale deployment/game-service --replicas=0
kubectl scale deployment/order-service --replicas=0
kubectl scale deployment/analytics-service --replicas=0
```

---

## 🧹 Cleanup Procedures

### Local Environment Cleanup

```powershell
# Stop and remove containers
docker-compose down
docker container prune -f
docker image prune -f

# Remove Docker network
docker network rm lugx-network
```

### Kubernetes Cleanup

```powershell
cd k8s

# Use cleanup script
.\cleanup.ps1

# Or manual cleanup
kubectl delete namespace lugx-gaming
```

### AWS EKS Cleanup

```powershell
# Cleanup monitoring stack
cd k8s
.\cleanup-monitoring.ps1

# Cleanup application
.\cleanup.ps1

# Destroy infrastructure (CAUTION: This will delete everything)
cd infrastructure
terraform destroy
```

---

## 📊 Deployment Checklist

### Pre-Deployment

- [ ] Prerequisites installed and configured
- [ ] Environment variables set
- [ ] Database accessible
- [ ] Docker images built successfully
- [ ] Kubernetes cluster accessible (if applicable)
- [ ] AWS credentials configured (if applicable)

### During Deployment

- [ ] All services start successfully
- [ ] Health checks pass
- [ ] Database connections established
- [ ] Load balancer configured (for cloud deployments)
- [ ] DNS/ingress configured

### Post-Deployment

- [ ] Integration tests pass
- [ ] End-to-end testing completed
- [ ] Monitoring alerts configured
- [ ] Performance benchmarks met
- [ ] Security scans completed
- [ ] Documentation updated

### Rollback Criteria

- [ ] Any service fails health checks for >5 minutes
- [ ] Integration tests fail
- [ ] Database connectivity issues
- [ ] Performance degradation >50%
- [ ] Security vulnerabilities detected

---

## 📞 Support and Escalation

### Support Contacts

- **DevOps Team**: devops@lugxgaming.com
- **Database Team**: dba@lugxgaming.com
- **Security Team**: security@lugxgaming.com

### Escalation Matrix

1. **Level 1**: Service restart, basic troubleshooting
2. **Level 2**: Configuration changes, rollback procedures
3. **Level 3**: Infrastructure changes, emergency procedures

### Emergency Procedures

1. **Immediate Response**: Scale down affected services
2. **Communication**: Notify stakeholders via Slack/Teams
3. **Investigation**: Gather logs and metrics
4. **Resolution**: Apply fix or rollback
5. **Post-Mortem**: Document lessons learned

---

## 📝 Change Log

| Version | Date       | Changes                  |
| ------- | ---------- | ------------------------ |
| 1.0.0   | 2025-08-05 | Initial runbook creation |

---

**End of Runbook**

For additional support, refer to the [CICD-README.md](./CICD-README.md) for detailed pipeline information.
