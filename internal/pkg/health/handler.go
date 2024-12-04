package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler represents the health check HTTP handler
type Handler struct {
	checker *Checker
}

// NewHandler creates a new health check handler
func NewHandler(checker *Checker) *Handler {
	return &Handler{
		checker: checker,
	}
}

// Register registers health check routes
func (h *Handler) Register(router *gin.Engine) {
	router.GET("/api/v1/health", h.Check)
}

// Check handles the health check request
func (h *Handler) Check(c *gin.Context) {
	resp := h.checker.Check(c.Request.Context())
	if resp.Status == "healthy" {
		c.JSON(http.StatusOK, resp)
	} else {
		c.JSON(http.StatusServiceUnavailable, resp)
	}
}
