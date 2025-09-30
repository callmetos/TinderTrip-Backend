package middleware

import (
	"net/http"

	"TinderTrip-Backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header is required",
				"message": "Please provide a valid token",
			})
			c.Abort()
			return
		}

		// Extract token from header
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization header",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_provider", claims.Provider)

		c.Next()
	}
}

// OptionalAuthMiddleware handles optional JWT authentication
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Extract token from header
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user information in context if valid
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_provider", claims.Provider)

		c.Next()
	}
}

// AdminMiddleware checks if user is admin (you can implement admin logic)
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authentication required",
				"message": "Please login to access this resource",
			})
			c.Abort()
			return
		}

		// TODO: Implement admin check logic
		// For now, we'll just check if user exists
		if userID == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"message": "Admin privileges required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUserID gets the current user ID from context
func GetCurrentUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", false
	}

	return userIDStr, true
}

// GetCurrentUserEmail gets the current user email from context
func GetCurrentUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", false
	}

	return emailStr, true
}

// GetCurrentUserProvider gets the current user provider from context
func GetCurrentUserProvider(c *gin.Context) (string, bool) {
	provider, exists := c.Get("user_provider")
	if !exists {
		return "", false
	}

	providerStr, ok := provider.(string)
	if !ok {
		return "", false
	}

	return providerStr, true
}

// RequireAuth is a helper function to check if user is authenticated
func RequireAuth(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// RequireProvider checks if user has specific provider
func RequireProvider(c *gin.Context, provider string) bool {
	userProvider, exists := GetCurrentUserProvider(c)
	if !exists {
		return false
	}

	return userProvider == provider
}

// RequireEmailAuth checks if user uses email authentication
func RequireEmailAuth(c *gin.Context) bool {
	return RequireProvider(c, "password")
}

// RequireGoogleAuth checks if user uses Google authentication
func RequireGoogleAuth(c *gin.Context) bool {
	return RequireProvider(c, "google")
}
