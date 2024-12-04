package middleware

import (
	"context"
	"net/http"

	"go-api-mono/internal/pkg/utils"
)

// RequestID 请求ID中间件
func RequestID() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = utils.GenerateRequestID()
			}
			w.Header().Set("X-Request-ID", requestID)
			ctx := r.Context()
			ctx = context.WithValue(ctx, RequestIDKey, requestID)
			next(w, r.WithContext(ctx))
		}
	}
}
