package cache

import (
	"context"
	"time"

	"github.com/lugondev/m3-storage/internal/infra/config"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the redis client
type RedisClient struct {
	client *redis.Client
	opt    *redis.Options
}

// NewRedisClient creates a new Redis client connection
func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	var client *redis.Client
	if cfg.Url != "" {
		opt, err := redis.ParseURL(cfg.Url)
		if err != nil {
			return nil, err
		}
		client = redis.NewClient(opt)
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Host + ":" + cfg.Port,
			Username: cfg.User,
			Password: cfg.Pass,
			DB:       cfg.DB,
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{client: client, opt: client.Options()}, nil
}

// Get retrieves a value from Redis
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set stores a value in Redis
func (r *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Delete removes a value from Redis
func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Client returns the underlying go-redis client.
// Added to allow access for operations like Ping in health checks.
func (r *RedisClient) Client() *redis.Client {
	return r.client
}

// GetRedisURL returns the Redis URL.
func (r *RedisClient) GetRedisURL() string {
	return r.client.Options().Addr
}

// GetRedisOptions returns the Redis options.
func (r *RedisClient) GetRedisOptions() redis.Options {
	return *(r.opt)
}
