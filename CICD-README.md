# LUGX Gaming - AWS EKS CI/CD Pipeline

This repository contains a comprehensive CI/CD pipeline for deploying the LUGX Gaming platform to AWS EKS with rolling releases, integration testing, and AWS managed PostgreSQL.

## 🚀 Features

- **Multi-Environment Support**: Development, Staging, and Production environments
- **Rolling Deployments**: Zero-downtime deployments with automatic rollback
- **Integration Testing**: Automated post-deployment testing
- **AWS Managed PostgreSQL**: RDS PostgreSQL for production-grade database
- **Container Registry**: AWS ECR for secure image storage
- **Infrastructure as Code**: Terraform for AWS infrastructure management
- **Monitoring & Alerting**: CloudWatch integration with Slack notifications
- **Resource Cleanup**: Automated cleanup of old images and resources

## 📋 Prerequisites

### AWS Setup

1. **AWS Account** with sufficient permissions
2. **AWS CLI** configured with appropriate credentials
3. **Terraform** for infrastructure provisioning
4. **kubectl** for Kubernetes management

### Required AWS Services

- **Amazon EKS** - Kubernetes cluster
- **Amazon RDS** - PostgreSQL database
- **Amazon ECR** - Container registry
- **AWS Secrets Manager** - Credential management
- **Application Load Balancer** - Load balancing
- **CloudWatch** - Monitoring and logging

### GitHub Secrets

Configure the following secrets in your GitHub repository:

```bash
# AWS Credentials
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY

# Database Credentials (will be automatically managed by Terraform)
AWS_RDS_ENDPOINT
AWS_RDS_USERNAME
AWS_RDS_PASSWORD
AWS_RDS_DATABASE

# Optional: Slack Integration
SLACK_WEBHOOK_URL
```

## 🏗️ Infrastructure Setup

### 1. Deploy Infrastructure

Use the Infrastructure workflow to create AWS resources:

```bash
# Via GitHub Actions
Go to Actions → Infrastructure Setup → Run workflow
- Environment: production
- Action: create

# Or manually with Terraform
cd infrastructure
terraform init
terraform apply
```

### 2. Configure EKS Cluster

The Terraform will create:

- VPC with public/private subnets
- EKS cluster with managed node groups
- RDS PostgreSQL instance
- ECR repositories for each service
- Application Load Balancer
- Security groups and IAM roles

## 🔄 CI/CD Workflows

### Main Deployment Workflow (`deploy-eks.yml`)

**Triggers:**

- Push to `main` (production)
- Push to `develop` (staging)
- Push to tags `v*` (production)
- Manual dispatch

**Process:**

1. **Build & Test**: Unit tests for all services
2. **Build Images**: Docker images pushed to ECR
3. **Deploy**: Rolling deployment to EKS
4. **Integration Tests**: Post-deployment testing
5. **Rollback**: Automatic rollback on failure

### Infrastructure Management (`infrastructure.yml`)

**Triggers:**

- Manual dispatch only

**Process:**

- Create/destroy AWS infrastructure
- Manage Terraform state
- Output configuration values

### Cleanup Workflow (`cleanup.yml`)

**Triggers:**

- Weekly schedule (Sundays 2 AM)
- Manual dispatch

**Process:**

- Remove old ECR images
- Clean failed Kubernetes resources
- Manage AWS Secrets lifecycle
- CloudWatch log retention

## 🎯 Deployment Environments

### Development (`lugx-gaming-dev`)

- **Database**: Shared RDS instance with dev schema
- **Resources**: 1 CPU, 2Gi Memory limit
- **Purpose**: Feature development and testing

### Staging (`lugx-gaming-staging`)

- **Database**: Shared RDS instance with staging schema
- **Resources**: 2 CPU, 4Gi Memory limit
- **Purpose**: Pre-production testing and validation

### Production (`lugx-gaming-prod`)

- **Database**: Dedicated RDS instance
- **Resources**: 4 CPU, 8Gi Memory limit
- **Security**: Network policies, resource quotas
- **Monitoring**: Enhanced logging and alerting

## 🔧 Local Development with AWS Integration

### Setup Local Environment

1. **Configure AWS CLI**:

   ```bash
   aws configure
   aws eks update-kubeconfig --region us-east-1 --name lugx-gaming-cluster
   ```

2. **Deploy to Development**:

   ```powershell
   .\scripts\deploy-aws-eks.ps1 -Environment development -ImageTag latest
   ```

3. **Run Integration Tests**:
   ```powershell
   .\integration-tests\run-tests.ps1 -Service all
   ```

### Testing Database Connection

```bash
# Connect to RDS PostgreSQL
psql -h your-rds-endpoint.amazonaws.com -U postgres -d lugx_gaming

# Port forward to development database
kubectl port-forward service/game-service 8080:8080 -n lugx-gaming-dev
curl http://localhost:8080/health
```

## 📊 Monitoring and Observability

### CloudWatch Integration

- **Container Insights**: Automatic metrics collection
- **Application Logs**: Centralized logging
- **Custom Metrics**: Application-specific monitoring

### Prometheus & Grafana (Optional)

Deploy monitoring stack:

```bash
kubectl apply -f k8s/monitoring/
```

Access dashboards:

- **Prometheus**: `http://localhost:30090`
- **Grafana**: `http://localhost:30300` (admin/admin)

## 🚨 Troubleshooting

### Common Issues

1. **Deployment Failures**

   ```bash
   # Check pod status
   kubectl get pods -n lugx-gaming-prod

   # View deployment logs
   kubectl logs deployment/game-service -n lugx-gaming-prod

   # Rollback if needed
   kubectl rollout undo deployment/game-service -n lugx-gaming-prod
   ```

2. **Database Connection Issues**

   ```bash
   # Check secrets
   kubectl get secrets -n lugx-gaming-prod

   # Verify RDS connectivity
   kubectl run -it --rm debug --image=postgres:15 --restart=Never -- \
     psql -h your-rds-endpoint -U postgres -d lugx_gaming
   ```

3. **Image Pull Errors**

   ```bash
   # Check ECR login
   aws ecr get-login-password --region us-east-1 | \
     docker login --username AWS --password-stdin your-account.dkr.ecr.us-east-1.amazonaws.com

   # Verify image exists
   aws ecr describe-images --repository-name lugx-gaming-frontend
   ```

### Emergency Procedures

1. **Immediate Rollback**

   ```bash
   # Rollback all services
   kubectl rollout undo deployment/frontend -n lugx-gaming-prod
   kubectl rollout undo deployment/game-service -n lugx-gaming-prod
   kubectl rollout undo deployment/order-service -n lugx-gaming-prod
   kubectl rollout undo deployment/analytics-service -n lugx-gaming-prod
   ```

2. **Scale Down Services**

   ```bash
   # Emergency scale down
   kubectl scale deployment/frontend --replicas=0 -n lugx-gaming-prod
   kubectl scale deployment/game-service --replicas=0 -n lugx-gaming-prod
   ```

3. **Database Failover**
   ```bash
   # Switch to read replica (manual process)
   aws rds promote-read-replica --db-instance-identifier lugx-gaming-postgres-prod-replica
   ```

## 🔐 Security Best Practices

### Container Security

- Multi-stage Docker builds for minimal attack surface
- Regular image scanning with ECR
- Non-root container execution
- Read-only root filesystems where possible

### Network Security

- VPC with private subnets for EKS nodes
- Security groups with minimal required access
- Network policies for pod-to-pod communication
- SSL/TLS termination at load balancer

### Secret Management

- AWS Secrets Manager for database credentials
- Kubernetes secrets for application configuration
- Rotation policies for long-lived credentials
- Encryption at rest and in transit

## 📈 Scaling and Performance

### Horizontal Pod Autoscaling

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: frontend-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: frontend
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

### Cluster Autoscaling

EKS managed node groups automatically scale based on pod resource requirements.

### Database Scaling

- RDS supports vertical scaling (instance size)
- Read replicas for read-heavy workloads
- Connection pooling in applications

## 🔄 Backup and Disaster Recovery

### Database Backups

- Automated RDS backups (7-day retention for production)
- Manual snapshots before major deployments
- Cross-region backup replication for production

### Application State

- Stateless application design
- External state storage (RDS, S3)
- Configuration stored in Kubernetes ConfigMaps/Secrets

### Disaster Recovery Plan

1. **RTO**: 15 minutes for application recovery
2. **RPO**: 1 hour for database recovery
3. **Procedures**: Documented runbooks for common scenarios

## 📝 Development Workflow

### Feature Development

1. Create feature branch from `develop`
2. Develop and test locally
3. Create pull request to `develop`
4. Automatic deployment to staging on merge
5. Manual promotion to production

### Release Process

1. Merge `develop` to `main` for release
2. Tag release with semantic version
3. Automatic deployment to production
4. Monitor deployment and run smoke tests

## 🤝 Contributing

1. Follow the branching strategy
2. Ensure all tests pass
3. Update documentation as needed
4. Follow security best practices
5. Test in staging before production

## 📞 Support

For issues with the CI/CD pipeline:

1. Check GitHub Actions logs
2. Review CloudWatch logs
3. Consult troubleshooting guide
4. Contact DevOps team

---

**Last Updated**: $(Get-Date -Format "yyyy-MM-dd")
**Version**: 1.0.0
