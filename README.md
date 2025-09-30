# TinderTrip Backend

[![CI](https://github.com/callmetos/TinderTrip-Backend/actions/workflows/ci.yml/badge.svg)](https://github.com/callmetos/TinderTrip-Backend/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](docker-compose.yml)

A Go-based backend API for a Tinder-like trip matching application built with Gin, PostgreSQL, and Redis.

## 🚀 Features

- **User Authentication**: JWT-based authentication with Google OAuth support
- **Event Management**: Create, join, and manage travel events
- **Swipe System**: Like/pass functionality for event discovery
- **Real-time Chat**: WebSocket-based chat system for event participants
- **User Profiles**: Comprehensive user profiles with preferences
- **Event History**: Track completed and upcoming events
- **Push Notifications**: Real-time notifications for events and messages
- **File Upload**: Nextcloud integration for image and file storage
- **Email Service**: SMTP-based email notifications and password reset
- **Rate Limiting**: API rate limiting with Redis
- **CORS Support**: Cross-origin resource sharing configuration
- **Audit Logging**: Comprehensive audit and API logging
- **Background Workers**: Email and notification processing

## 🏗️ Project Structure

```
├── cmd/                    # Application entry points
│   ├── api/               # Main API server
│   ├── migrate/           # Database migration tool
│   └── worker/            # Background job worker
├── internal/              # Private application code
│   ├── api/              # HTTP handlers and middleware
│   │   ├── handlers/     # Request handlers
│   │   ├── middleware/   # HTTP middleware
│   │   └── routes/       # Route definitions
│   ├── models/           # Data models
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic layer
│   ├── dto/              # Data transfer objects
│   ├── websocket/        # WebSocket handlers
│   ├── validator/        # Input validation
│   └── utils/            # Utility functions
├── pkg/                   # Public packages
│   ├── database/         # Database connections
│   ├── config/           # Configuration management
│   ├── email/            # Email service
│   ├── storage/          # File storage (Nextcloud)
│   └── notification/     # Push notifications
├── scripts/              # Utility scripts
├── tests/                # Test files
└── docs/                 # API documentation
```

## 🛠️ Tech Stack

- **Language**: Go 1.23+
- **Framework**: Gin
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Authentication**: JWT + Google OAuth
- **Email**: SMTP
- **Storage**: Nextcloud
- **Containerization**: Docker & Docker Compose
- **Linting**: golangci-lint
- **Testing**: Go testing framework

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

### Option 1: Using Docker Compose (Recommended)

1. **Clone the repository**
```bash
git clone <repository-url>
cd tinder-trip-backend
```

2. **Set up environment variables**
```bash
cp env.example .env
# Edit .env with your configuration
```

3. **Start the development environment**
```bash
docker-compose up -d
```

4. **Run database migrations**
```bash
go run cmd/migrate/main.go -action=up
```

5. **Access the API**
- API: http://localhost:9952
- Health Check: http://localhost:9952/health
- Swagger Docs: http://localhost:9952/swagger/index.html

### Option 2: Local Development

1. **Install dependencies**
```bash
make deps
```

2. **Start PostgreSQL and Redis**
```bash
# Using Docker
docker run -d --name postgres -e POSTGRES_PASSWORD=password -p 5432:5432 postgres:15-alpine
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

3. **Set up environment variables**
```bash
cp env.example .env
# Edit .env with your configuration
```

4. **Run database migrations**
```bash
make migrate-up
```

5. **Start the application**
```bash
make run
```

## 📚 API Documentation

### Authentication Endpoints

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/google` - Google OAuth URL
- `GET /api/v1/auth/google/callback` - Google OAuth callback
- `POST /api/v1/auth/forgot-password` - Request password reset
- `POST /api/v1/auth/reset-password` - Reset password
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Refresh JWT token

### Event Endpoints

- `GET /api/v1/events` - Get events (with filters)
- `POST /api/v1/events` - Create event
- `GET /api/v1/events/:id` - Get specific event
- `PUT /api/v1/events/:id` - Update event
- `DELETE /api/v1/events/:id` - Delete event
- `POST /api/v1/events/:id/join` - Join event
- `POST /api/v1/events/:id/leave` - Leave event
- `POST /api/v1/events/:id/swipe` - Swipe on event

### User Endpoints

- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile
- `DELETE /api/v1/users/profile` - Delete user profile

### Chat Endpoints

- `GET /api/v1/chat/rooms` - Get chat rooms
- `GET /api/v1/chat/rooms/:id/messages` - Get messages
- `POST /api/v1/chat/rooms/:id/messages` - Send message

### History Endpoints

- `GET /api/v1/history` - Get event history
- `POST /api/v1/history/:id/complete` - Mark event as complete

## 🔧 Development

### Running Tests
```bash
make test
```

### Running with Coverage
```bash
make test-coverage
```

### Linting
```bash
make lint
```

### Building for Production
```bash
make build-linux
```

### Hot Reloading (Development)
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reloading
air
```

## 🐳 Docker

### Development
```bash
# Start all services
docker-compose up -d

# Start only databases
docker-compose up -d postgres redis

# View logs
docker-compose logs -f app
```

### Production
```bash
# Build and start production environment
docker-compose -f docker-compose.prod.yml up -d

# Scale the application
docker-compose -f docker-compose.prod.yml up -d --scale app=3
```

## 📊 Database Schema

The application uses PostgreSQL with the following main tables:

- **users** - User accounts and authentication
- **user_profiles** - Extended user information
- **events** - Travel events and activities
- **event_members** - Event participants
- **event_swipes** - User swipes on events
- **chat_rooms** - Event chat rooms
- **chat_messages** - Chat messages
- **user_event_history** - Event participation history
- **tags** - Event and user tags
- **password_resets** - Password reset tokens
- **audit_logs** - System audit trail
- **api_logs** - API request logs

## 🔐 Environment Variables

See `env.example` for all available configuration options:

- **Database**: PostgreSQL connection settings
- **Redis**: Redis connection settings
- **JWT**: Token configuration
- **Email**: SMTP settings
- **AWS**: S3 storage settings
- **Firebase**: Push notification settings
- **Rate Limiting**: API rate limit configuration
- **CORS**: Cross-origin settings

## 🚀 Deployment

### Using Docker Compose

1. **Set up production environment variables**
```bash
cp env.example .env.production
# Edit .env.production with production values
```

2. **Deploy with Docker Compose**
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Using Kubernetes

1. **Create Kubernetes manifests**
2. **Apply configurations**
```bash
kubectl apply -f k8s/
```

## 🔄 CI Pipeline

This project uses GitHub Actions for continuous integration:

### Automated CI Pipeline

- **Runs on every push and PR**
  - Go 1.23 setup and dependency caching
  - Unit and integration tests with PostgreSQL and Redis
  - Code linting with golangci-lint
  - Docker image building and testing

### Quality Gates

All changes must pass:
- ✅ Unit tests (100% pass rate)
- ✅ Integration tests
- ✅ Code linting (golangci-lint)
- ✅ Docker build verification

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Run linting and tests locally
6. Commit your changes (`git commit -m 'Add some amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Guidelines

- Follow the [Contributing Guide](CONTRIBUTING.md)
- Use conventional commit messages
- Ensure all tests pass
- Update documentation as needed
- Follow the code style guidelines

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:

- Create an issue in the repository
- Check the API documentation
- Review the code examples

## 🔄 Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes and updates.
