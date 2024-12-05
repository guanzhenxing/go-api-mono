package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisCache(t *testing.T) {
	config := RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	cache, err := NewRedisCache(config)
	if err != nil {
		t.Skip("Redis is not available:", err)
		return
	}
	defer cache.Close()

	ctx := context.Background()

	t.Run("Test Set and Get", func(t *testing.T) {
		key := "test_key"
		value := "test_value"
		expiration := time.Minute

		err := cache.Set(ctx, key, value, expiration)
		assert.NoError(t, err)

		result, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Test Delete", func(t *testing.T) {
		key := "test_delete_key"
		value := "test_value"
		expiration := time.Minute

		err := cache.Set(ctx, key, value, expiration)
		assert.NoError(t, err)

		err = cache.Del(ctx, key)
		assert.NoError(t, err)

		_, err = cache.Get(ctx, key)
		assert.Error(t, err) // Should return redis.Nil
	})

	t.Run("Test Ping", func(t *testing.T) {
		err := cache.Ping(ctx)
		assert.NoError(t, err)
	})
}

func TestNewRedisCache_ConnectionError(t *testing.T) {
	config := RedisConfig{
		Host:     "nonexistent",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	_, err := NewRedisCache(config)
	assert.Error(t, err)
}
