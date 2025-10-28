# Integration Tests

## Overview
Integration tests for TinderTrip Backend that test complete workflows using real database connections.

## Tests Included

### 1. Complete Auth Flow Test (`TestCompleteAuthFlow`)
Tests the entire authentication workflow from registration to login:

**Steps:**
1. ✅ Register with email
2. ✅ Verify OTP with invalid code (should fail)
3. ✅ Verify OTP with valid code (should succeed)
4. ✅ Attempt duplicate registration (should fail with 409)
5. ✅ Login with wrong password (should fail)
6. ✅ Login with correct credentials (should succeed)
7. ✅ Access protected endpoint with valid token
8. ✅ Access protected endpoint without token (should fail)
9. ✅ Access protected endpoint with invalid token (should fail)
10. ✅ Logout

**What it tests:**
- User registration flow
- OTP generation and verification
- Email verification
- Duplicate email detection (409 Conflict)
- Password authentication
- JWT token generation
- Protected endpoint access
- Token validation
- User session management

### 2. Resend Verification Flow Test (`TestResendVerificationFlow`)
Tests the OTP resend functionality:

**Steps:**
1. ✅ Initial registration
2. ✅ Resend verification OTP
3. ✅ Verify with new OTP

**What it tests:**
- OTP regeneration
- Resend verification email
- Multiple OTP handling

## Running Tests

### Prerequisites
- PostgreSQL database running
- Environment variables configured (`.env` file)
- Database tables created

### Run All Integration Tests
```bash
go test -v ./tests/integration/...
```

### Run Specific Test
```bash
# Run complete auth flow test
go test -v -run TestCompleteAuthFlow ./tests/integration/

# Run resend verification test
go test -v -run TestResendVerificationFlow ./tests/integration/
```

### Run with Timeout
```bash
go test -v ./tests/integration/... -timeout 60s
```

### Run with Coverage
```bash
go test -v -coverprofile=coverage.out ./tests/integration/...
go tool cover -html=coverage.out
```

## Configuration

Integration tests use the **real database** specified in your `.env` file:
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`

Make sure your database is accessible before running tests.

## Data Cleanup

Tests automatically clean up test data after completion:
- Test users are deleted
- Email verifications are removed
- No test data is left in the database

Test emails use timestamp-based unique identifiers:
- Format: `test-{timestamp}@example.com`
- Format: `resend-{timestamp}@example.com`

## Test Structure

Each integration test follows this pattern:

```go
func TestFeature(t *testing.T) {
    // Setup
    loadTestConfig(t)
    db := setupTestDB(t)
    router := setupTestRouter(db)
    
    // Cleanup after test
    defer func() {
        // Delete test data
    }()
    
    // Test steps
    t.Run("Step 1: Description", func(t *testing.T) {
        // Test implementation
    })
}
```

## Expected Results

All tests should pass with:
- ✅ HTTP status codes matching expectations
- ✅ Response structure following API standards
- ✅ Token generation and validation working
- ✅ Database operations successful
- ✅ Error handling correct (409 for duplicates, 401 for unauthorized, etc.)

## Troubleshooting

### Database Connection Failed
```bash
Error: Failed to connect to test database
```
**Solution:** Check your `.env` file and ensure PostgreSQL is running.

### Test Data Cleanup Failed
```bash
Warning: Failed to delete test user
```
**Solution:** Check database permissions and constraints.

### SMTP Errors (If Email Sending is Enabled)
```bash
Error: Failed to send verification OTP
```
**Solution:** Email sending is optional for tests. OTP is captured directly from database.

## Future Tests

Planned integration tests:
- [ ] Password reset flow
- [ ] Event creation and management
- [ ] Event swipe and matching
- [ ] Chat functionality
- [ ] User profile management
- [ ] Google OAuth flow
- [ ] File upload (avatar)

## Notes

- Tests use real database connections
- OTPs are captured from database for testing
- Email sending is bypassed (OTP captured directly)
- All tests are independent and can run in parallel
- Unique email addresses prevent test conflicts

