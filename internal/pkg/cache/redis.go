package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache defines the interface for cache operations
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Ping(ctx context.Context) error
	Close() error
}

// RedisConfig defines the configuration for Redis connection
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache client
func NewRedisCache(config RedisConfig) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &redisCache{
		client: client,
	}, nil
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *redisCache) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *redisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *redisCache) Close() error {
	return c.client.Close()
}
