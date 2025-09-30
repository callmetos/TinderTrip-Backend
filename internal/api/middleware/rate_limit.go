package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"TinderTrip-Backend/internal/utils"
	"TinderTrip-Backend/pkg/config"
	"TinderTrip-Backend/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimiter represents a rate limiter
type RateLimiter struct {
	limiter *limiter.Limiter
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	cfg := config.AppConfig.RateLimit

	// Create rate limiter with memory store
	store := memory.NewStore()

	// Create rate limiter
	rate := limiter.Rate{
		Period: parseDuration(cfg.Window),
		Limit:  int64(cfg.Requests),
	}

	instance := limiter.New(store, rate)

	return &RateLimiter{
		limiter: instance,
	}
}

// RateLimit middleware for rate limiting
func RateLimit() gin.HandlerFunc {
	rateLimiter := NewRateLimiter()

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Create context
		ctx := context.Background()

		// Get rate limit info
		context, err := rateLimiter.limiter.Get(ctx, clientIP)
		if err != nil {
			utils.Logger().WithFields(map[string]interface{}{
				"error":     err,
				"client_ip": clientIP,
			}).Error("Rate limiter error")

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Rate limiter error",
				"message": "Unable to process request",
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		// Check if limit exceeded
		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"message": fmt.Sprintf("You have exceeded the rate limit of %d requests per %s",
					config.AppConfig.RateLimit.Requests,
					config.AppConfig.RateLimit.Window),
				"retry_after": context.Reset - time.Now().Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitWithRedis uses Redis for rate limiting
func RateLimitWithRedis() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Create Redis key
		key := fmt.Sprintf("rate_limit:%s", clientIP)

		// Get current count from Redis
		ctx := context.Background()
		count, err := database.GetRedisClient().Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			utils.Logger().WithFields(map[string]interface{}{
				"error":     err,
				"client_ip": clientIP,
			}).Error("Redis rate limiter error")

			c.Next()
			return
		}

		// Check if limit exceeded
		if count >= config.AppConfig.RateLimit.Requests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"message": fmt.Sprintf("You have exceeded the rate limit of %d requests per %s",
					config.AppConfig.RateLimit.Requests,
					config.AppConfig.RateLimit.Window),
			})
			c.Abort()
			return
		}

		// Increment counter
		windowDuration := parseDuration(config.AppConfig.RateLimit.Window)
		err = database.GetRedisClient().Set(ctx, key, count+1, windowDuration).Err()
		if err != nil {
			utils.Logger().WithFields(map[string]interface{}{
				"error":     err,
				"client_ip": clientIP,
			}).Error("Redis rate limiter set error")
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.AppConfig.RateLimit.Requests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(config.AppConfig.RateLimit.Requests-count-1))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(windowDuration).Unix(), 10))

		c.Next()
	}
}

// RateLimitByUser rate limits by user ID (for authenticated users)
func RateLimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			// If no user ID, use IP-based rate limiting
			RateLimit()(c)
			return
		}

		// Create Redis key for user
		key := fmt.Sprintf("rate_limit:user:%s", userID)

		// Get current count from Redis
		ctx := context.Background()
		count, err := database.GetRedisClient().Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			utils.Logger().WithFields(map[string]interface{}{
				"error":   err,
				"user_id": userID,
			}).Error("Redis user rate limiter error")

			c.Next()
			return
		}

		// Check if limit exceeded
		if count >= config.AppConfig.RateLimit.Requests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"message": fmt.Sprintf("You have exceeded the rate limit of %d requests per %s",
					config.AppConfig.RateLimit.Requests,
					config.AppConfig.RateLimit.Window),
			})
			c.Abort()
			return
		}

		// Increment counter
		windowDuration := parseDuration(config.AppConfig.RateLimit.Window)
		err = database.GetRedisClient().Set(ctx, key, count+1, windowDuration).Err()
		if err != nil {
			utils.Logger().WithFields(map[string]interface{}{
				"error":   err,
				"user_id": userID,
			}).Error("Redis user rate limiter set error")
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.AppConfig.RateLimit.Requests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(config.AppConfig.RateLimit.Requests-count-1))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(windowDuration).Unix(), 10))

		c.Next()
	}
}

// parseDuration parses duration string to time.Duration
func parseDuration(duration string) time.Duration {
	switch duration {
	case "1s":
		return time.Second
	case "1m":
		return time.Minute
	case "1h":
		return time.Hour
	case "1d":
		return 24 * time.Hour
	default:
		return time.Hour // Default to 1 hour
	}
}
