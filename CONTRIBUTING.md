# Contributing to TinderTrip Backend

Thank you for your interest in contributing to TinderTrip Backend! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Documentation](#documentation)
- [Issue Reporting](#issue-reporting)

## Code of Conduct

This project follows a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/TinderTrip-Backend.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.23 or later
- PostgreSQL 15 or later
- Redis 7 or later (optional)
- Docker and Docker Compose (for containerized development)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/TinderTrip-Backend.git
   cd TinderTrip-Backend
   ```

2. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

3. **Install dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

4. **Set up the database**
   ```bash
   # Using Docker Compose
   docker-compose up -d postgres
   
   # Or manually create a PostgreSQL database
   createdb tinder_trip_backend
   ```

5. **Run migrations**
   ```bash
   make migrate-up
   ```

6. **Start the development server**
   ```bash
   make run
   # Or
   go run cmd/api/main.go
   ```

### Docker Development

1. **Start all services**
   ```bash
   docker-compose up -d
   ```

2. **View logs**
   ```bash
   docker-compose logs -f api
   ```

3. **Stop services**
   ```bash
   docker-compose down
   ```

## Making Changes

### Branch Naming

Use descriptive branch names:
- `feature/add-user-authentication`
- `bugfix/fix-login-error`
- `hotfix/security-patch`
- `refactor/cleanup-database-layer`

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
type(scope): description

[optional body]

[optional footer(s)]
```

Examples:
- `feat(auth): add Google OAuth integration`
- `fix(api): resolve user profile update issue`
- `docs(readme): update installation instructions`
- `refactor(db): optimize user queries`

### Code Structure

- **Handlers**: API request/response handling
- **Services**: Business logic
- **Repositories**: Data access layer
- **Models**: Database entities
- **DTOs**: Data transfer objects
- **Middleware**: Cross-cutting concerns

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test package
go test ./internal/service/...

# Run tests with race detection
go test -race ./...
```

### Writing Tests

- Write unit tests for all business logic
- Write integration tests for API endpoints
- Aim for at least 80% test coverage
- Use table-driven tests where appropriate

### Test Structure

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserRequest
        want    *User
        wantErr bool
    }{
        {
            name: "valid user",
            input: CreateUserRequest{
                Email: "test@example.com",
                Password: "password123",
            },
            want: &User{
                Email: "test@example.com",
            },
            wantErr: false,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Submitting Changes

### Pull Request Process

1. **Create a pull request** from your feature branch to `main`
2. **Fill out the PR template** completely
3. **Ensure all checks pass** (tests, linting, security scans)
4. **Request review** from maintainers
5. **Address feedback** promptly
6. **Squash commits** if requested

### Pull Request Guidelines

- Keep PRs focused and small
- Include tests for new functionality
- Update documentation as needed
- Ensure CI/CD pipeline passes
- Follow the PR template

## Code Style

### Go Style Guide

Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines:

- Use `gofmt` and `goimports`
- Follow naming conventions
- Write clear, readable code
- Add comments for exported functions
- Use meaningful variable names

### Linting

We use `golangci-lint` for code quality:

```bash
# Run linter
make lint

# Fix auto-fixable issues
golangci-lint run --fix
```

### Code Review Checklist

- [ ] Code follows Go style guidelines
- [ ] Functions are well-documented
- [ ] Error handling is appropriate
- [ ] Tests are comprehensive
- [ ] No security vulnerabilities
- [ ] Performance is acceptable
- [ ] Database queries are optimized

## Documentation

### API Documentation

- Use Swagger/OpenAPI annotations
- Generate docs with `make swagger`
- Keep API docs up to date

### Code Documentation

- Document all exported functions
- Use clear, concise comments
- Include examples where helpful
- Update README for significant changes

## Issue Reporting

### Bug Reports

Use the bug report template and include:
- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details
- Error messages/logs

### Feature Requests

Use the feature request template and include:
- Clear description of the feature
- Problem it solves
- Proposed solution
- Use cases
- Acceptance criteria

## Getting Help

- Check existing issues and PRs
- Read the documentation
- Ask questions in discussions
- Contact maintainers if needed

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to TinderTrip Backend! 🚀
