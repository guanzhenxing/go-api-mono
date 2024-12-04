package middleware

import (
	"net/http"
	"time"

	"go-api-mono/internal/pkg/logger"

	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger(log *logger.Logger) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path
			query := r.URL.RawQuery
			method := r.Method

			defer func() {
				latency := time.Since(start)
				requestID, _ := r.Context().Value("request_id").(string)
				log.Info("HTTP Request",
					zap.String("method", method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("ip", r.RemoteAddr),
					zap.Duration("latency", latency),
					zap.String("request_id", requestID),
				)
			}()

			next(w, r)
		}
	}
}
