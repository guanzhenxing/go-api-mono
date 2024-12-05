package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-api-mono/pkg/logger"
)

// Logger 创建一个日志中间件
// 它会记录每个请求的详细信息，包括：
// - 请求方法和路径
// - 响应状态码
// - 处理时间
// - 请求ID
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
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
