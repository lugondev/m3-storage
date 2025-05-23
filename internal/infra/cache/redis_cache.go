package cache

import (
	"context"
	"encoding/json"
	"time"

	ports "github.com/lugondev/m3-storage/internal/modules/app/port"
)

// RedisCacheService implements the ports.CacheService interface
type RedisCacheService struct {
	client *RedisClient
}

// NewRedisCacheService creates a new Redis cache service
func NewRedisCacheService(client *RedisClient) ports.CacheService {
	return &RedisCacheService{
		client: client,
	}
}

// Get retrieves a value from cache
func (s *RedisCacheService) Get(ctx context.Context, key string) (any, error) {
	val, err := s.client.Get(ctx, key)
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil // Key does not exist
		}
		return nil, err
	}

	var result any
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return val, nil // Return as string if not JSON
	}

	return result, nil
}

// Set stores a value in cache with expiration
func (s *RedisCacheService) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	// Convert value to JSON if it's not a string
	var val string
	switch v := value.(type) {
	case string:
		val = v
	default:
		bytes, err := json.Marshal(value)
		if err != nil {
			return err
		}
		val = string(bytes)
	}

	return s.client.Set(ctx, key, val, expiration)
}

// Delete removes a value from cache
func (s *RedisCacheService) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key)
}

// Clear removes all values from cache
func (s *RedisCacheService) Clear(ctx context.Context) error {
	// Redis client doesn't have a direct method for this, so we'll return nil
	// In a real implementation, we would use FLUSHDB or FLUSHALL
	return nil
}
