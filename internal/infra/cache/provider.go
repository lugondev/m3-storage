package cache

import (
	"context"
	"fmt"
	"os"

	"github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/infra/config"
)

// InitializeRedisClient sets up the Redis connection based on the configuration.
// Returns the custom RedisClient wrapper.
func InitializeRedisClient(cfg config.Config, log logger.Logger) (*RedisClient, error) { // Changed return type
	ctx := context.Background()                   // Use background context for initialization
	redisClient, err := NewRedisClient(cfg.Redis) // Use existing NewRedisClient function
	if err != nil {
		log.Errorf(ctx, "Failed to connect to Redis: %v", err)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info(ctx, "Redis connection established successfully")
	return redisClient, nil
}

// CloseRedisClient safely closes the Redis client connection.
func CloseRedisClient(client *RedisClient, log logger.Logger) { // Changed parameter type
	ctx := context.Background() // Use background context for shutdown
	if client != nil {
		if err := client.Close(); err != nil { // Call Close() on the wrapper
			log.Errorf(ctx, "Failed to close Redis client gracefully: %v", err)
		} else {
			log.Info(ctx, "Redis client closed successfully.")
		}
	}
}

// ExitOnError logs a fatal error and exits if err is not nil.
// Reusing the pattern for consistency.
func ExitOnError(log logger.Logger, msg string, err error) {
	if err != nil {
		log.Errorf(context.Background(), "%s: %v", msg, err)
		os.Exit(1)
	}
}
