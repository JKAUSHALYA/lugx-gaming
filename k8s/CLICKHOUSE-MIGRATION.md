# ClickHouse Migration to Kubernetes

## Summary

Successfully migrated ClickHouse from Docker Compose to Kubernetes management. The analytics service now connects to a ClickHouse instance running within the Kubernetes cluster instead of an external Docker container.

## Changes Made

### 1. Created New ClickHouse Kubernetes Manifest (`k8s/clickhouse.yaml`)

- **ConfigMap**: `clickhouse-config` with database initialization scripts
- **Secret**: `clickhouse-secret` with ClickHouse credentials
- **PersistentVolumeClaim**: `clickhouse-data` with 10Gi storage
- **Deployment**: ClickHouse server with proper resource limits and health checks
- **Services**:
  - Internal service: `clickhouse` (ClusterIP)
  - External service: `clickhouse-external` (NodePort 30123 HTTP, 30900 Native)

### 2. Updated Analytics Service Configuration (`k8s/analytics-service.yaml`)

- Changed `CLICKHOUSE_HOST` from `host.docker.internal` to `clickhouse`
- Changed `CLICKHOUSE_PORT` from external ConfigMap to hardcoded `9000`
- Updated secret references to use `clickhouse-secret` instead of `external-db-secrets`
- Updated init container to use internal service names

### 3. Updated External Services Configuration (`k8s/external-services.yaml`)

- Removed ClickHouse configuration from `external-services-config` ConfigMap
- Removed ClickHouse credentials from `external-db-secrets` Secret
- Kept only PostgreSQL configuration for external database

### 4. Updated Kustomization (`k8s/kustomization.yaml`)

- Added `clickhouse.yaml` to the resources list

### 5. Updated Deployment Script (`k8s/deploy.ps1`)

- Removed Docker Compose startup for ClickHouse
- Added ClickHouse deployment step to Kubernetes deployment
- Updated service URLs to include ClickHouse endpoints

### 6. Updated Documentation (`k8s/README.md`)

- Updated architecture description to reflect internal ClickHouse management
- Removed ClickHouse from external database configuration
- Added ClickHouse service URLs and access information
- Updated deployment steps to include ClickHouse

## Benefits

1. **Unified Management**: All services now managed through Kubernetes
2. **Better Resource Control**: ClickHouse now has proper resource limits and requests
3. **Persistent Storage**: Data persists across pod restarts with PVC
4. **Health Monitoring**: Proper liveness and readiness probes
5. **Scalability**: Easier to scale and manage alongside other services
6. **Service Discovery**: Internal DNS resolution for service-to-service communication

## Service Endpoints

- **ClickHouse HTTP (External)**: `http://localhost:30123`
- **ClickHouse Native (External)**: `localhost:30900`
- **ClickHouse (Internal)**: `clickhouse:8123` (HTTP), `clickhouse:9000` (Native)

## Deployment

Use the existing deployment commands:

```powershell
# Deploy all services including ClickHouse
cd k8s
.\deploy.ps1

# Or deploy manually
kubectl apply -k .
```

## Verification

```powershell
# Check ClickHouse pod status
kubectl get pods -n lugx-gaming -l app=clickhouse

# Check ClickHouse logs
kubectl logs -n lugx-gaming deployment/clickhouse

# Test ClickHouse connectivity
curl http://localhost:30123/ping
```
