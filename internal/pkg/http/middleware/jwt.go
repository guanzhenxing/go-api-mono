package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-api-mono/internal/pkg/auth"
	"go-api-mono/internal/pkg/errors"
)

// JWTOptions JWT中间件选项
type JWTOptions struct {
	TokenPrefix    string   // 令牌前缀（如 "Bearer"）
	SkipPaths      []string // 跳过验证的路径
	ClaimsKey      string   // Claims在上下文中的键
	ContextUserKey string   // 用户信息在上下文中的键
}

// DefaultJWTOptions 默认JWT选项
var DefaultJWTOptions = JWTOptions{
	TokenPrefix:    "Bearer",
	SkipPaths:      []string{"/api/v1/auth/login", "/api/v1/auth/register"},
	ClaimsKey:      string(ClaimsKey),
	ContextUserKey: string(UserKey),
}

// JWT 创建JWT认证中间件
func JWT(jwt *auth.JWT, opts JWTOptions) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 检查是否需要跳过验证
			for _, path := range opts.SkipPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next(w, r)
					return
				}
			}

			// 从请求头获取令牌
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, errors.ErrUnauthorized.Error(), http.StatusUnauthorized)
				return
			}

			// 验证令牌前缀
			if !strings.HasPrefix(authHeader, opts.TokenPrefix+" ") {
				http.Error(w, errors.ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			// 提取令牌
			token := strings.TrimPrefix(authHeader, opts.TokenPrefix+" ")

			// 验证令牌并获取用户信息
			claims, err := jwt.GetUserFromToken(token)
			if err != nil {
				switch {
				case errors.Is(err, errors.ErrTokenExpired):
					http.Error(w, errors.ErrTokenExpired.Error(), http.StatusUnauthorized)
				case errors.Is(err, errors.ErrInvalidToken):
					http.Error(w, errors.ErrInvalidToken.Error(), http.StatusUnauthorized)
				default:
					http.Error(w, errors.ErrUnauthorized.Error(), http.StatusUnauthorized)
				}
				return
			}

			// 将认证信息注入上下文
			ctx := context.WithValue(r.Context(), ContextKey(opts.ClaimsKey), claims)
			ctx = context.WithValue(ctx, ContextKey(opts.ContextUserKey), claims)
			next(w, r.WithContext(ctx))
		}
	}
}
