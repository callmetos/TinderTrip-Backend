package middleware

import (
	"TinderTrip-Backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID exists in header (X-Request-ID)
		requestID := c.GetHeader("X-Request-ID")

		// If not provided, generate a new UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context for use in handlers and logging
		c.Set(utils.RequestIDKey, requestID)

		// Set response header
		c.Header("X-Request-ID", requestID)

		// Continue to next handler
		c.Next()
	}
}
