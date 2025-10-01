#!/usr/bin/env bash
set -euo pipefail

# End-to-end backend flow (no browser/CORS):
# register -> login -> profile -> travel pref -> food pref -> availability
# -> budget -> create event -> like -> join -> confirm -> chat -> complete -> history

BASE="http://localhost:9952/api/v1"
EMAIL="flow_tester_$(date +%s)@example.com"
PASSWORD="Passw0rd!"
FIRST_NAME="Flow"
LAST_NAME="Tester"

echo "1) Health" >&2
curl -sS "http://localhost:9952/health" | jq .

echo "2) Register" >&2
curl -sS -X POST "$BASE/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\",\"first_name\":\"$FIRST_NAME\",\"last_name\":\"$LAST_NAME\",\"display_name\":\"$FIRST_NAME $LAST_NAME\"}" | jq .

echo "3) Login" >&2
LOGIN_JSON=$(curl -sS -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
echo "$LOGIN_JSON" | jq .
ACCESS_TOKEN=$(echo "$LOGIN_JSON" | jq -r '.access_token // .accessToken // .token // empty')
if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
  echo "ERROR: cannot extract access_token from login response" >&2; exit 1
fi
AUTH_HEADER="Authorization: Bearer $ACCESS_TOKEN"

echo "4) Create/Update Profile" >&2
curl -sS -X PUT "$BASE/users/profile" \
  -H "$AUTH_HEADER" -H "Content-Type: application/json" \
  -d '{
    "display_name": "Flow Tester",
    "bio": "Test end-to-end flow",
    "gender": "male",
    "languages": "en,th",
    "age": 25,
    "location": "Bangkok"
  }' | jq .

echo "5) Travel Preferences (bulk)" >&2
curl -sS -X PUT "$BASE/users/travel-preferences/bulk" \
  -H "$AUTH_HEADER" -H "Content-Type: application/json" \
  -d '{
    "travel_styles": ["cafe_dessert","coffee"]
  }' | jq .

echo "6) Food Preferences (bulk)" >&2
curl -sS -X PUT "$BASE/users/food-preferences/bulk" \
  -H "$AUTH_HEADER" -H "Content-Type: application/json" \
  -d '{
    "preferences": [
      {"food_category":"thai_food","preference_level":3},
      {"food_category":"japanese_food","preference_level":3},
      {"food_category":"international_food","preference_level":1}
    ]
  }' | jq .

echo "7) Availability Preferences" >&2
curl -sS -X PUT "$BASE/users/preferences/availability" \
  -H "$AUTH_HEADER" -H "Content-Type: application/json" \
  -d '{
    "mon": true, "tue": true, "wed": true, "thu": false, "fri": true,
    "sat": true, "sun": false,
    "all_day": false,
    "morning": true, "afternoon": true
  }' | jq .

echo "8) Budget Preferences" >&2
curl -sS -X PUT "$BASE/users/preferences/budget" \
  -H "$AUTH_HEADER" -H "Content-Type: application/json" \
  -d '{
    "meal_min": 100, "meal_max": 400,
    "daytrip_min": 500, "daytrip_max": 2000,
    "overnight_min": 1000, "overnight_max": 5000,
    "currency": "THB",
    "unlimited": false
  }' | jq .

echo "9) Create Trip (Event)" >&2
CREATE_EVENT_JSON=$(curl -sS -X POST "$BASE/events" \
  -H "$AUTH_HEADER" -H "Content-Type: application/json" \
  -d '{
    "title":"Flow Test Trip",
    "description":"E2E flow test event",
    "location":"Ayutthaya",
    "start_time":"2025-12-01T09:00:00Z",
    "end_time":"2025-12-01T18:00:00Z",
    "max_members": 5,
    "event_type": "meal"
  }')
echo "$CREATE_EVENT_JSON" | jq .
EVENT_ID=$(echo "$CREATE_EVENT_JSON" | jq -r '.id // .event.id // .data.id // empty')
if [ -z "$EVENT_ID" ] || [ "$EVENT_ID" = "null" ]; then
  echo "Trying to fetch an event id from list..." >&2
  EVENTS=$(curl -sS -H "$AUTH_HEADER" "$BASE/events")
  EVENT_ID=$(echo "$EVENTS" | jq -r '.events[0].id // .[0].id // empty')
fi
if [ -z "$EVENT_ID" ] || [ "$EVENT_ID" = "null" ]; then
  echo "ERROR: cannot get EVENT_ID" >&2; exit 1
fi
echo "EVENT_ID=$EVENT_ID" >&2

echo "10) Swipe/Like Event" >&2
curl -sS -X POST "$BASE/events/$EVENT_ID/swipe" \
  -H "$AUTH_HEADER" -H "Content-Type: application/json" \
  -d '{"event_id":"'$EVENT_ID'","direction":"like"}' | jq .

echo "11) Join Event" >&2
curl -sS -X POST "$BASE/events/$EVENT_ID/join" -H "$AUTH_HEADER" | jq .

echo "12) Confirm Event" >&2
curl -sS -X POST "$BASE/events/$EVENT_ID/confirm" -H "$AUTH_HEADER" | jq .

echo "13) Chat - rooms" >&2
ROOMS=$(curl -sS -H "$AUTH_HEADER" "$BASE/chat/rooms")
echo "$ROOMS" | jq .
ROOM_ID=$(echo "$ROOMS" | jq -r '.rooms[0].id // .[0].id // empty')
if [ -n "$ROOM_ID" ] && [ "$ROOM_ID" != "null" ]; then
  echo "13.1) Send message" >&2
  curl -sS -X POST "$BASE/chat/rooms/$ROOM_ID/messages" \
    -H "$AUTH_HEADER" -H "Content-Type: application/json" \
    -d '{"room_id":"'$ROOM_ID'","body":"Hello from E2E flow!","message_type":"text"}' | jq .

  echo "13.2) Get messages" >&2
  curl -sS -H "$AUTH_HEADER" "$BASE/chat/rooms/$ROOM_ID/messages" | jq .
else
  echo "WARN: No chat room yet; skipping send." >&2
fi

echo "14) Complete Event" >&2
curl -sS -X POST "$BASE/events/$EVENT_ID/complete" -H "$AUTH_HEADER" | jq .

echo "15) History" >&2
HISTORY=$(curl -sS -H "$AUTH_HEADER" "$BASE/history")
echo "$HISTORY" | jq .
HISTORY_ID=$(echo "$HISTORY" | jq -r '.history[0].id // .[0].id // empty')
if [ -n "$HISTORY_ID" ] && [ "$HISTORY_ID" != "null" ]; then
  echo "15.1) Mark history item complete" >&2
  curl -sS -X POST "$BASE/history/$HISTORY_ID/complete" -H "$AUTH_HEADER" | jq .
fi

echo "Done ✅" >&2

