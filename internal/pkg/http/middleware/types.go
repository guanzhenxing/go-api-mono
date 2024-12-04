package middleware

// ContextKey 定义上下文键类型
type ContextKey string

const (
	// RequestIDKey 请求ID的上下文键
	RequestIDKey ContextKey = "request_id"
	// ClaimsKey JWT声明的上下文键
	ClaimsKey ContextKey = "claims"
	// UserKey 用户信息的上下文键
	UserKey ContextKey = "user"
)
