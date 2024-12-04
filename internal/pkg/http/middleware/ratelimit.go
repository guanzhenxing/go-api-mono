package middleware

import (
	"net/http"
	"strings"

	"go-api-mono/internal/pkg/errors"
	"go-api-mono/internal/pkg/security"
)

// RateLimitOptions 速率限制中间件选项
type RateLimitOptions struct {
	SkipPaths []string                   // 跳过限流的路径
	GetIPKey  func(*http.Request) string // 自定义获取IP的函数
}

// DefaultRateLimitOptions 默认速率限制选项
var DefaultRateLimitOptions = RateLimitOptions{
	SkipPaths: []string{"/health", "/metrics"},
	GetIPKey: func(r *http.Request) string {
		// 默认从X-Real-IP或X-Forwarded-For获取IP
		if ip := r.Header.Get("X-Real-IP"); ip != "" {
			return ip
		}
		if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
			i := strings.Index(ip, ",")
			if i == -1 {
				return ip
			}
			return ip[:i]
		}
		// 从RemoteAddr获取
		ip := r.RemoteAddr
		i := strings.LastIndex(ip, ":")
		if i == -1 {
			return ip
		}
		return ip[:i]
	},
}

// RateLimit 创建速率限制中间件
func RateLimit(limiter *security.IPRateLimiter, opts RateLimitOptions) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 检查是否需要跳过限流
			for _, path := range opts.SkipPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next(w, r)
					return
				}
			}

			// 获取客户端标识（默认使用IP）
			key := opts.GetIPKey(r)

			// 获取限流器并检查是否允许请求
			if !limiter.GetLimiter(key).Allow() {
				http.Error(w, errors.ErrTooManyRequests.Error(), http.StatusTooManyRequests)
				return
			}

			next(w, r)
		}
	}
}

// IPRateLimit 创建基于IP的速率限制中间件
func IPRateLimit(requests, burst int) Middleware {
	limiter := security.NewIPRateLimiter(float64(requests), float64(burst))
	opts := DefaultRateLimitOptions

	return RateLimit(limiter, opts)
}
