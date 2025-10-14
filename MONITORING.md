# TinderTrip Monitoring Setup

This document explains how to set up comprehensive monitoring for the TinderTrip backend using Prometheus, Grafana, and custom metrics.

## 🏗️ Architecture

The monitoring stack includes:

- **Prometheus**: Metrics collection and storage
- **Grafana**: Visualization and dashboards
- **Node Exporter**: System metrics
- **Redis Exporter**: Redis metrics
- **PostgreSQL Exporter**: Database metrics
- **Custom TinderTrip Metrics**: Application-specific metrics

## 🚀 Quick Start

### 1. Automated Setup

Run the monitoring setup script:

```bash
./scripts/setup-monitoring.sh
```

This script will:
- Create necessary directories
- Update your `.env` file with monitoring configuration
- Start the monitoring stack with Docker Compose
- Verify all services are running

### 2. Manual Setup

If you prefer to set up manually:

#### Step 1: Update Environment Variables

Add these to your `.env` file:

```env
# Monitoring Configuration
MONITORING_ENABLED=true
METRICS_PORT=9090
HEALTH_PORT=8080
```

#### Step 2: Start Monitoring Stack

```bash
docker-compose -f docker-compose.monitoring.yml up -d
```

#### Step 3: Start Your API

```bash
go run cmd/api/main.go
```

## 📊 Monitoring Endpoints

Once running, you can access:

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin123)
- **API Metrics**: http://localhost:9090/metrics
- **API Health**: http://localhost:8080/health
- **API Readiness**: http://localhost:8080/ready
- **API Liveness**: http://localhost:8080/live

## 📈 Metrics Collected

### HTTP Metrics
- Request rate by method and endpoint
- Response time percentiles
- Status code distribution
- Error rates

### Database Metrics
- Active connections
- Idle connections
- Connection pool status

### Business Metrics
- User registrations
- Event creations
- Active users
- Application uptime

### System Metrics
- CPU usage
- Memory usage
- Disk I/O
- Network I/O

## 🎛️ Grafana Dashboards

### TinderTrip API Dashboard

The main dashboard includes:

1. **HTTP Request Rate**: Requests per second by endpoint
2. **Response Time**: 95th and 50th percentile response times
3. **Status Codes**: Distribution of HTTP status codes
4. **Database Connections**: Active and idle connections
5. **Business Metrics**: User registrations and event creations
6. **Active Users**: Current number of active users
7. **Application Uptime**: Total uptime in seconds
8. **Error Rate**: Percentage of 5xx errors

### Importing Dashboards

1. Access Grafana at http://localhost:3000
2. Login with admin/admin123
3. Go to Dashboards → Import
4. Upload the `monitoring/grafana/dashboards/tindertrip-dashboard.json` file

## 🔧 Configuration

### Prometheus Configuration

The Prometheus configuration is in `monitoring/prometheus.yml`. It includes:

- Scrape intervals
- Target configurations
- Retention policies

### Grafana Configuration

Grafana is configured with:

- Auto-provisioned datasources
- Pre-configured dashboards
- Default admin user

## 🚨 Alerting (Optional)

To set up alerts, you can:

1. Create alert rules in Prometheus
2. Configure notification channels in Grafana
3. Set up alerting rules for:
   - High error rates
   - Slow response times
   - Database connection issues
   - High memory usage

## 📝 Custom Metrics

The application exposes custom metrics:

```go
// HTTP metrics
http_requests_total
http_request_duration_seconds

// Database metrics
db_connections_active
db_connections_idle

// Business metrics
user_registrations_total
event_creations_total
active_users
application_uptime_seconds_total
```

## 🐛 Troubleshooting

### Common Issues

1. **Port conflicts**: Make sure ports 9090, 3000, 9100 are available
2. **Permission issues**: Ensure Docker has proper permissions
3. **Service not starting**: Check Docker logs with `docker-compose logs`

### Checking Service Status

```bash
# Check all services
docker-compose -f docker-compose.monitoring.yml ps

# Check specific service logs
docker-compose -f docker-compose.monitoring.yml logs prometheus
docker-compose -f docker-compose.monitoring.yml logs grafana
```

### Restarting Services

```bash
# Restart all services
docker-compose -f docker-compose.monitoring.yml restart

# Restart specific service
docker-compose -f docker-compose.monitoring.yml restart prometheus
```

## 🧹 Cleanup

To stop and remove all monitoring services:

```bash
docker-compose -f docker-compose.monitoring.yml down -v
```

This will remove containers and volumes.

## 📚 Additional Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)

## 🔒 Security Notes

- Change default Grafana admin password in production
- Use proper authentication for Prometheus
- Consider using HTTPS in production
- Restrict access to monitoring endpoints

## 📊 Performance Impact

The monitoring setup has minimal performance impact:

- Prometheus metrics collection: ~1-2ms per request
- Health checks: Negligible
- System metrics: Low CPU usage
- Storage: ~1GB per day for typical usage

## 🎯 Next Steps

1. Set up log aggregation (ELK stack)
2. Configure alerting rules
3. Set up distributed tracing
4. Implement custom business metrics
5. Set up log-based monitoring
