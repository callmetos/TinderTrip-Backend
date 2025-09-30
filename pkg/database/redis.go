package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"TinderTrip-Backend/pkg/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() error {
	cfg := config.AppConfig.Redis

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return nil
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

// Cache helper functions
func SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func GetCache(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

func DeleteCache(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

func SetCacheWithJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func GetCacheWithJSON(ctx context.Context, key string, dest interface{}) error {
	return RedisClient.Get(ctx, key).Scan(dest)
}

// Session helper functions
func SetSession(ctx context.Context, sessionID string, userID uint, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return RedisClient.Set(ctx, key, strconv.Itoa(int(userID)), expiration).Err()
}

func GetSession(ctx context.Context, sessionID string) (uint, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	userIDStr, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(userID), nil
}

func DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return RedisClient.Del(ctx, key).Err()
}
