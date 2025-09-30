#!/bin/bash

# Development script for TinderTrip Backend

set -e

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

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    print_error "docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

# Function to start development environment
start_dev() {
    print_status "Starting development environment..."
    
    # Start services
    docker-compose up -d postgres redis
    
    # Wait for services to be ready
    print_status "Waiting for services to be ready..."
    sleep 10
    
    # Run migrations
    print_status "Running database migrations..."
    go run cmd/migrate/main.go -action=up
    
    # Start the application
    print_status "Starting application..."
    go run cmd/api/main.go
}

# Function to stop development environment
stop_dev() {
    print_status "Stopping development environment..."
    docker-compose down
}

# Function to run tests
run_tests() {
    print_status "Running tests..."
    go test ./...
}

# Function to run linter
run_linter() {
    print_status "Running linter..."
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run
    else
        print_warning "golangci-lint is not installed. Installing..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        golangci-lint run
    fi
}

# Function to run migrations
run_migrations() {
    print_status "Running migrations..."
    go run cmd/migrate/main.go -action=up
}

# Function to rollback migrations
rollback_migrations() {
    print_status "Rolling back migrations..."
    go run cmd/migrate/main.go -action=down
}

# Function to show help
show_help() {
    echo "TinderTrip Backend Development Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  start       Start development environment"
    echo "  stop        Stop development environment"
    echo "  test        Run tests"
    echo "  lint        Run linter"
    echo "  migrate     Run database migrations"
    echo "  rollback    Rollback database migrations"
    echo "  help        Show this help message"
    echo ""
}

# Main script logic
case "${1:-help}" in
    start)
        start_dev
        ;;
    stop)
        stop_dev
        ;;
    test)
        run_tests
        ;;
    lint)
        run_linter
        ;;
    migrate)
        run_migrations
        ;;
    rollback)
        rollback_migrations
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac
