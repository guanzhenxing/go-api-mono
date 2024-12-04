package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-api-mono/internal/pkg/logger"
)

// Logger creates a logger middleware
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method

		c.Next()

		cost := time.Since(start)
		log.Info("Request completed",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
			zap.Duration("cost", cost),
			zap.String("request_id", c.GetString("X-Request-ID")),
		)
	}
}
