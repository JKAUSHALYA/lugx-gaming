# Kubernetes Deployment for LUGX Gaming Platform

This directory contains Kubernetes manifests for deploying the LUGX Gaming platform with four microservices:

- **Frontend**: Static website served by Nginx
- **Game Service**: Go-based REST API for game management (connects to external PostgreSQL)
- **Order Service**: Go-based REST API for order management (connects to external PostgreSQL)
- **Analytics Service**: Go-based REST API for user behavior analytics (connects to Kubernetes-managed ClickHouse)
- **ClickHouse**: Analytics database deployed in Kubernetes cluster

## Prerequisites

1. **Kubernetes Cluster**: Ensure you have a running Kubernetes cluster (minikube, kind, or cloud provider)
2. **kubectl**: Kubernetes command-line tool configured to connect to your cluster
3. **External Databases**:
   - PostgreSQL database accessible from Kubernetes cluster
4. **Docker Images**: Build and tag the following Docker images:
   - `lugx-gaming-frontend:latest`
   - `lugx-game-service:latest`
   - `lugx-order-service:latest`
   - `lugx-analytics-service:latest`

## Database Setup

### PostgreSQL Configuration (External)

Before deploying, ensure PostgreSQL is running and accessible:

- Host: Update `POSTGRES_HOST` in `external-services.yaml`
- Port: Default 5432 (update if different)
- Database: `lugx_gaming` (or update `POSTGRES_DB` in ConfigMap)
- Credentials: Update base64 encoded values in `external-db-secrets`

### ClickHouse Configuration (Kubernetes-managed)

ClickHouse is now deployed as part of the Kubernetes cluster with:

- **Persistent Storage**: 10Gi PVC for data persistence
- **Internal Service**: `clickhouse:9000` (native) and `clickhouse:8123` (HTTP)
- **External Access**: NodePort 30900 (native) and 30123 (HTTP)
- **Auto-initialization**: Database schema and tables created automatically
- **Resource Limits**: 1Gi memory, 500m CPU limits

## Configure External Database Credentials

Before deploying, you need to update the PostgreSQL credentials in `external-services.yaml`:

1. **Encode your credentials to base64**:

   ```powershell
   # For PostgreSQL
   [System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes("your_postgres_username"))
   [System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes("your_postgres_password"))
   ```

2. **Update the Secret in `external-services.yaml`**:
   Replace the base64 encoded values in the `external-db-secrets` Secret section.

3. **Update PostgreSQL host**:
   If your PostgreSQL database is not accessible via `host.docker.internal`, update the host value in the ConfigMap section of `external-services.yaml`.

**Note**: ClickHouse credentials are managed internally by Kubernetes and don't need to be configured manually.

## Build Docker Images

Before deploying to Kubernetes, build the Docker images:

```powershell
# Build Frontend
cd front-end
docker build -t lugx-gaming-frontend:latest .

# Build Game Service
cd ../game-service
docker build -t lugx-game-service:latest .

# Build Order Service
cd ../order-service
docker build -t lugx-order-service:latest .

# Build Analytics Service
cd ../analytics-service
docker build -t lugx-analytics-service:latest .
```

## Deployment Steps

1. **Apply all manifests**:

   ```powershell
   kubectl apply -f k8s/
   ```

2. **Or deploy step by step**:

   ```powershell
   # Create namespace
   kubectl apply -f k8s/namespace.yaml

   # Deploy external services configuration and secrets
   # Configure external services
   kubectl apply -f k8s/external-services.yaml

   # Deploy ClickHouse
   kubectl apply -f k8s/clickhouse.yaml

   # Deploy application services
   kubectl apply -f k8s/game-service.yaml
   kubectl apply -f k8s/order-service.yaml
   kubectl apply -f k8s/analytics-service.yaml
   kubectl apply -f k8s/frontend.yaml
   ```

   Or use Kustomize:

   ```powershell
   kubectl apply -k k8s/
   ```

## Verify Deployment

```powershell
# Check all resources in the namespace
kubectl get all -n lugx-gaming

# Check pod status
kubectl get pods -n lugx-gaming

# Check services
kubectl get services -n lugx-gaming

# View logs
kubectl logs -n lugx-gaming deployment/game-service
kubectl logs -n lugx-gaming deployment/order-service
kubectl logs -n lugx-gaming deployment/analytics-service
kubectl logs -n lugx-gaming deployment/frontend
kubectl logs -n lugx-gaming deployment/clickhouse
```

## Access the Application

The services are exposed via NodePort services:

- **Frontend**: `http://localhost:30000` (or `http://<node-ip>:30000`)
- **Game Service**: `http://localhost:30080` (or `http://<node-ip>:30080`)
- **Order Service**: `http://localhost:30081` (or `http://<node-ip>:30081`)
- **Analytics Service**: `http://localhost:30082` (or `http://<node-ip>:30082`)
- **ClickHouse HTTP**: `http://localhost:30123` (or `http://<node-ip>:30123`)
- **ClickHouse Native**: `localhost:30900` (or `<node-ip>:30900`)

For minikube, you can also use:

```powershell
minikube service frontend-external -n lugx-gaming
minikube service game-service-external -n lugx-gaming
minikube service order-service-external -n lugx-gaming
minikube service analytics-service-external -n lugx-gaming
```

## Architecture

```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│    Frontend     │  │  Game Service   │  │  Order Service  │  │Analytics Service│
│   (Port 80)     │  │   (Port 8080)   │  │   (Port 8081)   │  │   (Port 8080)   │
│   Nginx/HTML    │  │   Go/Gin API    │  │   Go/Gin API    │  │   Go/Mux API    │
└─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────────┘
        │                        │                        │                    │
        │                        └────────┬───────────────┘                    │
        │                                 │                                    │
        │              ┌─────────────────────────────────┐      ┌─────────────────────────────────┐
        │              │         PostgreSQL              │      │         ClickHouse              │
        │              │        (Port 5432)              │      │        (Port 9000)              │
        │              │      Database Backend           │      │      Analytics Database         │
        │              └─────────────────────────────────┘      └─────────────────────────────────┘
        │
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                       │
│  - Namespace: lugx-gaming                                   │
│  - Secrets: Database credentials                            │
│  - PersistentVolume: Database storage                       │
│  - Services: Internal and external access                   │
│  - Deployments: Multi-replica services                      │
└─────────────────────────────────────────────────────────────┘
```

## Scaling

Scale individual services:

```powershell
kubectl scale deployment frontend --replicas=5 -n lugx-gaming
kubectl scale deployment game-service --replicas=3 -n lugx-gaming
kubectl scale deployment order-service --replicas=3 -n lugx-gaming
kubectl scale deployment analytics-service --replicas=3 -n lugx-gaming
```

## Cleanup

Remove all resources:

```powershell
kubectl delete namespace lugx-gaming
```

## Configuration Details

### Resource Limits

- **Frontend**: 64Mi-128Mi memory, 50m-100m CPU
- **Game Service**: 128Mi-256Mi memory, 100m-200m CPU
- **Order Service**: 128Mi-256Mi memory, 100m-200m CPU
- **PostgreSQL**: 256Mi-512Mi memory, 250m-500m CPU

### Health Checks

- All services include liveness and readiness probes
- Services wait for PostgreSQL to be ready via init containers

### Storage

- PostgreSQL uses persistent volumes for data persistence
- Storage class: `local-storage` (adjust for your cluster)

### Security

- Database credentials stored in Kubernetes secrets
- Services communicate internally via ClusterIP services
- External access controlled via NodePort services
