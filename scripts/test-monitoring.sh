#!/bin/bash

# TinderTrip Monitoring Test Script
# This script tests the monitoring endpoints

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Test function
test_endpoint() {
    local url=$1
    local name=$2
    local expected_status=${3:-200}
    
    print_status "Testing $name at $url..."
    
    if response=$(curl -s -w "%{http_code}" -o /dev/null "$url" 2>/dev/null); then
        if [ "$response" = "$expected_status" ]; then
            print_status "✅ $name is responding correctly (HTTP $response)"
            return 0
        else
            print_error "❌ $name returned HTTP $response, expected $expected_status"
            return 1
        fi
    else
        print_error "❌ $name is not responding"
        return 1
    fi
}

# Test monitoring endpoints
print_status "🧪 Testing TinderTrip Monitoring Endpoints..."
echo ""

# Test API health endpoints (using SERVER_PORT from .env)
SERVER_PORT=$(grep SERVER_PORT .env | cut -d '=' -f2)
HEALTH_PORT=$(grep HEALTH_PORT .env | cut -d '=' -f2)
test_endpoint "http://localhost:${HEALTH_PORT}/health" "Health Check"
test_endpoint "http://localhost:${HEALTH_PORT}/ready" "Readiness Check"
test_endpoint "http://localhost:${HEALTH_PORT}/live" "Liveness Check"

echo ""

# Test metrics endpoint
test_endpoint "http://localhost:9090/metrics" "Prometheus Metrics"

echo ""

# Test monitoring stack
test_endpoint "http://localhost:9090" "Prometheus UI" "302"
test_endpoint "http://localhost:3001" "Grafana UI" "302"

echo ""

# Test system metrics
test_endpoint "http://localhost:9100/metrics" "Node Exporter"

echo ""

print_status "🎉 Monitoring test completed!"
print_status "You can now access:"
print_status "  - Prometheus: http://localhost:9090"
print_status "  - Grafana: http://localhost:3001 (admin/admin123)"
print_status "  - API Health: http://localhost:${HEALTH_PORT}/health"
print_status "  - API Metrics: http://localhost:9090/metrics"
