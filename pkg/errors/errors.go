package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// ErrorCode 定义错误码类型
type ErrorCode int

const (
	// 系统级错误码 (1-99)
	ErrCodeUnknown ErrorCode = iota + 1
	ErrCodeInternal
	ErrCodeDatabase
	ErrCodeCache
	ErrCodeConfig

	// 请求相关错误码 (100-199)
	ErrCodeBadRequest = iota + 100
	ErrCodeUnauthorized
	ErrCodeForbidden
	ErrCodeNotFound
	ErrCodeMethodNotAllowed
	ErrCodeConflict
	ErrCodeTooManyRequests
	ErrCodeValidation

	// 认证相关错误码 (200-299)
	ErrCodeInvalidToken = iota + 200
	ErrCodeTokenExpired
	ErrCodeInvalidCredentials
	ErrCodeUserNotFound
	ErrCodeUserExists
	ErrCodeInvalidSignature

	// 业务相关错误码 (300-399)
	ErrCodeInvalidOperation = iota + 300
	ErrCodeResourceNotFound
	ErrCodeResourceExists
	ErrCodeResourceUnavailable
)

// Error 定义自定义错误类型
type Error struct {
	Code     ErrorCode   `json:"code"`              // 错误码
	Message  string      `json:"message"`           // 错误信息
	Details  interface{} `json:"details,omitempty"` // 错误详情
	HTTPCode int         `json:"-"`                 // HTTP状态码
	stack    []uintptr   // 错误堆栈
	cause    error       // 原始错误
}

// Error 实现error接口
func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.cause)
	}
	if e.Details != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Details)
	}
	return e.Message
}

// WithDetails 添加错误详情
func (e *Error) WithDetails(details interface{}) *Error {
	e.Details = details
	return e
}

// WithCause 添加原始错误
func (e *Error) WithCause(err error) *Error {
	e.cause = err
	return e
}

// Cause 获取原始错误
func (e *Error) Cause() error {
	return e.cause
}

// StackTrace 获取堆栈信息
func (e *Error) StackTrace() string {
	var builder strings.Builder
	frames := runtime.CallersFrames(e.stack)
	for {
		frame, more := frames.Next()
		builder.WriteString(fmt.Sprintf("\n\t%s:%d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	return builder.String()
}

// New 创建新的错误
func New(code ErrorCode, message string) *Error {
	err := &Error{
		Code:     code,
		Message:  message,
		HTTPCode: errorCodeToHTTPStatus(code),
		stack:    make([]uintptr, 32),
	}
	// 捕获堆栈信息
	runtime.Callers(2, err.stack)
	return err
}

// Wrap 包装已有错误
func Wrap(err error, code ErrorCode, message string) *Error {
	if err == nil {
		return nil
	}
	wrappedErr := New(code, message)
	wrappedErr.cause = err
	return wrappedErr
}

// Unwrap 解包错误
func Unwrap(err error) error {
	if e, ok := err.(*Error); ok {
		return e.cause
	}
	return err
}

// Is 判断错误类型
func Is(err, target error) bool {
	if err == target {
		return true
	}
	if e, ok := err.(*Error); ok {
		return Is(e.cause, target)
	}
	return false
}

// As 类型断言
func As(err error, target interface{}) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	if e, ok := err.(*Error); ok {
		if As(e.cause, target) {
			return true
		}
	}
	return false
}

// errorCodeToHTTPStatus 将错误码转换为HTTP状态码
func errorCodeToHTTPStatus(code ErrorCode) int {
	switch {
	case code >= 1 && code < 100:
		return http.StatusInternalServerError
	case code >= 100 && code < 200:
		switch code {
		case ErrCodeBadRequest:
			return http.StatusBadRequest
		case ErrCodeUnauthorized:
			return http.StatusUnauthorized
		case ErrCodeForbidden:
			return http.StatusForbidden
		case ErrCodeNotFound:
			return http.StatusNotFound
		case ErrCodeMethodNotAllowed:
			return http.StatusMethodNotAllowed
		case ErrCodeConflict:
			return http.StatusConflict
		case ErrCodeTooManyRequests:
			return http.StatusTooManyRequests
		case ErrCodeValidation:
			return http.StatusUnprocessableEntity
		default:
			return http.StatusBadRequest
		}
	case code >= 200 && code < 300:
		return http.StatusUnauthorized
	case code >= 300 && code < 400:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// 预定义错误
var (
	// 系统错误
	ErrUnknown  = New(ErrCodeUnknown, "Unknown error")
	ErrInternal = New(ErrCodeInternal, "Internal server error")
	ErrDatabase = New(ErrCodeDatabase, "Database error")
	ErrCache    = New(ErrCodeCache, "Cache error")
	ErrConfig   = New(ErrCodeConfig, "Configuration error")

	// 请求错误
	ErrBadRequest       = New(ErrCodeBadRequest, "Bad request")
	ErrUnauthorized     = New(ErrCodeUnauthorized, "Unauthorized")
	ErrForbidden        = New(ErrCodeForbidden, "Forbidden")
	ErrNotFound         = New(ErrCodeNotFound, "Not found")
	ErrMethodNotAllowed = New(ErrCodeMethodNotAllowed, "Method not allowed")
	ErrConflict         = New(ErrCodeConflict, "Conflict")
	ErrTooManyRequests  = New(ErrCodeTooManyRequests, "Too many requests")
	ErrValidation       = New(ErrCodeValidation, "Validation failed")

	// 认证错误
	ErrInvalidToken       = New(ErrCodeInvalidToken, "Invalid token")
	ErrTokenExpired       = New(ErrCodeTokenExpired, "Token expired")
	ErrInvalidCredentials = New(ErrCodeInvalidCredentials, "Invalid credentials")
	ErrUserNotFound       = New(ErrCodeUserNotFound, "User not found")
	ErrUserExists         = New(ErrCodeUserExists, "User already exists")
	ErrInvalidSignature   = New(ErrCodeInvalidSignature, "Invalid signature")

	// 业务错误
	ErrInvalidOperation    = New(ErrCodeInvalidOperation, "Invalid operation")
	ErrResourceNotFound    = New(ErrCodeResourceNotFound, "Resource not found")
	ErrResourceExists      = New(ErrCodeResourceExists, "Resource already exists")
	ErrResourceUnavailable = New(ErrCodeResourceUnavailable, "Resource unavailable")
)

// IsNotFound 判断是否是"未找到"错误
func IsNotFound(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrCodeNotFound || e.Code == ErrCodeUserNotFound || e.Code == ErrCodeResourceNotFound
	}
	return false
}

// IsValidation 判断是否是验证错误
func IsValidation(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrCodeValidation
	}
	return false
}

// IsUnauthorized 判断是否是未授权错误
func IsUnauthorized(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrCodeUnauthorized || e.Code == ErrCodeInvalidToken || e.Code == ErrCodeTokenExpired
	}
	return false
}

// IsForbidden 判断是否是禁止访问错误
func IsForbidden(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrCodeForbidden
	}
	return false
}
