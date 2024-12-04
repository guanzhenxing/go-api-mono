package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID creates a request ID middleware
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID is already set
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in header and context
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("X-Request-ID", requestID)

		c.Next()
	}
}
