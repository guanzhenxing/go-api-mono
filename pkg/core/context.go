package core

import (
	"context"
	"net/http"

	"go-api-mono/pkg/logger"
)

// Context 包装了HTTP请求的上下文信息
type Context struct {
	context.Context
	Request  *http.Request
	Response *ResponseHandler
	Logger   *logger.Logger
}

// NewContext 创建一个新的上下文
func NewContext(r *http.Request, w http.ResponseWriter, log *logger.Logger) *Context {
	return &Context{
		Context:  r.Context(),
		Request:  r,
		Response: NewResponse(w, r),
		Logger:   log,
	}
}

// WithValue 返回一个带有新值的上下文
func (c *Context) WithValue(key, val interface{}) *Context {
	return &Context{
		Context:  context.WithValue(c.Context, key, val),
		Request:  c.Request,
		Response: c.Response,
		Logger:   c.Logger,
	}
}

// Set 设置上下文值
func (c *Context) Set(key, val interface{}) {
	c.Context = context.WithValue(c.Context, key, val)
}

// GetString 从上下文中获取字符串值
func (c *Context) GetString(key interface{}) string {
	if val, ok := c.Value(key).(string); ok {
		return val
	}
	return ""
}

// GetInt 从上下文中获取整数值
func (c *Context) GetInt(key interface{}) int {
	if val, ok := c.Value(key).(int); ok {
		return val
	}
	return 0
}

// GetBool 从上下文中获取布尔值
func (c *Context) GetBool(key interface{}) bool {
	if val, ok := c.Value(key).(bool); ok {
		return val
	}
	return false
}

// GetInterface 从上下文中获取接口值
func (c *Context) GetInterface(key interface{}) interface{} {
	return c.Value(key)
}
