# Kubernetes Deployment Guide for CasWAF

This guide provides instructions for deploying CasWAF on Kubernetes.

## Prerequisites

- A running Kubernetes cluster (1.19+)
- `kubectl` configured to access your cluster
- A running Casdoor instance (can be in the same cluster or external)
- Basic understanding of Kubernetes resources

## Architecture

The deployment consists of:
- **CasWAF Application**: The main WAF application
- **MySQL Database**: Stores CasWAF configuration and data
- **ConfigMap**: Contains CasWAF configuration
- **Services**: Exposes CasWAF and MySQL within the cluster
- **Ingress** (optional): Exposes CasWAF externally

## Quick Start

### 1. Deploy Casdoor (if not already deployed)

CasWAF requires Casdoor for authentication. If you don't have Casdoor deployed:

```bash
# Follow Casdoor's Kubernetes deployment guide:
# https://casdoor.org/docs/deployment/k8s
```

### 2. Configure CasWAF

Edit `k8s/configmap.yaml` and update the following values:

```yaml
# MySQL connection
dataSourceName: root:YOUR_MYSQL_PASSWORD@tcp(caswaf-mysql:3306)/

# Casdoor configuration
casdoorEndpoint: http://casdoor.casdoor-system.svc.cluster.local:8000
clientId: "YOUR_CLIENT_ID"
clientSecret: "YOUR_CLIENT_SECRET"
casdoorOrganization: "YOUR_ORGANIZATION"
casdoorApplication: "YOUR_APPLICATION"
```

**Important**: 
- Replace `YOUR_MYSQL_PASSWORD` with your MySQL root password
- Get `clientId` and `clientSecret` from your Casdoor application settings
- Ensure `casdoorEndpoint` points to your Casdoor instance

### 3. Configure MySQL Password

Edit `k8s/mysql.yaml` and update the MySQL root password:

```bash
# Generate base64 encoded password
echo -n "your-secure-password" | base64
```

Then update the `mysql-root-password` in the Secret resource with the base64 value.

### 4. Deploy to Kubernetes

```bash
# Create namespace and deploy MySQL
kubectl apply -f k8s/mysql.yaml

# Wait for MySQL to be ready
kubectl wait --for=condition=ready pod -l app=caswaf-mysql -n caswaf --timeout=300s

# Deploy ConfigMap
kubectl apply -f k8s/configmap.yaml

# Deploy CasWAF
kubectl apply -f k8s/deployment.yaml

# (Optional) Deploy Ingress
kubectl apply -f k8s/ingress.yaml
```

### 5. Verify Deployment

```bash
# Check if pods are running
kubectl get pods -n caswaf

# Check logs
kubectl logs -f deployment/caswaf -n caswaf

# Check services
kubectl get svc -n caswaf
```

### 6. Access CasWAF

If using Ingress:
```bash
# Update your DNS or /etc/hosts to point to your ingress controller IP
# Then access: http://caswaf.example.com
```

If using port-forward for testing:
```bash
kubectl port-forward svc/caswaf 7000:7000 -n caswaf
# Access: http://localhost:7000
```

## Configuration Details

### ConfigMap (`configmap.yaml`)

Key configuration parameters:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `httpport` | CasWAF HTTP port | `7000` |
| `runmode` | Run mode (dev/prod) | `prod` |
| `driverName` | Database driver | `mysql` |
| `dataSourceName` | MySQL connection string | `root:password@tcp(caswaf-mysql:3306)/` |
| `dbName` | Database name | `caswaf` |
| `casdoorEndpoint` | Casdoor API endpoint | Required |
| `casdoorInsecureSkipVerify` | Skip TLS verification for Casdoor | `true` |
| `clientId` | Casdoor application client ID | Required |
| `clientSecret` | Casdoor application client secret | Required |
| `casdoorOrganization` | Casdoor organization name | `built-in` |
| `casdoorApplication` | Casdoor application name | Required |

### MySQL Deployment (`mysql.yaml`)

- Uses MySQL 8.0.25
- Persistent storage with PVC (10Gi)
- Includes health checks
- Root password stored in Kubernetes Secret

### CasWAF Deployment (`deployment.yaml`)

Features:
- Init container to wait for MySQL readiness
- Liveness and readiness probes
- Resource limits and requests
- Configuration mounted from ConfigMap

## Troubleshooting

### Common Issues

#### 1. "wait-for-it: timeout occurred after waiting 15 seconds for db:3306"

**Cause**: MySQL is not ready or not accessible

**Solution**:
```bash
# Check MySQL pod status
kubectl get pods -n caswaf -l app=caswaf-mysql

# Check MySQL logs
kubectl logs -n caswaf -l app=caswaf-mysql

# Verify MySQL service
kubectl get svc -n caswaf caswaf-mysql
```

#### 2. "casdoorsdk.GetCerts() error: Unauthorized operation"

**Cause**: Incorrect Casdoor configuration or credentials

**Solution**:
1. Verify Casdoor is accessible:
   ```bash
   kubectl run -it --rm debug --image=curlimages/curl --restart=Never -n caswaf -- \
     curl -v http://casdoor.casdoor-system.svc.cluster.local:8000
   ```

2. Verify `clientId` and `clientSecret` in ConfigMap match your Casdoor application

3. Ensure the Casdoor application is configured correctly:
   - Organization name matches `casdoorOrganization`
   - Application name matches `casdoorApplication`
   - Client ID and secret are correct

4. Check Casdoor logs for authentication errors

#### 3. Database Connection Issues

**Solution**:
```bash
# Test MySQL connection from CasWAF pod
kubectl exec -it deployment/caswaf -n caswaf -- sh
# Then inside the pod:
# nc -zv caswaf-mysql 3306
```

#### 4. Init Container Stuck

If the init container is stuck waiting for MySQL:
```bash
# Check init container logs
kubectl logs -n caswaf <pod-name> -c wait-for-mysql

# Force restart
kubectl rollout restart deployment/caswaf -n caswaf
```

### Viewing Logs

```bash
# CasWAF logs
kubectl logs -f deployment/caswaf -n caswaf

# MySQL logs
kubectl logs -f deployment/caswaf-mysql -n caswaf

# All logs in namespace
kubectl logs -f -n caswaf --all-containers=true
```

## Production Recommendations

1. **Use External MySQL**: For production, consider using a managed MySQL service (AWS RDS, Google Cloud SQL, etc.)

2. **Configure TLS**: 
   - Set `casdoorInsecureSkipVerify = false`
   - Use proper TLS certificates for Casdoor

3. **Resource Limits**: Adjust resource limits based on your traffic:
   ```yaml
   resources:
     requests:
       memory: "512Mi"
       cpu: "500m"
     limits:
       memory: "2Gi"
       cpu: "2000m"
   ```

4. **High Availability**:
   - Increase replicas for CasWAF
   - Use MySQL replication or managed service
   - Configure proper health checks

5. **Monitoring**: Set up monitoring and alerting:
   - Prometheus metrics
   - Application logs
   - Resource usage

6. **Backup**: Regular backups of MySQL data

7. **Security**:
   - Use Kubernetes Secrets for sensitive data
   - Enable RBAC
   - Network policies to restrict traffic
   - Regular security updates

## Updating CasWAF

```bash
# Update the image version in deployment.yaml, then:
kubectl set image deployment/caswaf caswaf=casbin/caswaf:NEW_VERSION -n caswaf

# Or apply updated deployment
kubectl apply -f k8s/deployment.yaml

# Check rollout status
kubectl rollout status deployment/caswaf -n caswaf
```

## Uninstall

```bash
# Delete all resources
kubectl delete -f k8s/ingress.yaml
kubectl delete -f k8s/deployment.yaml
kubectl delete -f k8s/configmap.yaml
kubectl delete -f k8s/mysql.yaml

# Or delete the entire namespace
kubectl delete namespace caswaf
```

## Advanced Configuration

### Using External MySQL

Edit `configmap.yaml`:
```yaml
dataSourceName: root:password@tcp(external-mysql.example.com:3306)/
```

Then skip deploying `mysql.yaml`.

### Using Redis for Sessions

Edit `configmap.yaml`:
```yaml
redisEndpoint: redis-host:6379
```

### Custom Domain

Edit `ingress.yaml`:
```yaml
spec:
  rules:
  - host: waf.yourdomain.com
```

## Support

For issues and questions:
- GitHub Issues: https://github.com/casbin/caswaf/issues
- Documentation: https://caswaf.org
- Casdoor Documentation: https://casdoor.org

## License

Apache-2.0
