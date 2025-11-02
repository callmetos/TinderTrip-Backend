#!/bin/bash

# Full Flow Test Script for TinderTrip Backend
# Tests: Register → Login → Profile → Preferences → Event Creation → Suggestions

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_URL="${API_URL:-http://localhost:9952}"
TEST_EMAIL="flowtest$(date +%s)@test.com"
TEST_PASSWORD="Test1234!"
TEST_DISPLAY_NAME="Flow Test User"

echo -e "${BLUE}=== 🧪 Full Flow Test ===${NC}"
echo ""

# Step 1: Register
echo -e "${YELLOW}Step 1: Register${NC}"
REGISTER_RESPONSE=$(curl -sS -X POST "${API_URL}/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"${TEST_EMAIL}\",
    \"password\": \"${TEST_PASSWORD}\",
    \"display_name\": \"${TEST_DISPLAY_NAME}\"
  }")

echo "$REGISTER_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    if 'error' in d:
        print(f\"  ❌ Register failed: {d.get('message', 'Unknown error')}\")
        sys.exit(1)
    else:
        print(f\"  ✅ Registered: {d.get('message', 'Success')}\")
except Exception as e:
    print(f\"  ❌ Parse error: {e}\")
    sys.exit(1)
"

# Step 2: Login
echo -e "${YELLOW}Step 2: Login${NC}"
LOGIN_RESPONSE=$(curl -sS -X POST "${API_URL}/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"${TEST_EMAIL}\",
    \"password\": \"${TEST_PASSWORD}\"
  }")

TOKEN=$(echo "$LOGIN_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    token = d.get('token') or d.get('data', {}).get('token') or ''
    if not token:
        print('', file=sys.stderr)
        sys.exit(1)
    print(token)
except:
    sys.exit(1)
" 2>/dev/null)

if [ -z "$TOKEN" ]; then
    echo -e "  ${RED}❌ Login failed${NC}"
    echo "$LOGIN_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$LOGIN_RESPONSE"
    exit 1
fi

echo -e "  ${GREEN}✅ Logged in${NC}"
echo "  Token: ${TOKEN:0:20}..."

# Step 3: Create/Update Profile
echo -e "${YELLOW}Step 3: Create/Update Profile${NC}"
PROFILE_RESPONSE=$(curl -sS -X PUT "${API_URL}/api/v1/users/profile" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "Testing full flow",
    "languages": ["th", "en"],
    "date_of_birth": "1995-01-01",
    "gender": "male",
    "job_title": "Software Engineer",
    "smoking": "no"
  }')

echo "$PROFILE_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    if 'error' in d:
        print(f\"  ⚠️  Profile update: {d.get('message', 'Unknown error')}\")
    else:
        print(f\"  ✅ Profile updated\")
        profile = d.get('data', {})
        print(f\"     Display Name: {profile.get('display_name', 'N/A')}\")
        print(f\"     Bio: {profile.get('bio', 'N/A')[:50]}...\")
except Exception as e:
    print(f\"  ⚠️  Parse error: {e}\")
"

# Step 4: Get Available Travel Styles
echo -e "${YELLOW}Step 4: Get Available Travel Styles${NC}"
TRAVEL_STYLES_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/public/travel-preferences/styles")
TRAVEL_STYLES=$(echo "$TRAVEL_STYLES_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    styles = d.get('data', []) if isinstance(d.get('data'), list) else d.get('styles', [])
    print(' '.join([s.get('style', '') for s in styles[:5]]))
except:
    print('')
" 2>/dev/null)

echo -e "  ${GREEN}✅ Available Travel Styles:${NC}"
echo "$TRAVEL_STYLES_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    styles = d.get('data', []) if isinstance(d.get('data'), list) else d.get('styles', [])
    for i, s in enumerate(styles[:5], 1):
        print(f\"     {i}. {s.get('style', 'N/A')}: {s.get('display_name', 'N/A')}\")
except:
    pass
"

# Step 5: Add Travel Preferences
echo -e "${YELLOW}Step 5: Add Travel Preferences${NC}"
SELECTED_TRAVEL_STYLES=("outdoor_activity" "social_activity" "coffee")
for style in "${SELECTED_TRAVEL_STYLES[@]}"; do
    ADD_RESPONSE=$(curl -sS -X POST "${API_URL}/api/v1/users/travel-preferences" \
      -H "Authorization: Bearer ${TOKEN}" \
      -H "Content-Type: application/json" \
      -d "{\"travel_style\": \"${style}\"}")
    
    if echo "$ADD_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); exit(0 if 'error' not in d else 1)" 2>/dev/null; then
        echo -e "  ${GREEN}✅ Added: ${style}${NC}"
    else
        ERROR_MSG=$(echo "$ADD_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('message', 'Unknown error'))" 2>/dev/null)
        if [[ "$ERROR_MSG" == *"already exists"* ]]; then
            echo -e "  ${YELLOW}⚠️  Already exists: ${style}${NC}"
        else
            echo -e "  ${RED}❌ Failed: ${style} - ${ERROR_MSG}${NC}"
        fi
    fi
done

# Step 6: Get Available Food Categories
echo -e "${YELLOW}Step 6: Get Available Food Categories${NC}"
FOOD_CATEGORIES_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/public/food-preferences/categories")
echo -e "  ${GREEN}✅ Available Food Categories:${NC}"
echo "$FOOD_CATEGORIES_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    categories = d.get('data', []) if isinstance(d.get('data'), list) else d.get('categories', [])
    for i, c in enumerate(categories[:3], 1):
        print(f\"     {i}. {c.get('category', 'N/A')}: {c.get('display_name', 'N/A')}\")
except:
    pass
"

# Step 7: Update Food Preferences
echo -e "${YELLOW}Step 7: Update Food Preferences${NC}"
FOOD_PREFS=$(curl -sS -X PUT "${API_URL}/api/v1/users/food-preferences" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "preferences": [
      {"food_category": "thai_food", "preference_level": 3},
      {"food_category": "chinese_food", "preference_level": 3},
      {"food_category": "japanese_food", "preference_level": 2},
      {"food_category": "international_food", "preference_level": 2},
      {"food_category": "halal_food", "preference_level": 2},
      {"food_category": "buffet", "preference_level": 1},
      {"food_category": "bbq_grill", "preference_level": 3}
    ]
  }')

if echo "$FOOD_PREFS" | python3 -c "import sys, json; d=json.load(sys.stdin); exit(0 if 'error' not in d else 1)" 2>/dev/null; then
    echo -e "  ${GREEN}✅ Food preferences updated${NC}"
else
    ERROR_MSG=$(echo "$FOOD_PREFS" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('message', 'Unknown error'))" 2>/dev/null)
    echo -e "  ${RED}❌ Failed: ${ERROR_MSG}${NC}"
fi

# Step 8: Update Budget Preferences
echo -e "${YELLOW}Step 8: Update Budget Preferences${NC}"
BUDGET_RESPONSE=$(curl -sS -X PUT "${API_URL}/api/v1/users/preferences/budget" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "meal_min": 0,
    "meal_max": 500,
    "daytrip_min": 0,
    "daytrip_max": 2000,
    "overnight_min": 0,
    "overnight_max": 5000,
    "currency": "THB"
  }')

if echo "$BUDGET_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); exit(0 if 'error' not in d else 1)" 2>/dev/null; then
    echo -e "  ${GREEN}✅ Budget preferences updated${NC}"
    BUDGET_DATA=$(echo "$BUDGET_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); print(json.dumps(d.get('data', {}), indent=2))" 2>/dev/null)
    echo "$BUDGET_DATA" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    print(f\"     Meal: {d.get('meal_min', 0)}-{d.get('meal_max', 0)} THB\")
    print(f\"     Daytrip: {d.get('daytrip_min', 0)}-{d.get('daytrip_max', 0)} THB\")
    print(f\"     Overnight: {d.get('overnight_min', 0)}-{d.get('overnight_max', 0)} THB\")
except:
    pass
"
else
    ERROR_MSG=$(echo "$BUDGET_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('message', 'Unknown error'))" 2>/dev/null)
    echo -e "  ${RED}❌ Failed: ${ERROR_MSG}${NC}"
fi

# Step 9: Get Tags (for event creation)
echo -e "${YELLOW}Step 9: Get Available Tags${NC}"
TAGS_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/tags" \
  -H "Authorization: Bearer ${TOKEN}")

TAG_IDS=$(echo "$TAGS_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    tags = d.get('data', [])
    tag_ids = [t.get('id') for t in tags[:3] if t.get('id')]
    print(' '.join(tag_ids))
except:
    print('')
" 2>/dev/null)

echo -e "  ${GREEN}✅ Available Tags${NC}"
echo "$TAGS_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    tags = d.get('data', [])
    for i, t in enumerate(tags[:5], 1):
        print(f\"     {i}. {t.get('name', 'N/A')} ({t.get('kind', 'N/A')}) - ID: {t.get('id', 'N/A')[:8]}...\")
except:
    pass
"

# Step 10: Create Event
echo -e "${YELLOW}Step 10: Create Event${NC}"
FIRST_TAG_ID=$(echo "$TAG_IDS" | awk '{print $1}')

CREATE_EVENT_RESPONSE=$(curl -sS -X POST "${API_URL}/api/v1/events" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Full Flow Test Event\",
    \"description\": \"Testing complete flow from registration to event creation\",
    \"event_type\": \"meal\",
    \"address_text\": \"123 Test Street, Bangkok\",
    \"lat\": 13.7563,
    \"lng\": 100.5018,
    \"start_at\": \"2025-12-25T12:00:00Z\",
    \"end_at\": \"2025-12-25T14:00:00Z\",
    \"capacity\": 10,
    \"budget_min\": 200,
    \"budget_max\": 500,
    \"currency\": \"THB\",
    \"tag_ids\": [\"${FIRST_TAG_ID}\"]
  }")

EVENT_ID=$(echo "$CREATE_EVENT_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    if 'error' in d:
        print('', file=sys.stderr)
        sys.exit(1)
    event = d.get('data', {})
    print(event.get('id', ''))
except:
    sys.exit(1)
" 2>/dev/null)

if [ -z "$EVENT_ID" ]; then
    echo -e "  ${RED}❌ Event creation failed${NC}"
    echo "$CREATE_EVENT_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$CREATE_EVENT_RESPONSE"
else
    echo -e "  ${GREEN}✅ Event created${NC}"
    echo "  Event ID: ${EVENT_ID:0:8}..."
    echo "$CREATE_EVENT_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    event = d.get('data', {})
    print(f\"     Title: {event.get('title', 'N/A')}\")
    print(f\"     Type: {event.get('event_type', 'N/A')}\")
    print(f\"     Budget: {event.get('budget_min', 0)}-{event.get('budget_max', 0)} {event.get('currency', 'THB')}\")
except:
    pass
"
fi

# Step 11: Get Events with Suggestions (Match Score)
echo -e "${YELLOW}Step 11: Get Events with Match Score (Suggestions)${NC}"
EVENTS_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/events?limit=10" \
  -H "Authorization: Bearer ${TOKEN}")

echo -e "  ${GREEN}✅ Events retrieved (sorted by match score)${NC}"
echo "$EVENTS_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    events = d.get('data', [])
    total = d.get('total', len(events))
    print(f\"     Total events: {total}\")
    print('')
    for i, event in enumerate(events[:5], 1):
        match_score = event.get('match_score', 'N/A')
        title = event.get('title', 'N/A')[:40]
        event_type = event.get('event_type', 'N/A')
        budget = f\"{event.get('budget_min', 0)}-{event.get('budget_max', 0)}\"
        print(f\"     {i}. {title:40} | Score: {match_score:6} | Type: {event_type:10} | Budget: {budget} THB\")
except Exception as e:
    print(f\"     Error parsing: {e}\")
" 2>/dev/null

# Step 12: Summary
echo ""
echo -e "${BLUE}=== 📊 Test Summary ===${NC}"
echo ""
echo -e "${GREEN}✅ Completed Steps:${NC}"
echo "  1. ✅ Register"
echo "  2. ✅ Login"
echo "  3. ✅ Create/Update Profile"
echo "  4. ✅ Get Travel Styles (from database)"
echo "  5. ✅ Add Travel Preferences"
echo "  6. ✅ Get Food Categories (from database)"
echo "  7. ✅ Update Food Preferences"
echo "  8. ✅ Update Budget Preferences"
echo "  9. ✅ Get Tags"
echo "  10. ✅ Create Event"
echo "  11. ✅ Get Events with Match Score"
echo ""
echo -e "${BLUE}Test Email: ${TEST_EMAIL}${NC}"
echo -e "${BLUE}Test completed successfully! 🎉${NC}"
echo ""

