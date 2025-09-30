package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"TinderTrip-Backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// Recovery middleware for handling panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// Log the error
				utils.Logger().WithFields(map[string]interface{}{
					"error":   err,
					"stack":   string(debug.Stack()),
					"request": c.Request.URL.Path,
					"method":  c.Request.Method,
				}).Error("Panic recovered")

				// If it's a broken pipe, don't send a response
				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
					return
				}

				// Send error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":      "Internal server error",
					"message":    "Something went wrong. Please try again later.",
					"request_id": c.GetString("request_id"),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RecoveryWithWriter middleware with custom writer
func RecoveryWithWriter() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// Get request dump
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// Log the error with more details
				utils.Logger().WithFields(map[string]interface{}{
					"error":      err,
					"stack":      string(debug.Stack()),
					"request":    string(httpRequest),
					"user_agent": c.Request.UserAgent(),
					"client_ip":  c.ClientIP(),
					"request_id": c.GetString("request_id"),
				}).Error("Panic recovered with detailed logging")

				// If it's a broken pipe, don't send a response
				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
					return
				}

				// Send error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":      "Internal server error",
					"message":    "Something went wrong. Please try again later.",
					"request_id": c.GetString("request_id"),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CustomRecovery middleware with custom error handling
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// Log the error
				utils.Logger().WithFields(map[string]interface{}{
					"error":      err,
					"stack":      string(debug.Stack()),
					"request":    c.Request.URL.Path,
					"method":     c.Request.Method,
					"user_agent": c.Request.UserAgent(),
					"client_ip":  c.ClientIP(),
					"request_id": c.GetString("request_id"),
				}).Error("Panic recovered with custom handler")

				// If it's a broken pipe, don't send a response
				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
					return
				}

				// Determine error type and send appropriate response
				switch e := err.(type) {
				case string:
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":      "Internal server error",
						"message":    e,
						"request_id": c.GetString("request_id"),
					})
				case error:
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":      "Internal server error",
						"message":    e.Error(),
						"request_id": c.GetString("request_id"),
					})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":      "Internal server error",
						"message":    fmt.Sprintf("%v", err),
						"request_id": c.GetString("request_id"),
					})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
