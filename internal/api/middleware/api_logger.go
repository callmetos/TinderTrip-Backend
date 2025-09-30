package middleware

import (
	"fmt"
	"time"

	"TinderTrip-Backend/internal/models"
	"TinderTrip-Backend/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// APILogger middleware for logging API requests to database
func APILogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("APILogger middleware called for %s %s\n", c.Request.Method, c.Request.URL.Path)

		// Start timer
		start := time.Now()

		// Generate request ID
		requestID := uuid.New().String()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		durationMs := int(duration.Milliseconds())

		// Get user ID from context if available
		var userID *uuid.UUID
		if userIDStr, exists := c.Get("user_id"); exists {
			if id, ok := userIDStr.(string); ok {
				if parsedID, err := uuid.Parse(id); err == nil {
					userID = &parsedID
				}
			}
		}

		// Create API log entry
		method := c.Request.Method
		path := c.Request.URL.Path
		status := c.Writer.Status()
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()

		apiLog := &models.APILog{
			RequestID:  &requestID,
			UserID:     userID,
			Method:     &method,
			Path:       &path,
			Status:     &status,
			DurationMs: &durationMs,
			IPAddress:  &ipAddress,
			UserAgent:  &userAgent,
		}

		// Save to database synchronously for debugging
		if err := database.GetDB().Create(apiLog).Error; err != nil {
			// Log error but don't fail the request
			fmt.Printf("Error saving API log: %v\n", err)
		} else {
			fmt.Printf("API log saved successfully: %s %s %d\n", *apiLog.Method, *apiLog.Path, *apiLog.Status)
		}

		// Add request ID to response headers
		c.Header("X-Request-ID", requestID)
	}
}
