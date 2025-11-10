# TinderTrip Monitoring Setup

## 🎯 Overview

This monitoring setup provides comprehensive observability for the TinderTrip backend API using Prometheus, Grafana, and custom application metrics.

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   TinderTrip    │    │   Prometheus    │    │    Grafana      │
│      API        │───▶│   (Metrics)     │───▶│  (Dashboard)    │
│  Port: 8080     │    │  Port: 9090     │    │  Port: 3000     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Health Checks  │    │  System Metrics │    │  Custom Metrics │
│  /health        │    │  Node Exporter  │    │  Business Logic │
│  /ready         │    │  Port: 9100     │    │  User Activity  │
│  /live          │    │                 │    │  Event Tracking │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### 1. Start Monitoring Stack

```bash
# Start all monitoring services
./scripts/setup-monitoring.sh

# Or manually
docker-compose -f docker-compose.monitoring.yml up -d
```

### 2. Start Your API

```bash
# Make sure monitoring is enabled in .env
MONITORING_ENABLED=true

# Start the API
go run cmd/api/main.go
```

### 3. Test Monitoring

```bash
# Test all endpoints
./scripts/test-monitoring.sh
```

## 📊 Available Endpoints

| Service | URL | Purpose |
|---------|-----|---------|
| **API Health** | http://localhost:8080/health | Application health status |
| **API Readiness** | http://localhost:8080/ready | Readiness for traffic |
| **API Liveness** | http://localhost:8080/live | Liveness probe |
| **API Metrics** | http://localhost:9090/metrics | Prometheus metrics |
| **Prometheus UI** | http://localhost:9090 | Metrics query interface |
| **Grafana** | http://localhost:3000 | Dashboards (admin/admin123) |
| **Node Exporter** | http://localhost:9100/metrics | System metrics |

## 📈 Metrics Collected

### HTTP Metrics
- `http_requests_total` - Total HTTP requests by method, endpoint, status
- `http_request_duration_seconds` - Request duration histogram

### Database Metrics
- `db_connections_active` - Active database connections
- `db_connections_idle` - Idle database connections

### Business Metrics
- `user_registrations_total` - Total user registrations
- `event_creations_total` - Total events created
- `active_users` - Current active users

### System Metrics
- `application_uptime_seconds` - Application uptime
- Node exporter metrics (CPU, memory, disk, network)

## 🎛️ Grafana Dashboards

### Main Dashboard Features

1. **Request Rate**: Requests per second by endpoint
2. **Response Time**: 95th and 50th percentile response times
3. **Status Codes**: HTTP status code distribution
4. **Database Health**: Connection pool status
5. **Business Metrics**: User registrations and event creations
6. **System Health**: Uptime and error rates

### Import Dashboard

1. Access Grafana at http://localhost:3000
2. Login with admin/admin123
3. Go to Dashboards → Import
4. Upload `monitoring/grafana/dashboards/tindertrip-dashboard.json`

## 🔧 Configuration

### Environment Variables

```env
# Monitoring Configuration
MONITORING_ENABLED=true
METRICS_PORT=9090
HEALTH_PORT=8080
```

### Prometheus Configuration

Located in `monitoring/prometheus.yml`:
- Scrape interval: 15s
- Evaluation interval: 15s
- Retention: 200h

### Grafana Configuration

- Auto-provisioned datasources
- Pre-configured dashboards
- Default admin user

## 🚨 Health Checks

### Health Endpoint (`/health`)
Returns overall application health:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "uptime": "1h30m45s",
  "version": "1.0.1"
}
```

### Readiness Endpoint (`/ready`)
Checks if application is ready to receive traffic:
```json
{
  "status": "ready",
  "checks": {
    "database": "ok"
  }
}
```

### Liveness Endpoint (`/live`)
Simple liveness check:
```json
{
  "status": "alive",
  "uptime": "1h30m45s"
}
```

## 🐛 Troubleshooting

### Common Issues

1. **Port Conflicts**
   ```bash
   # Check if ports are in use
   lsof -i :9090,3000,9100
   ```

2. **Services Not Starting**
   ```bash
   # Check Docker logs
   docker-compose -f docker-compose.monitoring.yml logs
   ```

3. **Metrics Not Appearing**
   - Ensure `MONITORING_ENABLED=true` in .env
   - Check API logs for monitoring service errors
   - Verify Prometheus can reach the API

### Debug Commands

```bash
# Check service status
docker-compose -f docker-compose.monitoring.yml ps

# View logs
docker-compose -f docker-compose.monitoring.yml logs -f

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:9090/metrics
```

## 📚 Custom Metrics

### Adding New Metrics

1. Define metric in `monitoring_service.go`:
```go
var customMetric = promauto.NewCounter(
    prometheus.CounterOpts{
        Name: "custom_metric_total",
        Help: "Description of custom metric",
    },
)
```

2. Record metric in your code:
```go
customMetric.Inc()
```

### Business Metrics Examples

- User login attempts
- Event swipes (left/right)
- Chat messages sent
- Photo uploads
- API errors by type

## 🔒 Security Considerations

- Change default Grafana password in production
- Use authentication for Prometheus
- Consider HTTPS for production
- Restrict access to monitoring endpoints
- Use proper network segmentation

## 📊 Performance Impact

- **Metrics Collection**: ~1-2ms per request
- **Health Checks**: Negligible impact
- **System Metrics**: Low CPU usage
- **Storage**: ~1GB per day for typical usage

## 🧹 Cleanup

```bash
# Stop and remove all monitoring services
docker-compose -f docker-compose.monitoring.yml down -v

# Remove monitoring data
rm -rf monitoring/grafana/data
```

## 🎯 Next Steps

1. **Alerting**: Set up Prometheus alerting rules
2. **Log Aggregation**: Add ELK stack for logs
3. **Distributed Tracing**: Implement Jaeger or Zipkin
4. **Custom Dashboards**: Create business-specific dashboards
5. **SLA Monitoring**: Set up SLA tracking and alerting
