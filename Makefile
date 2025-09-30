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

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

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
