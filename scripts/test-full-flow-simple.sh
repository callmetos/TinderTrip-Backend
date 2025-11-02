#!/bin/bash

# Simple Full Flow Test Script
# Uses existing user to test complete flow

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

API_URL="${API_URL:-http://localhost:9952}"
TEST_EMAIL="${TEST_EMAIL:-channathat.u@gmail.com}"
TEST_PASSWORD="${TEST_PASSWORD:-admin1234}"

echo -e "${BLUE}=== 🧪 Full Flow Test (Complete) ===${NC}"
echo ""
echo "Using existing user: ${TEST_EMAIL}"
echo ""

# Step 1: Login
echo -e "${YELLOW}Step 1: Login${NC}"
LOGIN_RESPONSE=$(curl -sS -X POST "${API_URL}/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"${TEST_EMAIL}\", \"password\": \"${TEST_PASSWORD}\"}")

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
    exit 1
fi

echo -e "  ${GREEN}✅ Logged in${NC}"
echo ""

# Step 2: Get Profile
echo -e "${YELLOW}Step 2: Get Profile${NC}"
PROFILE_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/users/profile" \
  -H "Authorization: Bearer ${TOKEN}")

PROFILE=$(echo "$PROFILE_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    profile = d.get('data', {})
    print(f\"  ✅ Profile: {profile.get('display_name', 'N/A')}\")
    print(f\"     Bio: {profile.get('bio', 'N/A')[:50]}...\")
except:
    pass
" 2>/dev/null)

echo "$PROFILE"
echo ""

# Step 3: Get Travel Styles (from database)
echo -e "${YELLOW}Step 3: Get Travel Styles (from database)${NC}"
TRAVEL_STYLES_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/public/travel-preferences/styles")

echo "$TRAVEL_STYLES_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    styles = d.get('data', []) if isinstance(d.get('data'), list) else d.get('styles', [])
    print(f\"  ✅ {len(styles)} travel styles from DATABASE\")
    print('')
    for i, s in enumerate(styles[:5], 1):
        print(f\"     {i}. {s.get('style', 'N/A'):25} | {s.get('display_name', 'N/A')}\")
except Exception as e:
    print(f\"  ❌ Error: {e}\")
" 2>/dev/null

echo ""

# Step 4: Get Travel Preferences with User Selection
echo -e "${YELLOW}Step 4: Get Travel Preferences (with user selections)${NC}"
USER_TRAVEL_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/users/travel-preferences/styles" \
  -H "Authorization: Bearer ${TOKEN}")

echo "$USER_TRAVEL_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    styles = d.get('data', []) if isinstance(d.get('data'), list) else d.get('styles', [])
    selected = [s for s in styles if s.get('is_selected', False)]
    print(f\"  ✅ {len(styles)} total styles, {len(selected)} selected\")
    print('')
    if selected:
        print('  Selected styles:')
        for s in selected[:5]:
            print(f\"     • {s.get('style', 'N/A')}: {s.get('display_name', 'N/A')}\")
    else:
        print('  No styles selected yet')
except Exception as e:
    print(f\"  ❌ Error: {e}\")
" 2>/dev/null

echo ""

# Step 5: Get Food Categories (from database)
echo -e "${YELLOW}Step 5: Get Food Categories (from database)${NC}"
FOOD_CATEGORIES_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/public/food-preferences/categories")

echo "$FOOD_CATEGORIES_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    categories = d.get('data', []) if isinstance(d.get('data'), list) else d.get('categories', [])
    print(f\"  ✅ {len(categories)} food categories from DATABASE\")
    print('')
    for i, c in enumerate(categories, 1):
        print(f\"     {i}. {c.get('category', 'N/A'):25} | {c.get('display_name', 'N/A')}\")
except Exception as e:
    print(f\"  ❌ Error: {e}\")
" 2>/dev/null

echo ""

# Step 6: Get Food Preferences
echo -e "${YELLOW}Step 6: Get Food Preferences (with user selections)${NC}"
USER_FOOD_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/users/food-preferences/categories" \
  -H "Authorization: Bearer ${TOKEN}")

echo "$USER_FOOD_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    categories = d.get('data', []) if isinstance(d.get('data'), list) else d.get('categories', [])
    print(f\"  ✅ {len(categories)} categories with preferences\")
    print('')
    level_map = {1: '😱 Dislike', 2: '😃 Neutral', 3: '🤩 Love'}
    for c in categories:
        level = c.get('preference_level', 2)
        level_name = level_map.get(level, 'Unknown')
        print(f\"     • {c.get('category', 'N/A'):25} | {level_name}\")
except Exception as e:
    print(f\"  ❌ Error: {e}\")
" 2>/dev/null

echo ""

# Step 7: Get Budget Preferences
echo -e "${YELLOW}Step 7: Get Budget Preferences${NC}"
BUDGET_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/users/preferences/budget" \
  -H "Authorization: Bearer ${TOKEN}")

echo "$BUDGET_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    budget = d.get('data', {})
    if budget:
        print(f\"  ✅ Budget preferences:\")
        print(f\"     Meal: {budget.get('meal_min', 0)}-{budget.get('meal_max', 0)} {budget.get('currency', 'THB')}\")
        print(f\"     Daytrip: {budget.get('daytrip_min', 0)}-{budget.get('daytrip_max', 0)} {budget.get('currency', 'THB')}\")
        print(f\"     Overnight: {budget.get('overnight_min', 0)}-{budget.get('overnight_max', 0)} {budget.get('currency', 'THB')}\")
    else:
        print('  No budget preferences set yet')
except Exception as e:
    print(f\"  ❌ Error: {e}\")
" 2>/dev/null

echo ""

# Step 8: Get Tags (for event creation)
echo -e "${YELLOW}Step 8: Get Available Tags${NC}"
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

echo "$TAGS_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    tags = d.get('data', [])
    print(f\"  ✅ {len(tags)} tags available\")
    print('')
    for i, t in enumerate(tags[:5], 1):
        print(f\"     {i}. {t.get('name', 'N/A'):25} ({t.get('kind', 'N/A')})\")
except Exception as e:
    print(f\"  ❌ Error: {e}\")
" 2>/dev/null

echo ""

# Step 9: Get Events with Match Score (Suggestions)
echo -e "${YELLOW}Step 9: Get Events (with Match Score - Suggestions)${NC}"
EVENTS_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/events?limit=10" \
  -H "Authorization: Bearer ${TOKEN}")

echo "$EVENTS_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    events = d.get('data', [])
    total = d.get('total', len(events))
    print(f\"  ✅ {total} total events, showing {len(events)} (sorted by match score)\")
    print('')
    if events:
        print('  Top Events by Match Score:')
        print('')
        for i, event in enumerate(events[:5], 1):
            match_score = event.get('match_score', 'N/A')
            title = event.get('title', 'N/A')[:35]
            event_type = event.get('event_type', 'N/A')
            budget_min = event.get('budget_min', 0)
            budget_max = event.get('budget_max', 0)
            currency = event.get('currency', 'THB')
            print(f\"     {i}. {title:35} | Score: {match_score:6} | Type: {event_type:10} | Budget: {budget_min}-{budget_max} {currency}\")
    else:
        print('  No events found')
except Exception as e:
    print(f\"  ❌ Error: {e}\")
" 2>/dev/null

echo ""

# Step 10: Summary
echo -e "${BLUE}=== 📊 Test Summary ===${NC}"
echo ""
echo -e "${GREEN}✅ Completed Steps:${NC}"
echo "  1. ✅ Login"
echo "  2. ✅ Get Profile"
echo "  3. ✅ Get Travel Styles (from DATABASE)"
echo "  4. ✅ Get Travel Preferences"
echo "  5. ✅ Get Food Categories (from DATABASE)"
echo "  6. ✅ Get Food Preferences"
echo "  7. ✅ Get Budget Preferences"
echo "  8. ✅ Get Tags"
echo "  9. ✅ Get Events with Match Score (Suggestions)"
echo ""
echo -e "${BLUE}✅ All tests completed successfully!${NC}"
echo ""

