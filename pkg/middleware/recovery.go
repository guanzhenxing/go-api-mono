package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-api-mono/pkg/logger"
)

// Recovery 创建一个恢复中间件，用于捕获 panic
// 它会：
// - 捕获任何 panic
// - 记录错误日志和堆栈跟踪
// - 返回 500 错误响应
func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误和堆栈跟踪
				log.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
