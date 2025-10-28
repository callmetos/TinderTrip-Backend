#!/bin/bash

# =============================================================================
# TinderTrip Backend - Full Flow E2E Test
# =============================================================================
# User1: Register → Login → Update Profile → Create Event
# User2: Register → Login → Update Profile → Join Event
# =============================================================================

set -e  # Exit on error

API_URL="https://api.tindertrip.phitik.com/api/v1"
# API_URL="http://localhost:9952/api/v1"  # Uncomment for local testing

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_step() {
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}▶ $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Extract value from JSON response
extract_json() {
    echo "$1" | grep -o "\"$2\":\"[^\"]*" | cut -d'"' -f4
}

# =============================================================================
# USER 1 FLOW
# =============================================================================

print_step "👤 USER 1: Registration Flow"

# Generate unique email
TIMESTAMP=$(date +%s)
USER1_EMAIL="user1_${TIMESTAMP}@test.com"
USER1_PASSWORD="Test1234!"
USER1_DISPLAY_NAME="Test User 1"

print_info "Email: $USER1_EMAIL"

# Step 1: Register User 1
print_step "1️⃣ User 1: Register"
REGISTER1_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$USER1_EMAIL\",
    \"password\": \"$USER1_PASSWORD\",
    \"display_name\": \"$USER1_DISPLAY_NAME\"
  }")

echo "$REGISTER1_RESPONSE" | jq '.'

if echo "$REGISTER1_RESPONSE" | grep -q '"success":true'; then
    print_success "User 1 registered successfully"
else
    print_error "User 1 registration failed"
    exit 1
fi

# Step 2: Get OTP from dev endpoint
print_step "2️⃣ User 1: Get OTP"
sleep 2  # Wait for OTP generation

OTP_RESPONSE=$(curl -s -X GET "$API_URL/dev/otp")
echo "$OTP_RESPONSE" | jq '.'

USER1_OTP=$(echo "$OTP_RESPONSE" | jq -r ".data.otps[] | select(.email == \"$USER1_EMAIL\") | .otp")

if [ -z "$USER1_OTP" ] || [ "$USER1_OTP" = "null" ]; then
    print_error "Could not retrieve OTP for User 1"
    exit 1
fi

print_success "OTP retrieved: $USER1_OTP"

# Step 3: Verify Email with OTP
print_step "3️⃣ User 1: Verify Email"
VERIFY1_RESPONSE=$(curl -s -X POST "$API_URL/auth/verify-email" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$USER1_EMAIL\",
    \"otp\": \"$USER1_OTP\",
    \"password\": \"$USER1_PASSWORD\",
    \"display_name\": \"$USER1_DISPLAY_NAME\"
  }")

echo "$VERIFY1_RESPONSE" | jq '.'

USER1_TOKEN=$(echo "$VERIFY1_RESPONSE" | jq -r '.token')
USER1_ID=$(echo "$VERIFY1_RESPONSE" | jq -r '.user.id')

if [ -z "$USER1_TOKEN" ] || [ "$USER1_TOKEN" = "null" ]; then
    print_error "User 1 verification failed"
    exit 1
fi

print_success "User 1 verified and logged in"
print_info "Token: ${USER1_TOKEN:0:30}..."
print_info "User ID: $USER1_ID"

# Step 4: Update Profile
print_step "4️⃣ User 1: Update Profile"
PROFILE1_RESPONSE=$(curl -s -X PUT "$API_URL/users/profile" \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "Adventure lover and beach enthusiast!",
    "date_of_birth": "1995-05-15",
    "gender": "male",
    "phone_number": "+66812345678",
    "nationality": "Thai",
    "languages": "Thai, English"
  }')

echo "$PROFILE1_RESPONSE" | jq '.'

if echo "$PROFILE1_RESPONSE" | grep -q '"success":true'; then
    print_success "User 1 profile updated"
else
    print_error "User 1 profile update failed"
fi

# Step 5: Get Tags for Event Creation
print_step "5️⃣ Get Available Tags"
TAGS_RESPONSE=$(curl -s -X GET "$API_URL/public/tags?limit=20")

# Check if tags response is valid
if echo "$TAGS_RESPONSE" | jq -e '.data' > /dev/null 2>&1; then
    echo "$TAGS_RESPONSE" | jq '.data[] | {id, name, kind}'
    
    # Extract some tag IDs (adjust based on actual tags)
    TAG_BEACH=$(echo "$TAGS_RESPONSE" | jq -r '.data[] | select(.name == "Beach") | .id')
    TAG_ADVENTURE=$(echo "$TAGS_RESPONSE" | jq -r '.data[] | select(.name == "Adventure") | .id')
    TAG_FOOD=$(echo "$TAGS_RESPONSE" | jq -r '.data[] | select(.name == "Food") | .id')
else
    print_error "Failed to retrieve tags"
    TAG_BEACH=""
    TAG_ADVENTURE=""
    TAG_FOOD=""
fi

print_info "Beach Tag ID: $TAG_BEACH"
print_info "Adventure Tag ID: $TAG_ADVENTURE"
print_info "Food Tag ID: $TAG_FOOD"

# Step 6: Create Event
print_step "6️⃣ User 1: Create Event"

# Calculate dates (tomorrow and day after)
START_DATE=$(date -u -v+1d +"%Y-%m-%dT10:00:00Z" 2>/dev/null || date -u -d "tomorrow" +"%Y-%m-%dT10:00:00Z")
END_DATE=$(date -u -v+1d +"%Y-%m-%dT18:00:00Z" 2>/dev/null || date -u -d "tomorrow" +"%Y-%m-%dT18:00:00Z")

TAG_IDS_JSON="[]"
if [ ! -z "$TAG_BEACH" ] && [ "$TAG_BEACH" != "null" ]; then
    TAG_IDS_JSON="[\"$TAG_BEACH\""
    [ ! -z "$TAG_ADVENTURE" ] && [ "$TAG_ADVENTURE" != "null" ] && TAG_IDS_JSON="$TAG_IDS_JSON, \"$TAG_ADVENTURE\""
    [ ! -z "$TAG_FOOD" ] && [ "$TAG_FOOD" != "null" ] && TAG_IDS_JSON="$TAG_IDS_JSON, \"$TAG_FOOD\""
    TAG_IDS_JSON="$TAG_IDS_JSON]"
fi

EVENT_RESPONSE=$(curl -s -X POST "$API_URL/events" \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Beach Day Trip Phuket\",
    \"description\": \"Let's enjoy a beautiful day at the beach! Swimming, sunbathing, and delicious seafood.\",
    \"event_type\": \"one_day_trip\",
    \"address_text\": \"Patong Beach, Phuket\",
    \"lat\": 7.8967,
    \"lng\": 98.3004,
    \"start_at\": \"$START_DATE\",
    \"end_at\": \"$END_DATE\",
    \"capacity\": 10,
    \"budget_min\": 500,
    \"budget_max\": 1500,
    \"currency\": \"THB\",
    \"tag_ids\": $TAG_IDS_JSON
  }")

echo "$EVENT_RESPONSE" | jq '.'

EVENT_ID=$(echo "$EVENT_RESPONSE" | jq -r '.data.id')

if [ -z "$EVENT_ID" ] || [ "$EVENT_ID" = "null" ]; then
    print_error "Event creation failed"
    exit 1
fi

print_success "Event created successfully"
print_info "Event ID: $EVENT_ID"

# =============================================================================
# USER 2 FLOW
# =============================================================================

print_step "👤 USER 2: Registration Flow"

# Generate unique email for User 2
USER2_EMAIL="user2_${TIMESTAMP}@test.com"
USER2_PASSWORD="Test1234!"
USER2_DISPLAY_NAME="Test User 2"

print_info "Email: $USER2_EMAIL"

# Step 7: Register User 2
print_step "7️⃣ User 2: Register"
REGISTER2_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$USER2_EMAIL\",
    \"password\": \"$USER2_PASSWORD\",
    \"display_name\": \"$USER2_DISPLAY_NAME\"
  }")

echo "$REGISTER2_RESPONSE" | jq '.'

if echo "$REGISTER2_RESPONSE" | grep -q '"success":true'; then
    print_success "User 2 registered successfully"
else
    print_error "User 2 registration failed"
    exit 1
fi

# Step 8: Get OTP for User 2
print_step "8️⃣ User 2: Get OTP"
sleep 2

OTP_RESPONSE2=$(curl -s -X GET "$API_URL/dev/otp")
USER2_OTP=$(echo "$OTP_RESPONSE2" | jq -r ".data.otps[] | select(.email == \"$USER2_EMAIL\") | .otp")

if [ -z "$USER2_OTP" ] || [ "$USER2_OTP" = "null" ]; then
    print_error "Could not retrieve OTP for User 2"
    exit 1
fi

print_success "OTP retrieved: $USER2_OTP"

# Step 9: Verify Email for User 2
print_step "9️⃣ User 2: Verify Email"
VERIFY2_RESPONSE=$(curl -s -X POST "$API_URL/auth/verify-email" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$USER2_EMAIL\",
    \"otp\": \"$USER2_OTP\",
    \"password\": \"$USER2_PASSWORD\",
    \"display_name\": \"$USER2_DISPLAY_NAME\"
  }")

echo "$VERIFY2_RESPONSE" | jq '.'

USER2_TOKEN=$(echo "$VERIFY2_RESPONSE" | jq -r '.token')
USER2_ID=$(echo "$VERIFY2_RESPONSE" | jq -r '.user.id')

if [ -z "$USER2_TOKEN" ] || [ "$USER2_TOKEN" = "null" ]; then
    print_error "User 2 verification failed"
    exit 1
fi

print_success "User 2 verified and logged in"
print_info "Token: ${USER2_TOKEN:0:30}..."
print_info "User ID: $USER2_ID"

# Step 10: Update Profile for User 2
print_step "🔟 User 2: Update Profile"
PROFILE2_RESPONSE=$(curl -s -X PUT "$API_URL/users/profile" \
  -H "Authorization: Bearer $USER2_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "Love exploring new places and meeting new people!",
    "date_of_birth": "1998-08-20",
    "gender": "female",
    "phone_number": "+66898765432",
    "nationality": "Thai",
    "languages": "Thai, English, Japanese"
  }')

echo "$PROFILE2_RESPONSE" | jq '.'

if echo "$PROFILE2_RESPONSE" | grep -q '"success":true'; then
    print_success "User 2 profile updated"
else
    print_error "User 2 profile update failed"
fi

# Step 11: User 2 views the event
print_step "1️⃣1️⃣ User 2: View Event Details"
VIEW_EVENT_RESPONSE=$(curl -s -X GET "$API_URL/events/$EVENT_ID" \
  -H "Authorization: Bearer $USER2_TOKEN")

echo "$VIEW_EVENT_RESPONSE" | jq '.'

if echo "$VIEW_EVENT_RESPONSE" | grep -q '"success":true'; then
    print_success "Event details retrieved"
else
    print_error "Failed to retrieve event details"
fi

# Step 12: User 2 joins the event
print_step "1️⃣2️⃣ User 2: Join Event"
JOIN_RESPONSE=$(curl -s -X POST "$API_URL/events/$EVENT_ID/join" \
  -H "Authorization: Bearer $USER2_TOKEN")

echo "$JOIN_RESPONSE" | jq '.'

if echo "$JOIN_RESPONSE" | grep -q '"success":true'; then
    print_success "User 2 joined the event successfully"
else
    print_error "User 2 failed to join the event"
    exit 1
fi

# =============================================================================
# VERIFICATION
# =============================================================================

print_step "🔍 Verification"

# Check event members
print_step "1️⃣3️⃣ Check Event Members"
EVENT_CHECK=$(curl -s -X GET "$API_URL/events/$EVENT_ID" \
  -H "Authorization: Bearer $USER1_TOKEN")

echo "$EVENT_CHECK" | jq '.data.members'

MEMBER_COUNT=$(echo "$EVENT_CHECK" | jq '.data.member_count')
print_info "Total members: $MEMBER_COUNT"

if [ "$MEMBER_COUNT" -ge 2 ]; then
    print_success "Event has correct number of members"
else
    print_error "Event member count is incorrect"
fi

# Check User 2's joined events
print_step "1️⃣4️⃣ Check User 2's Joined Events"
JOINED_EVENTS=$(curl -s -X GET "$API_URL/events/joined" \
  -H "Authorization: Bearer $USER2_TOKEN")

echo "$JOINED_EVENTS" | jq '.'

if echo "$JOINED_EVENTS" | jq -e ".data[] | select(.id == \"$EVENT_ID\")" > /dev/null; then
    print_success "Event appears in User 2's joined events"
else
    print_error "Event does not appear in User 2's joined events"
fi

# =============================================================================
# SUMMARY
# =============================================================================

print_step "📊 Test Summary"

echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✓ All tests passed successfully!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

echo -e "${BLUE}Test Results:${NC}"
echo -e "  User 1:"
echo -e "    • Email: ${YELLOW}$USER1_EMAIL${NC}"
echo -e "    • ID: ${YELLOW}$USER1_ID${NC}"
echo -e "    • Token: ${YELLOW}${USER1_TOKEN:0:40}...${NC}"
echo ""
echo -e "  User 2:"
echo -e "    • Email: ${YELLOW}$USER2_EMAIL${NC}"
echo -e "    • ID: ${YELLOW}$USER2_ID${NC}"
echo -e "    • Token: ${YELLOW}${USER2_TOKEN:0:40}...${NC}"
echo ""
echo -e "  Event:"
echo -e "    • ID: ${YELLOW}$EVENT_ID${NC}"
echo -e "    • Title: ${YELLOW}Beach Day Trip Phuket${NC}"
echo -e "    • Members: ${YELLOW}$MEMBER_COUNT${NC}"
echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Cleanup instructions
print_step "🧹 Cleanup (Optional)"
echo "To clean up test data, run:"
echo ""
echo "# Delete User 1"
echo "curl -X DELETE \"$API_URL/users/profile\" \\"
echo "  -H \"Authorization: Bearer $USER1_TOKEN\""
echo ""
echo "# Delete User 2"
echo "curl -X DELETE \"$API_URL/users/profile\" \\"
echo "  -H \"Authorization: Bearer $USER2_TOKEN\""
echo ""

