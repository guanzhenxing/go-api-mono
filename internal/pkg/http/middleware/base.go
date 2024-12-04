package middleware

import (
	"net/http"

	"go-api-mono/internal/pkg/logger"
)

// HandlerFunc 定义了处理函数的类型
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Middleware 定义了中间件函数类型
type Middleware func(HandlerFunc) HandlerFunc

// Chain 将多个中间件串联成一个
func Chain(middlewares ...Middleware) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// WrapHandler 将处理函数包装为中间件
func WrapHandler(handler HandlerFunc, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}
