package core

import (
	"net/http"

	"go-api-mono/pkg/logger"
)

// HandlerFunc 定义了处理函数的类型
type HandlerFunc func(*Context)

// WrapHandler 将HandlerFunc包装为http.HandlerFunc
func WrapHandler(handler HandlerFunc, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(r, w, log)
		handler(ctx)
	}
}

// WrapMiddleware 将http.Handler包装为HandlerFunc
func WrapMiddleware(handler http.Handler) HandlerFunc {
	return func(c *Context) {
		handler.ServeHTTP(c.Response.Writer, c.Request)
	}
}

// Chain 将多个HandlerFunc串联成一个
func Chain(handlers ...HandlerFunc) HandlerFunc {
	return func(c *Context) {
		for _, handler := range handlers {
			handler(c)
		}
	}
}

// Adapt 将http.HandlerFunc适配为HandlerFunc
func Adapt(handler http.HandlerFunc) HandlerFunc {
	return func(c *Context) {
		handler(c.Response.Writer, c.Request)
	}
}

// AsHandler 将HandlerFunc转换为http.Handler
func AsHandler(handler HandlerFunc, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(r, w, log)
		handler(ctx)
	})
}
