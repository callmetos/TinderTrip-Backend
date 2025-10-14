package middleware

import (
	"strconv"
	"time"

	"TinderTrip-Backend/internal/service"
	"TinderTrip-Backend/internal/utils"

	"github.com/gin-gonic/gin"
)

var monitoringService *service.MonitoringService

// SetMonitoringService sets the monitoring service instance
func SetMonitoringService(ms *service.MonitoringService) {
	monitoringService = ms
}

// PrometheusMetrics middleware for collecting Prometheus metrics
func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		if monitoringService == nil {
			c.Next()
			return
		}

		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get endpoint (remove query parameters and IDs)
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		// Record metrics
		statusCode := strconv.Itoa(c.Writer.Status())
		monitoringService.RecordHTTPRequest(
			c.Request.Method,
			endpoint,
			statusCode,
			duration,
		)

		// Log slow requests
		if duration > 5.0 {
			utils.Logger().WithFields(map[string]interface{}{
				"method":     c.Request.Method,
				"endpoint":   endpoint,
				"duration":   duration,
				"status":     c.Writer.Status(),
				"client_ip":  c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			}).Warn("Slow request detected")
		}
	}
}

// BusinessMetrics middleware for tracking business-specific metrics
func BusinessMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		if monitoringService == nil {
			return
		}

		// Track specific business events
		switch {
		case c.Request.Method == "POST" && c.FullPath() == "/api/v1/auth/register":
			monitoringService.RecordUserRegistration()
		case c.Request.Method == "POST" && c.FullPath() == "/api/v1/events":
			monitoringService.RecordEventCreation()
		}
	}
}
