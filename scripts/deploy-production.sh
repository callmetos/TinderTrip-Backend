#!/bin/bash

# TinderTrip Production Deployment Script
set -e

echo "🚀 Starting TinderTrip Production Deployment..."

# Check if .env.production exists
if [ ! -f ".env.production" ]; then
    echo "❌ Error: .env.production file not found!"
    echo "Please copy env.production to .env.production and configure your settings."
    exit 1
fi

# Load environment variables
export $(cat .env.production | grep -v '^#' | xargs)

echo "📋 Configuration:"
echo "  - API Port: ${SERVER_PORT:-8080}"
echo "  - Metrics Port: ${METRICS_PORT:-9091}"
echo "  - Health Port: ${HEALTH_PORT:-8082}"
echo "  - Grafana Domain: ${GRAFANA_DOMAIN:-localhost}"

# Create monitoring network
echo "🌐 Creating monitoring network..."
docker network create monitoring 2>/dev/null || echo "Network already exists"

# Build and start main application
echo "🏗️  Building and starting main application..."
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Start monitoring stack
echo "📊 Starting monitoring stack..."
docker-compose -f docker-compose.prod.monitoring.yml down
docker-compose -f docker-compose.prod.monitoring.yml up -d

# Wait for monitoring services
echo "⏳ Waiting for monitoring services..."
sleep 20

# Health check
echo "🔍 Performing health checks..."

# Check API health
if curl -f -s "http://localhost:${HEALTH_PORT:-8082}/health" > /dev/null; then
    echo "✅ API Health Check: OK"
else
    echo "❌ API Health Check: FAILED"
fi

# Check Prometheus
if curl -f -s "http://localhost:9090/-/healthy" > /dev/null; then
    echo "✅ Prometheus: OK"
else
    echo "❌ Prometheus: FAILED"
fi

# Check Grafana
if curl -f -s "http://localhost:3000/api/health" > /dev/null; then
    echo "✅ Grafana: OK"
else
    echo "❌ Grafana: FAILED"
fi

echo ""
echo "🎉 Deployment completed!"
echo ""
echo "📊 Access URLs:"
echo "  - API: http://localhost:${SERVER_PORT:-8080}"
echo "  - API Health: http://localhost:${HEALTH_PORT:-8082}/health"
echo "  - Prometheus: http://localhost:9090"
echo "  - Grafana: http://localhost:3000"
echo "    - Username: ${GRAFANA_ADMIN_USER:-admin}"
echo "    - Password: ${GRAFANA_ADMIN_PASSWORD:-admin123}"
echo ""
echo "🔧 Management Commands:"
echo "  - View logs: docker-compose -f docker-compose.prod.yml logs -f"
echo "  - Stop services: docker-compose -f docker-compose.prod.yml down"
echo "  - Restart monitoring: docker-compose -f docker-compose.prod.monitoring.yml restart"
