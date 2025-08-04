# LUGX Gaming Platform Monitoring

This directory contains comprehensive monitoring setup using Prometheus and Grafana for the LUGX Gaming Platform, providing service health monitoring, performance metrics, and availability tracking.

## 📊 Monitoring Stack Overview

### Components

- **Prometheus**: Metrics collection and alerting
- **Grafana**: Visualization and dashboards
- **Node Exporter**: System-level metrics collection
- **Alert Manager**: Alert routing and management (via Prometheus rules)

### Metrics Collected

- **Service Health**: Endpoint availability and response status
- **Performance**: CPU usage, memory consumption, response times
- **Availability**: Uptime tracking and service status
- **System Metrics**: Node-level CPU, memory, disk, and network metrics
- **Kubernetes Metrics**: Pod, deployment, and cluster-level metrics

## 🚀 Quick Start

### Deploy Monitoring Stack

```powershell
# Deploy the entire LUGX platform including monitoring
.\deploy.ps1

# Or deploy only monitoring components
.\deploy-monitoring.ps1
```

### Access Monitoring Services

- **Prometheus**: http://localhost:30090
- **Grafana**: http://localhost:30300
  - Username: `admin`
  - Password: `admin`

### Clean Up

```powershell
# Remove only monitoring components
.\cleanup-monitoring.ps1

# Remove entire platform
.\cleanup.ps1
```

## 📈 Pre-configured Dashboards

### LUGX Gaming Platform Monitoring Dashboard

The deployment includes a pre-configured Grafana dashboard that provides:

1. **Service Availability Panel**: Real-time status of all LUGX services
2. **Service Status Overview**: Visual indicators for each service (green = up, red = down)
3. **CPU Usage Monitoring**: Resource utilization across all pods
4. **Memory Usage Tracking**: Memory consumption patterns
5. **Performance Trends**: Historical data for capacity planning

## 🔧 Monitoring Configuration

### Prometheus Configuration

Located in `prometheus-config.yaml`, includes:

- **Scrape Configs**: Auto-discovery of LUGX Gaming services
- **Alert Rules**: Pre-defined alerts for service health and performance
- **Retention**: 200 hours of metric data storage

### Key Scrape Jobs

```yaml
- lugx-frontend: Frontend application metrics
- lugx-game-service: Game management API health
- lugx-order-service: Order processing API health
- lugx-analytics-service: Analytics API health
- clickhouse: Database metrics
- kubernetes-pods: Auto-discovered pod metrics
- kubernetes-nodes: Node-level system metrics
```

### Alert Rules

The monitoring stack includes pre-configured alerts:

- **ServiceDown**: Triggers when a service is unavailable for 2+ minutes
- **HighErrorRate**: Alerts when error rate exceeds 10% for 5+ minutes
- **HighMemoryUsage**: Warns when memory usage exceeds 80% for 5+ minutes
- **HighCPUUsage**: Warns when CPU usage exceeds 80% for 5+ minutes

## 🎯 Service Discovery

Services are automatically discovered using Kubernetes annotations:

```yaml
annotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "8080"
  prometheus.io/path: "/health"
```

These annotations are already added to all LUGX Gaming services.

## 📊 Available Metrics

### Service Health Metrics

- `up`: Service availability (1 = up, 0 = down)
- `probe_success`: Health check success rate
- `probe_duration_seconds`: Health check response time

### Resource Metrics (via Node Exporter)

- `node_cpu_seconds_total`: CPU usage statistics
- `node_memory_MemTotal_bytes`: Total system memory
- `node_memory_MemAvailable_bytes`: Available memory
- `node_filesystem_size_bytes`: Disk usage information

### Kubernetes Metrics

- `container_cpu_usage_seconds_total`: Container CPU usage
- `container_memory_usage_bytes`: Container memory usage
- `kube_pod_status_phase`: Pod status information

## 🔍 Custom Queries

### Useful Prometheus Queries

**Service Availability**:

```promql
up{job=~"lugx-.*"}
```

**CPU Usage by Service**:

```promql
rate(container_cpu_usage_seconds_total{namespace="lugx-gaming"}[5m]) * 100
```

**Memory Usage by Service**:

```promql
container_memory_usage_bytes{namespace="lugx-gaming"}
```

**Service Response Time**:

```promql
probe_duration_seconds{job=~"lugx-.*"}
```

**Error Rate**:

```promql
rate(http_requests_total{status=~"5.."}[5m])
```

## 📋 Troubleshooting

### Common Issues

**Prometheus Can't Discover Services**:

- Verify services have proper annotations
- Check RBAC permissions for Prometheus
- Ensure services are in the correct namespace

**Grafana Dashboard Not Loading**:

- Verify Prometheus datasource is configured
- Check if Prometheus is reachable from Grafana
- Restart Grafana pod if needed

**Metrics Not Appearing**:

- Check if target services are running
- Verify health endpoints are responding
- Review Prometheus targets page for errors

### Debugging Commands

```powershell
# Check monitoring pod status
kubectl get pods -n monitoring

# View Prometheus logs
kubectl logs -n monitoring deployment/prometheus

# View Grafana logs
kubectl logs -n monitoring deployment/grafana

# Check service discovery
kubectl get pods -n lugx-gaming -o wide

# Test health endpoints
kubectl port-forward -n lugx-gaming deployment/game-service 8080:8080
curl http://localhost:8080/api/v1/health
```

## 🔧 Advanced Configuration

### Adding Custom Dashboards

1. Create dashboard JSON file
2. Add to `grafana-config.yaml` ConfigMap
3. Redeploy Grafana configuration

### Modifying Alert Rules

1. Edit `prometheus-config.yaml`
2. Update `alerting_rules.yml` section
3. Redeploy Prometheus configuration

### Scaling Monitoring

```powershell
# Scale Prometheus for high availability
kubectl scale deployment prometheus --replicas=2 -n monitoring

# Increase storage for metrics retention
kubectl patch pvc prometheus-pvc -n monitoring -p '{"spec":{"resources":{"requests":{"storage":"50Gi"}}}}'
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Monitoring Stack                        │
├─────────────────────────────────────────────────────────────┤
│  Grafana (Port 30300)                                      │
│    ├── Datasource: Prometheus                              │
│    ├── Dashboard: LUGX Gaming Platform                     │
│    └── Alerts: Visual notifications                        │
├─────────────────────────────────────────────────────────────┤
│  Prometheus (Port 30090)                                   │
│    ├── Scrape: LUGX services health endpoints              │
│    ├── Scrape: Node Exporter metrics                       │
│    ├── Scrape: Kubernetes API metrics                      │
│    ├── Storage: 20Gi persistent volume                     │
│    └── Alerts: Service & performance rules                 │
├─────────────────────────────────────────────────────────────┤
│  Node Exporter (DaemonSet)                                 │
│    ├── Metrics: CPU, Memory, Disk, Network                 │
│    └── Endpoint: :9100/metrics                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ Scrapes
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     LUGX Gaming Services                    │
├─────────────────────────────────────────────────────────────┤
│  Frontend (:80)              │  Game Service (:8080)       │
│  Order Service (:8080)       │  Analytics Service (:8080)  │
│  ClickHouse (:8123/:9000)    │  PostgreSQL (External)     │
└─────────────────────────────────────────────────────────────┘
```

## 📈 Benefits

### For Developers

- **Real-time Health Monitoring**: Immediate visibility into service status
- **Performance Insights**: CPU and memory usage patterns
- **Error Tracking**: Quick identification of failing components
- **Historical Data**: Trend analysis for capacity planning

### For Operations

- **Proactive Monitoring**: Alerts before issues impact users
- **Resource Optimization**: Data-driven scaling decisions
- **Incident Response**: Faster troubleshooting with detailed metrics
- **Compliance**: Audit trails and availability reporting

### For Business

- **Service Reliability**: Improved uptime and user experience
- **Cost Optimization**: Efficient resource utilization
- **Scalability Planning**: Data-driven growth decisions
- **Quality Assurance**: Continuous performance validation

## 🔗 Related Files

- `prometheus-config.yaml`: Prometheus server configuration
- `prometheus.yaml`: Prometheus deployment and services
- `grafana-config.yaml`: Grafana configuration and dashboards
- `grafana.yaml`: Grafana deployment and services
- `node-exporter.yaml`: System metrics collection
- `deploy-monitoring.ps1`: Automated deployment script
- `cleanup-monitoring.ps1`: Cleanup script

For more information about the LUGX Gaming Platform, see the main [README.md](README.md).
