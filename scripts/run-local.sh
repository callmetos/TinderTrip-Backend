#!/bin/bash

# TinderTrip Backend - Local Development Script
# This script runs the backend without Docker

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

echo "🚀 Starting TinderTrip Backend (Local Development)"
echo "=================================================="

# Check if .env file exists
if [ ! -f ".env" ]; then
    print_warning ".env file not found. Creating from env.example..."
    cp env.example .env
    print_warning "Please update .env file with your configuration before running again."
    exit 1
fi

# Load environment variables
print_step "Loading environment variables..."
source .env

# Check if PostgreSQL is running
print_step "Checking PostgreSQL connection..."
if ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER > /dev/null 2>&1; then
    print_error "PostgreSQL is not running or not accessible at $DB_HOST:$DB_PORT"
    print_warning "Please start PostgreSQL first:"
    print_warning "  - Using Docker: docker run --name postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15"
    print_warning "  - Using Homebrew: brew services start postgresql"
    exit 1
fi
print_success "PostgreSQL is running"

# Check if Redis is running (optional)
print_step "Checking Redis connection..."
if ! redis-cli -h $REDIS_HOST -p $REDIS_PORT ping > /dev/null 2>&1; then
    print_warning "Redis is not running at $REDIS_HOST:$REDIS_PORT"
    print_warning "The application will continue without Redis (some features may be limited)"
else
    print_success "Redis is running"
fi

# Install dependencies
print_step "Installing Go dependencies..."
go mod tidy
print_success "Dependencies installed"

# Build the application
print_step "Building the application..."
go build -o bin/api cmd/api/main.go
print_success "Application built successfully"

# Run database migrations
print_step "Running database migrations..."
if [ -f "bin/migrate" ]; then
    ./bin/migrate up
    print_success "Database migrations completed"
else
    print_warning "Migration binary not found. Skipping migrations."
fi

# Start the server
print_step "Starting the server..."
print_success "Server will be available at:"
print_success "  - Local: http://localhost:$SERVER_PORT"
print_success "  - LAN: http://$(hostname -I | awk '{print $1}'):$SERVER_PORT"
print_success "  - Swagger: http://localhost:$SERVER_PORT/swagger/index.html"
print_success ""
print_success "Press Ctrl+C to stop the server"

# Run the application
./bin/api
