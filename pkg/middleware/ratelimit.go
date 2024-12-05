package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitConfig 定义限流配置
type RateLimitConfig struct {
	Enabled bool    `yaml:"enabled"` // 是否启用限流
	Limit   float64 `yaml:"limit"`   // 每秒允许的请求数
	Burst   int     `yaml:"burst"`   // 允许的突发请求数
}

// RateLimit 创建一个限流中间件
// 它使用令牌桶算法来限制请求速率：
// - limit: 每秒允许的请求数
// - burst: 允许的突发请求数
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	// 如果未启用限流，返回空中间件
	if !config.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 创建限流器
	limiter := rate.NewLimiter(rate.Limit(config.Limit), config.Burst)

	return func(c *gin.Context) {
		// 尝试获取令牌
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			return
		}
		c.Next()
	}
}
