package middleware

import (
	"strings"

	"TinderTrip-Backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthContext middleware to extract user ID from JWT token and add to context
func AuthContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token and extract claims
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Add user ID to context
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
