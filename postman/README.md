# Postman Collection for TinderTrip API

This directory contains Postman collections and environment files for testing the TinderTrip API.

## Files

- `TinderTrip-API.postman_collection.json` - Complete API collection
- `TinderTrip-API.postman_environment.json` - Environment variables
- `README.md` - This documentation

## Setup

### 1. Import Collection

1. Open Postman
2. Click "Import" button
3. Select `TinderTrip-API.postman_collection.json`
4. Click "Import"

### 2. Import Environment

1. In Postman, go to "Environments"
2. Click "Import"
3. Select `TinderTrip-API.postman_environment.json`
4. Click "Import"
5. Select "TinderTrip API Environment" from the environment dropdown

### 3. Start the API Server

```bash
# Using Docker Compose
docker-compose up -d

# Or using Go directly
go run cmd/api/main.go
```

## Usage

### 1. Health Check

Start by testing the health endpoint:
- **GET** `/health`
- Should return: `{"message":"TinderTrip API is running","status":"ok"}`

### 2. Authentication Flow

#### Register a User
1. Go to **Authentication** → **Register User**
2. Update the request body with your details
3. Click "Send"
4. Copy the `token` from the response

#### Set JWT Token
1. Go to **Environments** → **TinderTrip API Environment**
2. Set `jwt_token` to the token you received
3. Save the environment

#### Login (Alternative)
1. Go to **Authentication** → **Login User**
2. Update the request body with your credentials
3. Click "Send"
4. Copy the `token` from the response

### 3. Test Protected Endpoints

Now you can test any protected endpoint:
- **User Profile** - Get, update, delete profile
- **User Preferences** - Manage availability and budget preferences
- **Events** - Create, join, swipe on events
- **Chat** - Send and receive messages
- **History** - View event history

### 4. Using Variables

The collection uses several variables:

- `{{base_url}}` - API base URL (default: http://localhost:9952)
- `{{jwt_token}}` - JWT authentication token
- `{{event_id}}` - Event ID for testing
- `{{room_id}}` - Chat room ID
- `{{history_id}}` - History record ID
- `{{user_id}}` - User ID
- `{{email}}` - User email
- `{{password}}` - User password

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
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

## Testing Workflow

### 1. Complete User Registration Flow
1. **Register User** - Create new account
2. **Set JWT Token** - Save token from response
3. **Get Profile** - Verify user profile
4. **Update Profile** - Add profile information
5. **Set Preferences** - Configure availability and budget

### 2. Event Management Flow
1. **Create Event** - Create a new event
2. **Get Events** - List all events
3. **Get Event by ID** - View specific event
4. **Join Event** - Join an event
5. **Swipe Event** - Like or pass on event

### 3. Chat Flow
1. **Get Chat Rooms** - List user's chat rooms
2. **Get Messages** - View messages in a room
3. **Send Message** - Send a new message

### 4. History Flow
1. **Get Event History** - View past events
2. **Mark Event Complete** - Mark event as completed

## Troubleshooting

### Common Issues

1. **401 Unauthorized**
   - Check if JWT token is set correctly
   - Verify token hasn't expired
   - Try logging in again

2. **404 Not Found**
   - Check if the API server is running
   - Verify the base URL is correct
   - Check if the endpoint exists

3. **500 Internal Server Error**
   - Check server logs
   - Verify database connection
   - Check if all required fields are provided

### Debug Tips

1. **Check Response Headers** - Look for error details
2. **Verify Request Body** - Ensure JSON is valid
3. **Check Environment Variables** - Make sure all variables are set
4. **Test with Health Check** - Verify server is running

## Environment Configuration

### Development
- **base_url**: `http://localhost:9952`
- **Database**: PostgreSQL (Docker)
- **Redis**: Optional

### Production
- **base_url**: `https://your-api-domain.com`
- **Database**: Production PostgreSQL
- **Redis**: Production Redis

## Collection Features

- **Organized by Feature** - Endpoints grouped by functionality
- **Pre-filled Examples** - Sample request bodies included
- **Variable Support** - Dynamic values for IDs and tokens
- **Environment Support** - Easy switching between environments
- **Documentation** - Each request includes description
- **Test Scripts** - Ready for automated testing

## Contributing

To add new endpoints or improve existing ones:

1. Update the collection JSON file
2. Add proper descriptions
3. Include example request bodies
4. Update this README
5. Test all endpoints

## Support

For issues or questions:
- Check the API documentation at `/swagger/index.html`
- Review server logs
- Check the GitHub repository
