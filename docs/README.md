# API Documentation

This directory contains the auto-generated API documentation for the TinderTrip Backend.

## Files

- `swagger.json` - OpenAPI 3.0 specification in JSON format
- `swagger.yaml` - OpenAPI 3.0 specification in YAML format
- `docs.go` - Go code containing the Swagger definitions

## Viewing the Documentation

### Local Development
1. Start the server:
   ```bash
   go run cmd/api/main.go
   ```

2. Open your browser and navigate to:
   ```
   http://localhost:8080/swagger/index.html
   ```

### Docker
1. Start the Docker container:
   ```bash
   docker-compose up -d
   ```

2. Open your browser and navigate to:
   ```
   http://localhost:9952/swagger/index.html
   ```

## Regenerating Documentation

To regenerate the API documentation after making changes to the code:

```bash
# Install swag if not already installed
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g cmd/api/main.go -o docs/
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/google` - Get Google OAuth URL
- `GET /api/v1/auth/google/callback` - Google OAuth callback
- `POST /api/v1/auth/forgot-password` - Request password reset
- `POST /api/v1/auth/reset-password` - Reset password
- `POST /api/v1/auth/logout` - Logout user
- `POST /api/v1/auth/refresh` - Refresh JWT token

### User Profile
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile
- `DELETE /api/v1/users/profile` - Delete user profile

### User Preferences
- `GET /api/v1/users/preferences/availability` - Get availability preferences
- `PUT /api/v1/users/preferences/availability` - Update availability preferences
- `GET /api/v1/users/preferences/budget` - Get budget preferences
- `PUT /api/v1/users/preferences/budget` - Update budget preferences

### Events
- `GET /api/v1/events` - Get events list
- `POST /api/v1/events` - Create new event
- `GET /api/v1/events/:id` - Get specific event
- `PUT /api/v1/events/:id` - Update event
- `DELETE /api/v1/events/:id` - Delete event
- `POST /api/v1/events/:id/join` - Join event
- `POST /api/v1/events/:id/leave` - Leave event
- `POST /api/v1/events/:id/swipe` - Swipe on event

### Chat
- `GET /api/v1/chat/rooms` - Get chat rooms
- `GET /api/v1/chat/rooms/:id/messages` - Get messages
- `POST /api/v1/chat/rooms/:id/messages` - Send message

### History
- `GET /api/v1/history` - Get event history
- `POST /api/v1/history/:id/complete` - Mark event as complete

### Public
- `GET /api/v1/public/events` - Get public events
- `GET /api/v1/public/events/:id` - Get public event details

## Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "message": "Success message",
  "data": { ... },
  "status": "success"
}
```

### Error Response
```json
{
  "error": "Error type",
  "message": "Error description",
  "status": "error"
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse. Default limits:
- 100 requests per hour per IP
- 10 requests per minute per user

## CORS

The API supports Cross-Origin Resource Sharing (CORS) for web applications. Default allowed origins:
- `http://localhost:3000`
- `http://localhost:3001`

## Health Check

Check if the API is running:
```
GET /health
```

Response:
```json
{
  "message": "TinderTrip API is running",
  "status": "ok"
}
```
