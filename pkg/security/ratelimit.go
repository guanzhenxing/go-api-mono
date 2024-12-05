package security

import (
	"sync"
	"time"
)

// TokenBucket 实现令牌桶算法
type TokenBucket struct {
	rate       float64    // 令牌产生速率
	capacity   float64    // 桶的容量
	tokens     float64    // 当前令牌数量
	lastUpdate time.Time  // 上次更新时间
	mu         sync.Mutex // 互斥锁
}

// NewTokenBucket 创建一个新的令牌桶
func NewTokenBucket(rate float64, capacity float64) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastUpdate: time.Now(),
	}
}

// Allow 检查是否允许请求通过
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastUpdate).Seconds()
	tb.tokens = min(tb.capacity, tb.tokens+elapsed*tb.rate)
	tb.lastUpdate = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow() bool
}

// IPRateLimiter IP限流器
type IPRateLimiter struct {
	limiters sync.Map
	rate     float64
	capacity float64
}

// NewIPRateLimiter 创建一个新的IP限流器
func NewIPRateLimiter(rate float64, capacity float64) *IPRateLimiter {
	return &IPRateLimiter{
		rate:     rate,
		capacity: capacity,
	}
}

// GetLimiter 获取指定IP的限流器
func (rl *IPRateLimiter) GetLimiter(ip string) RateLimiter {
	limiter, exists := rl.limiters.Load(ip)
	if !exists {
		limiter = NewTokenBucket(rl.rate, rl.capacity)
		rl.limiters.Store(ip, limiter)
	}
	rateLimiter, ok := limiter.(RateLimiter)
	if !ok {
		return nil
	}
	return rateLimiter
}

// min returns the smaller of x or y.
func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}
