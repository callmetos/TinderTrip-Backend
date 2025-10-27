# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=api
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/api

# Build for Linux
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/api

# Run the application
.PHONY: run
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/api
	./$(BINARY_NAME)

# Run all tests
.PHONY: test
test:
	$(GOTEST) -v ./tests/...

# Run unit tests only
.PHONY: test-unit
test-unit:
	$(GOTEST) -v ./tests/unit/...

# Run integration tests only
.PHONY: test-integration
test-integration:
	$(GOTEST) -v ./tests/integration/...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./tests/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with coverage report in terminal
.PHONY: test-coverage-report
test-coverage-report:
	$(GOTEST) -v -coverprofile=coverage.out ./tests/...
	$(GOCMD) tool cover -func=coverage.out

# Run specific test
.PHONY: test-specific
test-specific:
	@read -p "Enter test pattern: " pattern; \
	$(GOTEST) -v -run $$pattern ./tests/...

# Run benchmarks
.PHONY: benchmark
benchmark:
	$(GOTEST) -bench=. -benchmem ./tests/...

# Clean build files
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run migrations up
.PHONY: migrate-up
migrate-up:
	$(GOBUILD) -o migrate -v ./cmd/migrate
	./migrate up

# Run migrations down
.PHONY: migrate-down
migrate-down:
	$(GOBUILD) -o migrate -v ./cmd/migrate
	./migrate down

# Generate swagger docs
.PHONY: swagger
swagger:
	swag init -g cmd/api/main.go -o docs/

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Run the worker
.PHONY: worker
worker:
	$(GOBUILD) -o worker -v ./cmd/worker
	./worker

# Docker build
.PHONY: docker-build
docker-build:
	docker build -t tinder-trip-backend .

# Docker run
.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 tinder-trip-backend

# Development setup
.PHONY: dev-setup
dev-setup:
	cp env.example .env
	$(GOMOD) download
	$(GOMOD) tidy

# Monitoring setup
.PHONY: monitoring-setup
monitoring-setup:
	./scripts/setup-monitoring.sh

# Start monitoring stack
.PHONY: monitoring-start
monitoring-start:
	docker-compose -f docker-compose.monitoring.yml up -d

# Stop monitoring stack
.PHONY: monitoring-stop
monitoring-stop:
	docker-compose -f docker-compose.monitoring.yml down

# Production deployment
.PHONY: deploy-prod
deploy-prod:
	./scripts/deploy-production.sh

# Deploy monitoring to production
.PHONY: deploy-prod-monitoring
deploy-prod-monitoring:
	docker-compose -f docker-compose.prod.monitoring.yml up -d

# Stop production services
.PHONY: stop-prod
stop-prod:
	docker-compose -f docker-compose.prod.yml down
	docker-compose -f docker-compose.prod.monitoring.yml down

# Production logs
.PHONY: logs-prod
logs-prod:
	docker-compose -f docker-compose.prod.yml logs -f

# Production monitoring logs
.PHONY: logs-prod-monitoring
logs-prod-monitoring:
	docker-compose -f docker-compose.prod.monitoring.yml logs -f

# Test monitoring
.PHONY: monitoring-test
monitoring-test:
	./scripts/test-monitoring.sh

# View monitoring logs
.PHONY: monitoring-logs
monitoring-logs:
	docker-compose -f docker-compose.monitoring.yml logs -f
