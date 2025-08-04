# External Database Configuration Guide

This guide helps you configure the Kubernetes services to connect to external PostgreSQL and ClickHouse databases.

## Quick Setup Steps

### 1. Update Database Hosts

Edit `external-services.yaml` and update the following values in the ConfigMap:

```yaml
data:
  # Replace with your actual database hosts
  POSTGRES_HOST: "your-postgres-host.example.com" # or IP address
  CLICKHOUSE_HOST: "your-clickhouse-host.example.com" # or IP address
```

### 2. Update Database Credentials

Generate base64 encoded credentials and update the Secret in `external-services.yaml`:

**PowerShell commands to encode credentials:**

```powershell
# PostgreSQL credentials
[System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes("your_postgres_username"))
[System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes("your_postgres_password"))

# ClickHouse credentials
[System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes("your_clickhouse_username"))
[System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes("your_clickhouse_password"))
```

**Linux/macOS commands:**

```bash
# PostgreSQL credentials
echo -n "your_postgres_username" | base64
echo -n "your_postgres_password" | base64

# ClickHouse credentials
echo -n "your_clickhouse_username" | base64
echo -n "your_clickhouse_password" | base64
```

### 3. Common Host Values

- **Local development**: `host.docker.internal` (if running databases locally)
- **Docker Compose**: Use service names like `postgres` or `clickhouse`
- **External cloud databases**: Use the provided connection hostnames
- **On-premise**: Use IP addresses or hostnames accessible from Kubernetes cluster

### 4. Network Considerations

Ensure your Kubernetes cluster can reach the external databases:

1. **Firewall rules**: Allow connections from Kubernetes nodes to database ports
2. **Network connectivity**: Ensure network routes exist between Kubernetes and databases
3. **DNS resolution**: Verify hostnames resolve correctly from within Kubernetes pods

### 5. Testing Connectivity

After deployment, you can test connectivity using the init containers:

```bash
# Check init container logs to see if database connections succeed
kubectl logs -n lugx-gaming -l app=game-service -c wait-for-postgres
kubectl logs -n lugx-gaming -l app=order-service -c wait-for-postgres
kubectl logs -n lugx-gaming -l app=analytics-service -c wait-for-clickhouse
```

### 6. Troubleshooting

**Connection timeouts:**

- Verify database hosts are correct
- Check network connectivity from Kubernetes nodes
- Ensure firewall rules allow the traffic

**Authentication failures:**

- Verify base64 encoded credentials are correct
- Check if the database users have proper permissions
- Ensure SSL/TLS settings match between client and server

**DNS resolution issues:**

- Test DNS resolution from within a Kubernetes pod
- Consider using IP addresses instead of hostnames
- Check if custom DNS servers are needed
