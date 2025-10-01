# TinderTrip API Documentation

## Overview
TinderTrip Backend API provides endpoints for user management, event creation, chat functionality, and travel preferences.

## Base URL
- **Development**: `http://localhost:9952`
- **Swagger UI**: `http://localhost:9952/swagger/index.html`

## Authentication
Most endpoints require Bearer token authentication. Get your token from the login endpoint.

```bash
Authorization: Bearer <your_access_token>
```

## API Endpoints

### Health Check
- **GET** `/health` - Check API status

### Authentication
- **POST** `/api/v1/auth/register` - Register new user
- **POST** `/api/v1/auth/login` - Login user
- **POST** `/api/v1/auth/logout` - Logout user (requires auth)
- **POST** `/api/v1/auth/refresh` - Refresh access token (requires auth)

### User Profile
- **GET** `/api/v1/users/profile` - Get user profile (requires auth)
- **PUT** `/api/v1/users/profile` - Update user profile (requires auth)
- **DELETE** `/api/v1/users/profile` - Delete user profile (requires auth)

### Travel Preferences
- **GET** `/api/v1/users/travel-preferences` - Get travel preferences (requires auth)
- **PUT** `/api/v1/users/travel-preferences/bulk` - Update travel preferences (requires auth)

**Valid travel styles:**
- `cafe_dessert`, `bubble_tea`, `bakery_cake`, `bingsu_ice_cream`
- `coffee`, `matcha`, `pancakes`, `social_activity`
- `karaoke`, `gaming`, `movie`, `board_game`
- `outdoor_activity`, `party_celebration`, `swimming`, `skateboarding`

### Food Preferences
- **GET** `/api/v1/users/food-preferences` - Get food preferences (requires auth)
- **PUT** `/api/v1/users/food-preferences/bulk` - Update food preferences (requires auth)

**Valid food categories:**
- `thai_food`, `japanese_food`, `chinese_food`, `international_food`
- `halal_food`, `buffet`, `bbq_grill`

**Preference levels:**
- `1` - Dislike
- `2` - Neutral
- `3` - Like

### User Preferences
- **PUT** `/api/v1/users/preferences/availability` - Update availability (requires auth)
- **PUT** `/api/v1/users/preferences/budget` - Update budget (requires auth)

### Events
- **GET** `/api/v1/events` - Get events (requires auth)
- **POST** `/api/v1/events` - Create event (requires auth)
- **GET** `/api/v1/events/{id}` - Get event by ID (requires auth)
- **POST** `/api/v1/events/{id}/swipe` - Swipe event (requires auth)
- **POST** `/api/v1/events/{id}/join` - Join event (requires auth)
- **POST** `/api/v1/events/{id}/confirm` - Confirm event participation (requires auth)
- **POST** `/api/v1/events/{id}/complete` - Complete event (requires auth)

**Valid event types:**
- `meal` - Meal events
- `one_day_trip` - One day trips
- `overnight` - Overnight trips

### Chat
- **GET** `/api/v1/chat/rooms` - Get chat rooms (requires auth)
- **POST** `/api/v1/chat/rooms/{room_id}/messages` - Send message (requires auth)
- **GET** `/api/v1/chat/rooms/{room_id}/messages` - Get messages (requires auth)

### History
- **GET** `/api/v1/history` - Get user event history (requires auth)
- **POST** `/api/v1/history/{event_id}/complete` - Mark event as complete (requires auth)

## Request/Response Examples

### Register User
```json
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "password": "Passw0rd!",
  "display_name": "John Doe"
}
```

### Update Travel Preferences
```json
PUT /api/v1/users/travel-preferences/bulk
{
  "travel_styles": ["cafe_dessert", "coffee", "social_activity"]
}
```

### Update Food Preferences
```json
PUT /api/v1/users/food-preferences/bulk
{
  "preferences": [
    {"food_category": "thai_food", "preference_level": 3},
    {"food_category": "japanese_food", "preference_level": 3},
    {"food_category": "international_food", "preference_level": 1}
  ]
}
```

### Create Event
```json
POST /api/v1/events
{
  "title": "Coffee Meetup",
  "description": "Let's grab coffee together!",
  "event_type": "meal",
  "start_at": "2025-12-01T09:00:00Z",
  "end_at": "2025-12-01T11:00:00Z",
  "capacity": 4
}
```

### Send Chat Message
```json
POST /api/v1/chat/rooms/{room_id}/messages
{
  "room_id": "room-uuid",
  "body": "Hello everyone!",
  "message_type": "text"
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Validation failed",
  "message": "Invalid request data"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing token"
}
```

### 404 Not Found
```json
{
  "error": "Not found",
  "message": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "message": "Something went wrong. Please try again later."
}
```

## CORS Configuration
The API supports CORS for the following origins:
- `http://localhost:3000`
- `http://localhost:3001`
- `http://localhost:5173`
- `http://127.0.0.1:5501`
- `http://localhost:5501`

## Rate Limiting
- **Development**: Rate limiting is disabled
- **Production**: 100 requests per hour per IP

## Testing
Use the provided test script to verify all endpoints:
```bash
bash scripts/be_flow_test.sh
```

## Notes
- All timestamps are in ISO 8601 format (UTC)
- UUIDs are used for all entity IDs
- Password must be at least 6 characters
- Email addresses must be valid email format
