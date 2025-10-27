package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"TinderTrip-Backend/internal/api/middleware"
	"TinderTrip-Backend/internal/utils"
	"TinderTrip-Backend/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupMiddlewareTests() {
	gin.SetMode(gin.TestMode)
	config.AppConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret-key-for-middleware-testing",
			ExpireHours: 24,
		},
	}
}

func TestAuthMiddleware(t *testing.T) {
	setupMiddlewareTests()

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
		checkContext   func(t *testing.T, c *gin.Context)
		expectAbort    bool
	}{
		{
			name: "Valid token",
			setupRequest: func() *http.Request {
				token, _ := utils.GenerateToken("user-123", "test@example.com", "password")
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				return req
			},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				userID, exists := c.Get("user_id")
				assert.True(t, exists)
				assert.Equal(t, "user-123", userID)

				email, exists := c.Get("user_email")
				assert.True(t, exists)
				assert.Equal(t, "test@example.com", email)

				provider, exists := c.Get("user_provider")
				assert.True(t, exists)
				assert.Equal(t, "password", provider)
			},
			expectAbort: false,
		},
		{
			name: "Missing Authorization header",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name: "Invalid Authorization header format",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "InvalidFormat token123")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name: "Token without Bearer prefix",
			setupRequest: func() *http.Request {
				token, _ := utils.GenerateToken("user-123", "test@example.com", "password")
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", token)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name: "Invalid token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer invalid.token.here")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name: "Empty token after Bearer",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer ")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, router := gin.CreateTestContext(w)
			c.Request = tt.setupRequest()

			// Set up middleware and handler
			handlerCalled := false
			router.Use(middleware.AuthMiddleware())
			router.GET("/test", func(c *gin.Context) {
				handlerCalled = true
				if tt.checkContext != nil {
					tt.checkContext(t, c)
				}
				c.Status(http.StatusOK)
			})

			router.ServeHTTP(w, c.Request)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, !tt.expectAbort, handlerCalled)
		})
	}
}

func TestOptionalAuthMiddleware(t *testing.T) {
	setupMiddlewareTests()

	tests := []struct {
		name         string
		setupRequest func() *http.Request
		checkContext func(t *testing.T, c *gin.Context)
	}{
		{
			name: "Valid token",
			setupRequest: func() *http.Request {
				token, _ := utils.GenerateToken("user-123", "test@example.com", "password")
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				return req
			},
			checkContext: func(t *testing.T, c *gin.Context) {
				userID, exists := c.Get("user_id")
				assert.True(t, exists)
				assert.Equal(t, "user-123", userID)
			},
		},
		{
			name: "No token - should continue",
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/test", nil)
			},
			checkContext: func(t *testing.T, c *gin.Context) {
				_, exists := c.Get("user_id")
				assert.False(t, exists)
			},
		},
		{
			name: "Invalid token - should continue without setting context",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer invalid.token")
				return req
			},
			checkContext: func(t *testing.T, c *gin.Context) {
				_, exists := c.Get("user_id")
				assert.False(t, exists)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, router := gin.CreateTestContext(w)
			c.Request = tt.setupRequest()

			handlerCalled := false
			router.Use(middleware.OptionalAuthMiddleware())
			router.GET("/test", func(c *gin.Context) {
				handlerCalled = true
				tt.checkContext(t, c)
				c.Status(http.StatusOK)
			})

			router.ServeHTTP(w, c.Request)

			assert.True(t, handlerCalled, "Handler should always be called")
		})
	}
}

func TestGetCurrentUserID(t *testing.T) {
	tests := []struct {
		name      string
		setupCtx  func() *gin.Context
		wantID    string
		wantExist bool
	}{
		{
			name: "User ID exists in context",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_id", "user-123")
				return c
			},
			wantID:    "user-123",
			wantExist: true,
		},
		{
			name: "User ID does not exist",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				return c
			},
			wantID:    "",
			wantExist: false,
		},
		{
			name: "User ID is wrong type",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_id", 12345) // int instead of string
				return c
			},
			wantID:    "",
			wantExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupCtx()
			id, exists := middleware.GetCurrentUserID(c)
			assert.Equal(t, tt.wantExist, exists)
			assert.Equal(t, tt.wantID, id)
		})
	}
}

func TestGetCurrentUserEmail(t *testing.T) {
	tests := []struct {
		name      string
		setupCtx  func() *gin.Context
		wantEmail string
		wantExist bool
	}{
		{
			name: "Email exists in context",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_email", "test@example.com")
				return c
			},
			wantEmail: "test@example.com",
			wantExist: true,
		},
		{
			name: "Email does not exist",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				return c
			},
			wantEmail: "",
			wantExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupCtx()
			email, exists := middleware.GetCurrentUserEmail(c)
			assert.Equal(t, tt.wantExist, exists)
			assert.Equal(t, tt.wantEmail, email)
		})
	}
}

func TestGetCurrentUserProvider(t *testing.T) {
	tests := []struct {
		name         string
		setupCtx     func() *gin.Context
		wantProvider string
		wantExist    bool
	}{
		{
			name: "Provider exists in context",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_provider", "google")
				return c
			},
			wantProvider: "google",
			wantExist:    true,
		},
		{
			name: "Provider does not exist",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				return c
			},
			wantProvider: "",
			wantExist:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupCtx()
			provider, exists := middleware.GetCurrentUserProvider(c)
			assert.Equal(t, tt.wantExist, exists)
			assert.Equal(t, tt.wantProvider, provider)
		})
	}
}

func TestRequireAuth(t *testing.T) {
	tests := []struct {
		name     string
		setupCtx func() *gin.Context
		want     bool
	}{
		{
			name: "User authenticated",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_id", "user-123")
				return c
			},
			want: true,
		},
		{
			name: "User not authenticated",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				return c
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupCtx()
			result := middleware.RequireAuth(c)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestRequireProvider(t *testing.T) {
	tests := []struct {
		name     string
		setupCtx func() *gin.Context
		provider string
		want     bool
	}{
		{
			name: "Correct provider",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_provider", "password")
				return c
			},
			provider: "password",
			want:     true,
		},
		{
			name: "Wrong provider",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_provider", "google")
				return c
			},
			provider: "password",
			want:     false,
		},
		{
			name: "Provider not set",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				return c
			},
			provider: "password",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupCtx()
			result := middleware.RequireProvider(c, tt.provider)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestRequireEmailAuth(t *testing.T) {
	tests := []struct {
		name     string
		setupCtx func() *gin.Context
		want     bool
	}{
		{
			name: "Email/password auth",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_provider", "password")
				return c
			},
			want: true,
		},
		{
			name: "Google auth",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_provider", "google")
				return c
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupCtx()
			result := middleware.RequireEmailAuth(c)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestRequireGoogleAuth(t *testing.T) {
	tests := []struct {
		name     string
		setupCtx func() *gin.Context
		want     bool
	}{
		{
			name: "Google auth",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_provider", "google")
				return c
			},
			want: true,
		},
		{
			name: "Email/password auth",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_provider", "password")
				return c
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setupCtx()
			result := middleware.RequireGoogleAuth(c)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestAdminMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupCtx       func() *gin.Context
		expectedStatus int
		expectAbort    bool
	}{
		{
			name: "Authenticated user",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_id", "user-123")
				return c
			},
			expectedStatus: http.StatusOK,
			expectAbort:    false,
		},
		{
			name: "Not authenticated",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				return c
			},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name: "Empty user ID",
			setupCtx: func() *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Set("user_id", "")
				return c
			},
			expectedStatus: http.StatusForbidden,
			expectAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, router := gin.CreateTestContext(w)

			// Setup middleware that sets user_id in context before AdminMiddleware
			router.Use(func(c *gin.Context) {
				setupCtx := tt.setupCtx()
				if userID, exists := setupCtx.Get("user_id"); exists {
					c.Set("user_id", userID)
				}
				c.Next()
			})

			c.Request = httptest.NewRequest("GET", "/test", nil)

			handlerCalled := false
			router.Use(middleware.AdminMiddleware())
			router.GET("/test", func(c *gin.Context) {
				handlerCalled = true
				c.Status(http.StatusOK)
			})

			router.ServeHTTP(w, c.Request)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, !tt.expectAbort, handlerCalled)
		})
	}
}
