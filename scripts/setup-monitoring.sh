#!/bin/bash

# TinderTrip Monitoring Setup Script
# This script sets up monitoring infrastructure for TinderTrip backend

set -e

echo "🚀 Setting up TinderTrip Monitoring Infrastructure..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Create monitoring directories if they don't exist
print_status "Creating monitoring directories..."
mkdir -p monitoring/grafana/provisioning/datasources
mkdir -p monitoring/grafana/provisioning/dashboards
mkdir -p monitoring/grafana/dashboards

# Check if .env file exists
if [ ! -f .env ]; then
    print_warning ".env file not found. Creating from env.example..."
    cp env.example .env
    print_warning "Please update .env file with your actual configuration values."
fi

# Update .env file with monitoring settings if not already present
if ! grep -q "MONITORING_ENABLED" .env; then
    print_status "Adding monitoring configuration to .env file..."
    echo "" >> .env
    echo "# Monitoring Configuration" >> .env
    echo "MONITORING_ENABLED=true" >> .env
    echo "METRICS_PORT=9090" >> .env
    echo "HEALTH_PORT=8080" >> .env
fi

# Start monitoring stack
print_status "Starting monitoring stack with Docker Compose..."
docker-compose -f docker-compose.monitoring.yml up -d

# Wait for services to be ready
print_status "Waiting for services to be ready..."
sleep 10

# Check if services are running
print_status "Checking service status..."

# Check Prometheus
if curl -s http://localhost:9090 > /dev/null; then
    print_status "✅ Prometheus is running at http://localhost:9090"
else
    print_error "❌ Prometheus is not responding"
fi

# Check Grafana
if curl -s http://localhost:3000 > /dev/null; then
    print_status "✅ Grafana is running at http://localhost:3000"
    print_status "   Default credentials: admin / admin123"
else
    print_error "❌ Grafana is not responding"
fi

# Check Node Exporter
if curl -s http://localhost:9100 > /dev/null; then
    print_status "✅ Node Exporter is running at http://localhost:9100"
else
    print_warning "⚠️  Node Exporter is not responding"
fi

# Display monitoring URLs
echo ""
print_status "🎉 Monitoring setup complete!"
echo ""
echo "📊 Monitoring URLs:"
echo "   Prometheus: http://localhost:9090"
echo "   Grafana:    http://localhost:3000 (admin/admin123)"
echo "   Node Exporter: http://localhost:9100"
echo ""
echo "🔧 Next steps:"
echo "   1. Start your TinderTrip API with monitoring enabled"
echo "   2. Access Grafana and import the TinderTrip dashboard"
echo "   3. Configure alerts as needed"
echo ""
echo "📝 To stop monitoring:"
echo "   docker-compose -f docker-compose.monitoring.yml down"
echo ""
echo "📝 To view logs:"
echo "   docker-compose -f docker-compose.monitoring.yml logs -f"
