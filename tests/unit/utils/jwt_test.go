package utils_test

import (
	"testing"
	"time"

	"TinderTrip-Backend/internal/utils"
	"TinderTrip-Backend/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupJWTTests() {
	// Initialize config for testing
	config.AppConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret-key-for-testing-only",
			ExpireHours: 24,
		},
	}
}

func TestGenerateToken(t *testing.T) {
	setupJWTTests()

	tests := []struct {
		name     string
		userID   string
		email    string
		provider string
		wantErr  bool
	}{
		{
			name:     "Valid token generation",
			userID:   "550e8400-e29b-41d4-a716-446655440000",
			email:    "test@example.com",
			provider: "password",
			wantErr:  false,
		},
		{
			name:     "Generate token with Google provider",
			userID:   "550e8400-e29b-41d4-a716-446655440001",
			email:    "google@example.com",
			provider: "google",
			wantErr:  false,
		},
		{
			name:     "Generate token with empty user ID",
			userID:   "",
			email:    "test@example.com",
			provider: "password",
			wantErr:  false, // Should still work, just not recommended
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.GenerateToken(tt.userID, tt.email, tt.provider)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	setupJWTTests()

	tests := []struct {
		name        string
		setupToken  func() string
		wantErr     bool
		errContains string
	}{
		{
			name: "Valid token",
			setupToken: func() string {
				token, _ := utils.GenerateToken("user-123", "test@example.com", "password")
				return token
			},
			wantErr: false,
		},
		{
			name: "Invalid token format",
			setupToken: func() string {
				return "invalid.token.format"
			},
			wantErr:     true,
			errContains: "failed to parse token",
		},
		{
			name: "Empty token",
			setupToken: func() string {
				return ""
			},
			wantErr:     true,
			errContains: "failed to parse token",
		},
		{
			name: "Token with wrong signature",
			setupToken: func() string {
				// Generate token with different secret
				originalSecret := config.AppConfig.JWT.Secret
				config.AppConfig.JWT.Secret = "different-secret"
				token, _ := utils.GenerateToken("user-123", "test@example.com", "password")
				config.AppConfig.JWT.Secret = originalSecret
				return token
			},
			wantErr:     true,
			errContains: "failed to parse token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			claims, err := utils.ValidateToken(token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, "user-123", claims.UserID)
				assert.Equal(t, "test@example.com", claims.Email)
				assert.Equal(t, "password", claims.Provider)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	setupJWTTests()

	tests := []struct {
		name        string
		setupToken  func() string
		wantErr     bool
		errContains string
	}{
		{
			name: "Refresh valid token",
			setupToken: func() string {
				token, _ := utils.GenerateToken("user-123", "test@example.com", "password")
				return token
			},
			wantErr: false,
		},
		{
			name: "Refresh invalid token",
			setupToken: func() string {
				return "invalid.token"
			},
			wantErr:     true,
			errContains: "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldToken := tt.setupToken()
			newToken, err := utils.RefreshToken(oldToken)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, newToken)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newToken)
				// Note: New token might be the same if generated in the same second
				// The important thing is that it's valid

				// Verify new token is valid
				claims, err := utils.ValidateToken(newToken)
				assert.NoError(t, err)
				assert.Equal(t, "user-123", claims.UserID)
				assert.Equal(t, "test@example.com", claims.Email)
				assert.Equal(t, "password", claims.Provider)
			}
		})
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		wantToken   string
		wantErr     bool
		errContains string
	}{
		{
			name:       "Valid Bearer token",
			authHeader: "Bearer abc123xyz",
			wantToken:  "abc123xyz",
			wantErr:    false,
		},
		{
			name:        "Missing Bearer prefix",
			authHeader:  "abc123xyz",
			wantErr:     true,
			errContains: "must start with 'Bearer '",
		},
		{
			name:        "Empty header",
			authHeader:  "",
			wantErr:     true,
			errContains: "authorization header is required",
		},
		{
			name:        "Only Bearer without token",
			authHeader:  "Bearer",
			wantErr:     true,
			errContains: "must start with 'Bearer '",
		},
		{
			name:       "Bearer with space only",
			authHeader: "Bearer ",
			wantToken:  "",
			wantErr:    false,
		},
		{
			name:       "Bearer with long token",
			authHeader: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			wantToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.ExtractTokenFromHeader(tt.authHeader)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}
		})
	}
}

func TestGetTokenExpiration(t *testing.T) {
	setupJWTTests()

	tests := []struct {
		name        string
		setupToken  func() string
		wantErr     bool
		checkExpiry func(t *testing.T, expiry time.Time)
	}{
		{
			name: "Get expiration from valid token",
			setupToken: func() string {
				token, _ := utils.GenerateToken("user-123", "test@example.com", "password")
				return token
			},
			wantErr: false,
			checkExpiry: func(t *testing.T, expiry time.Time) {
				expected := time.Now().Add(24 * time.Hour)
				// Allow 1 minute difference for test execution time
				assert.WithinDuration(t, expected, expiry, time.Minute)
			},
		},
		{
			name: "Get expiration from invalid token",
			setupToken: func() string {
				return "invalid.token"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			expiry, err := utils.GetTokenExpiration(token)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkExpiry != nil {
					tt.checkExpiry(t, expiry)
				}
			}
		})
	}
}

func TestIsTokenExpired(t *testing.T) {
	setupJWTTests()

	tests := []struct {
		name        string
		setupToken  func() string
		wantExpired bool
	}{
		{
			name: "Non-expired token",
			setupToken: func() string {
				token, _ := utils.GenerateToken("user-123", "test@example.com", "password")
				return token
			},
			wantExpired: false,
		},
		{
			name: "Invalid token returns expired",
			setupToken: func() string {
				return "invalid.token"
			},
			wantExpired: true,
		},
		{
			name: "Empty token returns expired",
			setupToken: func() string {
				return ""
			},
			wantExpired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			expired := utils.IsTokenExpired(token)
			assert.Equal(t, tt.wantExpired, expired)
		})
	}
}

func TestJWTClaimsIntegrity(t *testing.T) {
	setupJWTTests()

	t.Run("Token contains all expected claims", func(t *testing.T) {
		userID := "user-123"
		email := "test@example.com"
		provider := "password"

		token, err := utils.GenerateToken(userID, email, provider)
		require.NoError(t, err)

		claims, err := utils.ValidateToken(token)
		require.NoError(t, err)

		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, provider, claims.Provider)
		assert.Equal(t, "TinderTrip-Backend", claims.Issuer)
		assert.Equal(t, userID, claims.Subject)
		assert.NotZero(t, claims.ExpiresAt)
		assert.True(t, time.Unix(claims.ExpiresAt, 0).After(time.Now()))
	})
}

func TestTokenExpiration(t *testing.T) {
	setupJWTTests()

	t.Run("Token expires after configured hours", func(t *testing.T) {
		config.AppConfig.JWT.ExpireHours = 1 // 1 hour

		token, err := utils.GenerateToken("user-123", "test@example.com", "password")
		require.NoError(t, err)

		expiry, err := utils.GetTokenExpiration(token)
		require.NoError(t, err)

		expected := time.Now().Add(1 * time.Hour)
		assert.WithinDuration(t, expected, expiry, time.Minute)

		// Reset to default
		config.AppConfig.JWT.ExpireHours = 24
	})
}
