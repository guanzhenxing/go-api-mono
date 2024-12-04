package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache 定义了缓存操作的接口
// 这个接口抽象了基本的缓存操作，使得可以轻松替换不同的缓存实现
type Cache interface {
	// Get 获取缓存的值
	Get(ctx context.Context, key string) (string, error)
	// Set 设置缓存的值，可以指定过期时间
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Del 删除缓存的键
	Del(ctx context.Context, key string) error
	// Ping 检查缓存服务是否可用
	Ping(ctx context.Context) error
	// Close 关闭缓存连接
	Close() error
}

// RedisConfig 定义了Redis连接的配置参数
type RedisConfig struct {
	Host     string // Redis服务器地址
	Port     int    // Redis服务器端口
	Password string // Redis认证密码
	DB       int    // Redis数据库索引
}

// redisCache 是Cache接口的Redis实现
type redisCache struct {
	client *redis.Client // Redis客户端实例
}

// NewRedisCache 创建一个新的Redis缓存客户端
// 它会立即尝试连接Redis服务器，如果连接失败则返回错误
func NewRedisCache(config RedisConfig) (Cache, error) {
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// 创建一个带超时的上下文来测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试连接是否成功
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &redisCache{
		client: client,
	}, nil
}

// Get 实现了Cache接口的Get方法
func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set 实现了Cache接口的Set方法
func (c *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Del 实现了Cache接口的Del方法
func (c *redisCache) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Ping 实现了Cache接口的Ping方法
func (c *redisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close 实现了Cache接口的Close方法
func (c *redisCache) Close() error {
	return c.client.Close()
}
