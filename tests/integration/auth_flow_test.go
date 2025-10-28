package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"TinderTrip-Backend/internal/api/handlers"
	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/dto"
	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/config"
	"TinderTrip-Backend/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupTestDB connects to the real PostgreSQL database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	// Connect to existing database
	err := database.ConnectPostgres()
	if err != nil {
		t.Skipf("Skipping integration test: Failed to connect to database: %v", err)
	}

	db := database.GetDB()
	if db == nil {
		t.Skip("Skipping integration test: Database connection is nil")
	}

	return db
}

// setupTestRouter creates a test router with all routes
func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestID())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()

	// Auth routes
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/verify-email", authHandler.VerifyEmail)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/resend-verification", authHandler.ResendVerification)
		authGroup.POST("/forgot-password", authHandler.ForgotPassword)
		authGroup.POST("/reset-password", authHandler.ResetPassword)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.GET("/google", authHandler.GoogleAuth)
		authGroup.GET("/google/callback", authHandler.GoogleCallback)

		// Protected routes
		authProtected := authGroup.Group("")
		authProtected.Use(middleware.AuthMiddleware())
		{
			authProtected.POST("/logout", authHandler.Logout)
			authProtected.GET("/check", authHandler.Check)
		}
	}

	return router
}

// loadTestConfig loads test configuration
func loadTestConfig(t *testing.T) {
	// Try to load configuration from environment
	// If config is not available (e.g., in CI without .env), skip the test
	defer func() {
		if r := recover(); r != nil {
			t.Skip("Skipping integration test: Configuration not available (missing .env file or environment variables)")
		}
	}()

	config.LoadConfig()
}

// TestCompleteAuthFlow tests the complete authentication flow
func TestCompleteAuthFlow(t *testing.T) {
	// Setup
	loadTestConfig(t)
	db := setupTestDB(t)
	router := setupTestRouter(db)

	testEmail := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())
	testPassword := "SecurePass123!"
	testDisplayName := "Test User"
	var capturedOTP string

	// Cleanup test data after test completes
	defer func() {
		// Delete test user and related data
		db.Where("email = ?", testEmail).Delete(&models.User{})
		db.Where("email = ?", testEmail).Delete(&models.EmailVerification{})
	}()

	t.Run("Step 1: Register with email", func(t *testing.T) {
		// Prepare request
		registerReq := dto.RegisterRequest{
			Email:       testEmail,
			Password:    testPassword,
			DisplayName: testDisplayName,
		}
		body, _ := json.Marshal(registerReq)

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code, "Registration should succeed")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Contains(t, response["message"].(string), "Verification OTP sent")

		// Capture OTP from database for testing
		var verification models.EmailVerification
		err = db.Where("email = ?", testEmail).First(&verification).Error
		require.NoError(t, err, "Should find verification record")
		capturedOTP = verification.OTP
		t.Logf("Captured OTP: %s", capturedOTP)
	})

	t.Run("Step 2: Verify OTP with invalid code (should fail)", func(t *testing.T) {
		// Prepare request with wrong OTP
		verifyReq := dto.RegisterWithOTPRequest{
			Email:       testEmail,
			Password:    testPassword,
			DisplayName: testDisplayName,
			OTP:         "000000", // Wrong OTP
		}
		body, _ := json.Marshal(verifyReq)

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/verify-email", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code, "Should fail with wrong OTP")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response["success"].(bool))
	})

	var authToken string

	t.Run("Step 3: Verify OTP with valid code (should succeed)", func(t *testing.T) {
		// Prepare request with correct OTP
		verifyReq := dto.RegisterWithOTPRequest{
			Email:       testEmail,
			Password:    testPassword,
			DisplayName: testDisplayName,
			OTP:         capturedOTP,
		}
		body, _ := json.Marshal(verifyReq)

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/verify-email", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusCreated, w.Code, "Registration should complete successfully")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.NotEmpty(t, response["token"], "Should return authentication token")

		// Extract token and user info
		authToken = response["token"].(string)
		userData := response["user"].(map[string]interface{})
		assert.Equal(t, testEmail, userData["email"])
		assert.Equal(t, testDisplayName, userData["display_name"])
		assert.True(t, userData["is_verified"].(bool), "User should be verified")

		t.Logf("Registration successful. Token: %s", authToken)
	})

	t.Run("Step 4: Attempt duplicate registration (should fail)", func(t *testing.T) {
		// Try to register again with same email
		registerReq := dto.RegisterRequest{
			Email:       testEmail,
			Password:    testPassword,
			DisplayName: testDisplayName,
		}
		body, _ := json.Marshal(registerReq)

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response - should return 409 Conflict
		assert.Equal(t, http.StatusConflict, w.Code, "Should fail with conflict")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["message"].(string), "already exists")
	})

	t.Run("Step 5: Login with wrong password (should fail)", func(t *testing.T) {
		// Prepare request with wrong password
		loginReq := dto.LoginRequest{
			Email:    testEmail,
			Password: "WrongPassword123!",
		}
		body, _ := json.Marshal(loginReq)

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should fail with wrong password")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response["success"].(bool))
	})

	t.Run("Step 6: Login with correct credentials (should succeed)", func(t *testing.T) {
		// Prepare request with correct credentials
		loginReq := dto.LoginRequest{
			Email:    testEmail,
			Password: testPassword,
		}
		body, _ := json.Marshal(loginReq)

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code, "Login should succeed")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.NotEmpty(t, response["token"], "Should return new authentication token")

		// Update token for next tests
		authToken = response["token"].(string)
		userData := response["user"].(map[string]interface{})
		assert.Equal(t, testEmail, userData["email"])
		assert.Equal(t, testDisplayName, userData["display_name"])

		t.Logf("Login successful. New token: %s", authToken)
	})

	t.Run("Step 7: Access protected endpoint with valid token", func(t *testing.T) {
		// Make request to protected endpoint
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/auth/check", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code, "Should access protected endpoint")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))

		// Check data structure
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "valid", data["status"])
		assert.Equal(t, testEmail, data["email"])
	})

	t.Run("Step 8: Access protected endpoint without token (should fail)", func(t *testing.T) {
		// Make request without token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/auth/check", nil)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should fail without token")
	})

	t.Run("Step 9: Access protected endpoint with invalid token (should fail)", func(t *testing.T) {
		// Make request with invalid token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/auth/check", nil)
		req.Header.Set("Authorization", "Bearer invalid-token-12345")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should fail with invalid token")
	})

	t.Run("Step 10: Logout", func(t *testing.T) {
		// Make request to logout
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code, "Logout should succeed")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Contains(t, response["message"].(string), "Logged out")
	})
}

// TestResendVerificationFlow tests the resend verification OTP flow
func TestResendVerificationFlow(t *testing.T) {
	// Setup
	loadTestConfig(t)
	db := setupTestDB(t)
	router := setupTestRouter(db)

	testEmail := fmt.Sprintf("resend-%d@example.com", time.Now().UnixNano())
	var firstOTP, secondOTP string

	// Cleanup test data after test completes
	defer func() {
		// Delete test user and related data
		db.Where("email = ?", testEmail).Delete(&models.User{})
		db.Where("email = ?", testEmail).Delete(&models.EmailVerification{})
	}()

	t.Run("Step 1: Initial registration", func(t *testing.T) {
		registerReq := dto.RegisterRequest{
			Email:       testEmail,
			Password:    "SecurePass123!",
			DisplayName: "Resend Test User",
		}
		body, _ := json.Marshal(registerReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Capture first OTP
		var verification models.EmailVerification
		err := db.Where("email = ?", testEmail).First(&verification).Error
		require.NoError(t, err)
		firstOTP = verification.OTP
		t.Logf("First OTP: %s", firstOTP)
	})

	t.Run("Step 2: Resend verification OTP", func(t *testing.T) {
		// Wait a bit to ensure different OTP
		time.Sleep(100 * time.Millisecond)

		resendReq := dto.ResendVerificationRequest{
			Email: testEmail,
		}
		body, _ := json.Marshal(resendReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/resend-verification", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response["success"].(bool))

		// Capture second OTP
		var verification models.EmailVerification
		err = db.Where("email = ?", testEmail).First(&verification).Error
		require.NoError(t, err)
		secondOTP = verification.OTP
		t.Logf("Second OTP: %s", secondOTP)

		// OTPs might be the same due to random generation, but the record should be updated
		assert.NotEmpty(t, secondOTP)
	})

	t.Run("Step 3: Verify with new OTP", func(t *testing.T) {
		verifyReq := dto.RegisterWithOTPRequest{
			Email:       testEmail,
			Password:    "SecurePass123!",
			DisplayName: "Resend Test User",
			OTP:         secondOTP,
		}
		body, _ := json.Marshal(verifyReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/verify-email", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Should verify successfully with new OTP")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response["success"].(bool))
		assert.NotEmpty(t, response["token"])
	})
}
