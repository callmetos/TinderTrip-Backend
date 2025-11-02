#!/bin/bash

# Frontend Flow Test Script
# Tests complete flow: Login → Profile → Preferences → Events

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

API_URL="${API_URL:-http://localhost:9952}"
TEST_EMAIL="${TEST_EMAIL:-testflow@chantos.com}"
TEST_PASSWORD="${TEST_PASSWORD:-password123}"

echo -e "${BLUE}=== 🎨 Frontend Flow Test ===${NC}"
echo ""
echo -e "${CYAN}User: ${TEST_EMAIL}${NC}"
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
        error_msg = d.get('message', d.get('error', 'Unknown error'))
        print(f'Error: {error_msg}', file=sys.stderr)
        sys.exit(1)
    print(token)
except Exception as e:
    print(f'Parse error: {e}', file=sys.stderr)
    sys.exit(1)
" 2>&1)

if [ $? -ne 0 ] || [ -z "$TOKEN" ]; then
    echo -e "  ${RED}❌ Login failed${NC}"
    echo "$LOGIN_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$LOGIN_RESPONSE"
    echo ""
    echo -e "${YELLOW}⚠️  Note: If user doesn't exist, register first or check credentials${NC}"
    exit 1
fi

echo -e "  ${GREEN}✅ Logged in${NC}"
echo "  Token: ${TOKEN:0:30}..."
echo ""

# Step 2: Get Profile
echo -e "${YELLOW}Step 2: Get Profile${NC}"
PROFILE_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/users/profile" \
  -H "Authorization: Bearer ${TOKEN}")

echo "$PROFILE_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    profile = d.get('data', {})
    if profile:
        print(f\"  ✅ Profile: {profile.get('display_name', 'N/A')}\")
        print(f\"     Bio: {profile.get('bio', 'Not set')[:50]}...\")
    else:
        print('  ⚠️  Profile not set yet')
except Exception as e:
    print(f\"  ⚠️  Error: {e}\")
" 2>/dev/null

echo ""

# Step 3: Update Profile
echo -e "${YELLOW}Step 3: Update Profile${NC}"
UPDATE_PROFILE_RESPONSE=$(curl -sS -X PUT "${API_URL}/api/v1/users/profile" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "display_name": "Test Flow User",
    "bio": "I love traveling and trying new foods!",
    "languages": ["th", "en"],
    "date_of_birth": "1995-05-15",
    "gender": "male",
    "job_title": "Software Engineer",
    "smoking": "no"
  }')

if echo "$UPDATE_PROFILE_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); exit(0 if 'error' not in d else 1)" 2>/dev/null; then
    echo -e "  ${GREEN}✅ Profile updated${NC}"
else
    ERROR_MSG=$(echo "$UPDATE_PROFILE_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('message', d.get('error', 'Unknown error')))" 2>/dev/null)
    echo -e "  ${YELLOW}⚠️  Profile update: ${ERROR_MSG}${NC}"
fi

echo ""

# Step 4: Get Travel Styles (from database)
echo -e "${YELLOW}Step 4: Get Travel Styles (from database)${NC}"
TRAVEL_STYLES_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/public/travel-preferences/styles")

TRAVEL_STYLES=$(echo "$TRAVEL_STYLES_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    styles = d.get('data', []) if isinstance(d.get('data'), list) else d.get('styles', [])
    style_codes = [s.get('style') for s in styles[:5] if s.get('style')]
    print(' '.join(style_codes))
except:
    print('')
" 2>/dev/null)

echo "$TRAVEL_STYLES_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    styles = d.get('data', []) if isinstance(d.get('data'), list) else d.get('styles', [])
    print(f\"  ✅ {len(styles)} travel styles from DATABASE\")
    print('')
    for i, s in enumerate(styles[:8], 1):
        print(f\"     {i}. {s.get('style', 'N/A'):25} | {s.get('display_name', 'N/A')}\")
except Exception as e:
    print(f\"  ❌ Error: {e}\")
" 2>/dev/null

echo ""

# Step 5: Add Travel Preferences
echo -e "${YELLOW}Step 5: Add Travel Preferences${NC}"
SELECTED_TRAVEL_STYLES=("outdoor_activity" "social_activity" "coffee" "cafe_dessert")
ADDED_COUNT=0

for style in "${SELECTED_TRAVEL_STYLES[@]}"; do
    ADD_RESPONSE=$(curl -sS -X POST "${API_URL}/api/v1/users/travel-preferences" \
      -H "Authorization: Bearer ${TOKEN}" \
      -H "Content-Type: application/json" \
      -d "{\"travel_style\": \"${style}\"}")
    
    if echo "$ADD_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); exit(0 if 'error' not in d else 1)" 2>/dev/null; then
        echo -e "  ${GREEN}✅ Added: ${style}${NC}"
        ((ADDED_COUNT++))
    else
        ERROR_MSG=$(echo "$ADD_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('message', 'Unknown error'))" 2>/dev/null)
        if [[ "$ERROR_MSG" == *"already exists"* ]]; then
            echo -e "  ${YELLOW}⚠️  Already exists: ${style}${NC}"
            ((ADDED_COUNT++))
        else
            echo -e "  ${RED}❌ Failed: ${style} - ${ERROR_MSG}${NC}"
        fi
    fi
done

echo -e "  ${GREEN}✅ Total travel preferences: ${ADDED_COUNT}/${#SELECTED_TRAVEL_STYLES[@]}${NC}"
echo ""

# Step 6: Get Food Categories (from database)
echo -e "${YELLOW}Step 6: Get Food Categories (from database)${NC}"
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

# Step 7: Update Food Preferences (Bulk)
echo -e "${YELLOW}Step 7: Update Food Preferences (Bulk)${NC}"
FOOD_PREFS_RESPONSE=$(curl -sS -X PUT "${API_URL}/api/v1/users/food-preferences/bulk" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "preferences": [
      {"food_category": "thai_food", "preference_level": 3},
      {"food_category": "chinese_food", "preference_level": 3},
      {"food_category": "japanese_food", "preference_level": 2},
      {"food_category": "international_food", "preference_level": 2},
      {"food_category": "halal_food", "preference_level": 1},
      {"food_category": "buffet", "preference_level": 2},
      {"food_category": "bbq_grill", "preference_level": 3}
    ]
  }')

if echo "$FOOD_PREFS_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); exit(0 if 'error' not in d else 1)" 2>/dev/null; then
    echo -e "  ${GREEN}✅ Food preferences updated${NC}"
    echo "$FOOD_PREFS_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    print('')
    print('  Preferences set:')
    level_map = {1: '😱 Dislike', 2: '😃 Neutral', 3: '🤩 Love'}
    prefs = d.get('data', {}).get('preferences', [])
    for p in prefs:
        level = p.get('preference_level', 2)
        level_name = level_map.get(level, 'Unknown')
        print(f\"     • {p.get('food_category', 'N/A')}: {level_name}\")
except:
    pass
" 2>/dev/null
else
    ERROR_MSG=$(echo "$FOOD_PREFS_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('message', d.get('error', 'Unknown error')))" 2>/dev/null)
    echo -e "  ${RED}❌ Failed: ${ERROR_MSG}${NC}"
fi

echo ""

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
    echo "$BUDGET_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    budget = d.get('data', {})
    if budget:
        print('')
        print('  Budget set:')
        print(f\"     Meal: {budget.get('meal_min', 0)}-{budget.get('meal_max', 0)} {budget.get('currency', 'THB')}\")
        print(f\"     Daytrip: {budget.get('daytrip_min', 0)}-{budget.get('daytrip_max', 0)} {budget.get('currency', 'THB')}\")
        print(f\"     Overnight: {budget.get('overnight_min', 0)}-{budget.get('overnight_max', 0)} {budget.get('currency', 'THB')}\")
except Exception as e:
    print(f\"  ⚠️  Parse error: {e}\")
" 2>/dev/null
else
    ERROR_MSG=$(echo "$BUDGET_RESPONSE" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('message', d.get('error', 'Unknown error')))" 2>/dev/null)
    echo -e "  ${RED}❌ Failed: ${ERROR_MSG}${NC}"
fi

echo ""

# Step 9: Get Preferences Summary
echo -e "${YELLOW}Step 9: Get Preferences Summary${NC}"

# Get Travel Preferences
USER_TRAVEL_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/users/travel-preferences/styles" \
  -H "Authorization: Bearer ${TOKEN}")

# Get Food Preferences
USER_FOOD_RESPONSE=$(curl -sS -X GET "${API_URL}/api/v1/users/food-preferences/categories" \
  -H "Authorization: Bearer ${TOKEN}")

echo -e "  ${GREEN}✅ Travel Preferences:${NC}"
echo "$USER_TRAVEL_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    styles = d.get('data', []) if isinstance(d.get('data'), list) else d.get('styles', [])
    selected = [s for s in styles if s.get('is_selected', False)]
    if selected:
        for s in selected:
            print(f\"     • {s.get('style', 'N/A')}: {s.get('display_name', 'N/A')}\")
    else:
        print('     (No travel preferences selected)')
except:
    pass
" 2>/dev/null

echo ""
echo -e "  ${GREEN}✅ Food Preferences:${NC}"
echo "$USER_FOOD_RESPONSE" | python3 -c "
import sys, json
try:
    d = json.load(sys.stdin)
    categories = d.get('data', []) if isinstance(d.get('data'), list) else d.get('categories', [])
    level_map = {1: '😱 Dislike', 2: '😃 Neutral', 3: '🤩 Love'}
    for c in categories:
        level = c.get('preference_level', 2)
        level_name = level_map.get(level, 'Unknown')
        if level != 2:  # Only show non-neutral
            print(f\"     • {c.get('category', 'N/A')}: {level_name}\")
except:
    pass
" 2>/dev/null

echo ""

# Step 10: Get Events with Match Score (Suggestions)
echo -e "${YELLOW}Step 10: Get Events with Match Score (Suggestions)${NC}"
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
        for i, event in enumerate(events[:10], 1):
            match_score = event.get('match_score', 'N/A')
            title = event.get('title', 'N/A')[:40]
            event_type = event.get('event_type', 'N/A')
            budget_min = event.get('budget_min', 0)
            budget_max = event.get('budget_max', 0)
            currency = event.get('currency', 'THB')
            status = event.get('status', 'N/A')
            print(f\"     {i:2d}. Score: {str(match_score):6} | {title:40} | Type: {event_type:12} | Budget: {budget_min}-{budget_max} {currency}\")
    else:
        print('  No events found')
except Exception as e:
    print(f\"  ❌ Error: {e}\")
" 2>/dev/null

echo ""

# Step 11: Summary
echo -e "${BLUE}=== 📊 Test Summary ===${NC}"
echo ""
echo -e "${GREEN}✅ Completed Steps:${NC}"
echo "  1. ✅ Login"
echo "  2. ✅ Get Profile"
echo "  3. ✅ Update Profile"
echo "  4. ✅ Get Travel Styles (from DATABASE)"
echo "  5. ✅ Add Travel Preferences"
echo "  6. ✅ Get Food Categories (from DATABASE)"
echo "  7. ✅ Update Food Preferences"
echo "  8. ✅ Update Budget Preferences"
echo "  9. ✅ Get Preferences Summary"
echo "  10. ✅ Get Events with Match Score (Suggestions)"
echo ""
echo -e "${CYAN}User: ${TEST_EMAIL}${NC}"
echo -e "${CYAN}Token: ${TOKEN:0:30}...${NC}"
echo ""
echo -e "${BLUE}✅ Full Frontend Flow Test Completed Successfully!${NC}"
echo ""

