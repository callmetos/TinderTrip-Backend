# 🚀 TinderTrip Production Deployment Guide

## 📋 Prerequisites

- Docker & Docker Compose installed
- Domain name configured (optional)
- SSL certificates (for HTTPS)
- Environment variables configured

## 🔧 Configuration Changes for Production

### 1. **URL Changes Required:**

| Component | Development | Production |
|-----------|-------------|------------|
| API Base URL | `http://localhost:9952` | `https://your-domain.com/api` |
| Prometheus | `http://localhost:9090` | `https://your-domain.com/prometheus` |
| Grafana | `http://localhost:3001` | `https://your-domain.com/grafana` |
| Health Check | `http://localhost:8082` | `https://your-domain.com/health` |

### 2. **Network Configuration:**
- **Development**: Uses `host.docker.internal` to access host services
- **Production**: Uses Docker service names (`api`, `postgres`, `redis`)

### 3. **Port Configuration:**
- **Development**: Various ports (9090, 3001, 8082, etc.)
- **Production**: Standard ports (80, 443, 9090, 3000)

## 🚀 Deployment Steps

### Step 1: Configure Environment
```bash
# Copy production environment template
cp env.production .env.production

# Edit with your production values
nano .env.production
```

### Step 2: Deploy Application
```bash
# Deploy everything (main app + monitoring)
make deploy-prod

# Or deploy step by step
docker-compose -f docker-compose.prod.yml up -d
make deploy-prod-monitoring
```

### Step 3: Verify Deployment
```bash
# Check all services
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.monitoring.yml ps

# Test health endpoints
curl https://your-domain.com/health
curl https://your-domain.com/prometheus
curl https://your-domain.com/grafana
```

## 🔒 Security Considerations

### 1. **Environment Variables:**
- Use strong passwords for all services
- Store secrets in secure vault (not in files)
- Rotate credentials regularly

### 2. **Network Security:**
- Use reverse proxy (Nginx) for SSL termination
- Restrict access to monitoring ports
- Use firewall rules

### 3. **Grafana Security:**
- Change default admin password
- Enable authentication
- Use HTTPS only

## 📊 Monitoring URLs

| Service | URL | Purpose |
|---------|-----|---------|
| API | `https://your-domain.com/api` | Main API endpoints |
| Health | `https://your-domain.com/health` | Health check |
| Prometheus | `https://your-domain.com/prometheus` | Metrics collection |
| Grafana | `https://your-domain.com/grafana` | Monitoring dashboard |

## 🔧 Management Commands

```bash
# View logs
make logs-prod
make logs-prod-monitoring

# Restart services
docker-compose -f docker-compose.prod.yml restart
docker-compose -f docker-compose.prod.monitoring.yml restart

# Stop all services
make stop-prod

# Update and redeploy
git pull
make deploy-prod
```

## 🌐 Nginx Configuration

Add to your Nginx configuration:

```nginx
# API
location /api/ {
    proxy_pass http://api:8080/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}

# Health check
location /health {
    proxy_pass http://api:8082/health;
}

# Prometheus
location /prometheus/ {
    proxy_pass http://prometheus:9090/;
}

# Grafana
location /grafana/ {
    proxy_pass http://grafana:3000/;
}
```

## 📈 Scaling Considerations

### 1. **Horizontal Scaling:**
- Use load balancer for multiple API instances
- Configure Prometheus to scrape multiple targets
- Use shared storage for Grafana

### 2. **Resource Limits:**
```yaml
# Add to docker-compose.prod.yml
services:
  api:
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
```

### 3. **Monitoring Scaling:**
- Use external Prometheus for large deployments
- Consider Grafana Enterprise for advanced features
- Implement alerting rules

## 🚨 Troubleshooting

### Common Issues:

1. **Port Conflicts:**
   ```bash
   # Check port usage
   netstat -tulpn | grep :80
   netstat -tulpn | grep :443
   ```

2. **Network Issues:**
   ```bash
   # Check Docker networks
   docker network ls
   docker network inspect monitoring
   ```

3. **Service Health:**
   ```bash
   # Check service logs
   docker-compose -f docker-compose.prod.yml logs api
   docker-compose -f docker-compose.prod.monitoring.yml logs prometheus
   ```

## 📝 Environment Variables Reference

| Variable | Description | Example |
|----------|-------------|---------|
| `SERVER_PORT` | API server port | `8080` |
| `METRICS_PORT` | Prometheus metrics port | `9091` |
| `HEALTH_PORT` | Health check port | `8082` |
| `GRAFANA_DOMAIN` | Grafana domain | `your-domain.com` |
| `GRAFANA_ROOT_URL` | Full Grafana URL | `https://your-domain.com/grafana` |
| `MONITORING_ENABLED` | Enable monitoring | `true` |

## 🔄 Updates and Maintenance

### Regular Tasks:
1. Update Docker images
2. Backup Grafana dashboards
3. Monitor disk usage
4. Review logs for errors
5. Update SSL certificates

### Backup Commands:
```bash
# Backup Grafana data
docker run --rm -v grafana_data:/data -v $(pwd):/backup alpine tar czf /backup/grafana-backup.tar.gz -C /data .

# Backup Prometheus data
docker run --rm -v prometheus_data:/data -v $(pwd):/backup alpine tar czf /backup/prometheus-backup.tar.gz -C /data .
```

---

**Note**: This guide assumes you have basic Docker and server administration knowledge. For production deployments, consider consulting with a DevOps engineer.
