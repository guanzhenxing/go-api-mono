package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitConfig defines rate limit configuration
type RateLimitConfig struct {
	Requests int `yaml:"requests"` // Number of requests per second
	Burst    int `yaml:"burst"`    // Maximum burst size
}

// RateLimit creates a rate limit middleware
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(config.Requests), config.Burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}
		c.Next()
	}
}
