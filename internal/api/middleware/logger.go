package middleware

import (
	"time"

	"TinderTrip-Backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request
		utils.Logger().WithFields(map[string]interface{}{
			"timestamp":  start.Format(time.RFC3339),
			"status":     c.Writer.Status(),
			"latency":    latency,
			"client_ip":  c.ClientIP(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"user_agent": c.Request.UserAgent(),
			"error":      c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}).Info("HTTP Request")
	}
}
