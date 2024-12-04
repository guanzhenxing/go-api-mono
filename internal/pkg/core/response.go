package core

import (
	"encoding/json"
	"net/http"

	"go-api-mono/internal/pkg/errors"
)

// H 是一个快捷的 map[string]interface{} 类型
type H map[string]interface{}

// Response 定义了统一的响应格式
type Response struct {
	Code    int         `json:"code"`               // 错误码
	Message string      `json:"message"`            // 响应信息
	Data    interface{} `json:"data,omitempty"`     // 响应数据
	Details interface{} `json:"details,omitempty"`  // 错误详情
	TraceID string      `json:"trace_id,omitempty"` // 请求追踪ID
}

// ResponseHandler 处理HTTP响应
type ResponseHandler struct {
	Writer  http.ResponseWriter
	Request *http.Request
	status  int
}

// NewResponse 创建一个新的响应处理器
func NewResponse(w http.ResponseWriter, r *http.Request) *ResponseHandler {
	return &ResponseHandler{
		Writer:  w,
		Request: r,
		status:  http.StatusOK,
	}
}

// Status 设置响应状态码
func (r *ResponseHandler) Status(code int) *ResponseHandler {
	r.status = code
	return r
}

// JSON 发送JSON响应
func (r *ResponseHandler) JSON(data interface{}) {
	r.Writer.Header().Set("Content-Type", "application/json")
	r.Writer.WriteHeader(r.status)

	resp := Response{
		Code:    r.status,
		Message: http.StatusText(r.status),
		Data:    data,
	}

	// 获取 trace_id，如果不存在则忽略
	if traceID, ok := r.Request.Context().Value("trace_id").(string); ok {
		resp.TraceID = traceID
	}

	if err := json.NewEncoder(r.Writer).Encode(resp); err != nil {
		http.Error(r.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Error 发送错误响应
func (r *ResponseHandler) Error(err error) {
	r.Writer.Header().Set("Content-Type", "application/json")

	var resp Response
	if e, ok := err.(*errors.Error); ok {
		r.Writer.WriteHeader(e.HTTPCode)
		resp = Response{
			Code:    int(e.Code),
			Message: e.Message,
			Details: e.Details,
		}
	} else {
		r.Writer.WriteHeader(http.StatusInternalServerError)
		resp = Response{
			Code:    int(errors.ErrCodeUnknown),
			Message: err.Error(),
		}
	}

	// 获取 trace_id，如果不存在则忽略
	if traceID, ok := r.Request.Context().Value("trace_id").(string); ok {
		resp.TraceID = traceID
	}

	if err := json.NewEncoder(r.Writer).Encode(resp); err != nil {
		http.Error(r.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Success 发送成功响应
func (r *ResponseHandler) Success(data interface{}) {
	r.Status(http.StatusOK).JSON(data)
}

// Created 发送创建成功响应
func (r *ResponseHandler) Created(data interface{}) {
	r.Status(http.StatusCreated).JSON(data)
}

// NoContent 发送无内容响应
func (r *ResponseHandler) NoContent() {
	r.Writer.WriteHeader(http.StatusNoContent)
}

// BadRequest 发送请求错误响应
func (r *ResponseHandler) BadRequest(message string, details interface{}) {
	r.Error(errors.New(errors.ErrCodeBadRequest, message).WithDetails(details))
}

// Unauthorized 发送未授权响应
func (r *ResponseHandler) Unauthorized(message string, details interface{}) {
	r.Error(errors.New(errors.ErrCodeUnauthorized, message).WithDetails(details))
}

// Forbidden 发送禁止访问响应
func (r *ResponseHandler) Forbidden(message string, details interface{}) {
	r.Error(errors.New(errors.ErrCodeForbidden, message).WithDetails(details))
}

// NotFound 发送未找到响应
func (r *ResponseHandler) NotFound(message string, details interface{}) {
	r.Error(errors.New(errors.ErrCodeNotFound, message).WithDetails(details))
}

// Conflict 发送冲突响应
func (r *ResponseHandler) Conflict(message string, details interface{}) {
	r.Error(errors.New(errors.ErrCodeConflict, message).WithDetails(details))
}

// InternalError 发送内部错误响应
func (r *ResponseHandler) InternalError(message string, details interface{}) {
	r.Error(errors.New(errors.ErrCodeInternal, message).WithDetails(details))
}

// ValidationError 发送验证错误响应
func (r *ResponseHandler) ValidationError(message string, details interface{}) {
	r.Error(errors.New(errors.ErrCodeValidation, message).WithDetails(details))
}

// DatabaseError 发送数据库错误响应
func (r *ResponseHandler) DatabaseError(message string, details interface{}) {
	r.Error(errors.New(errors.ErrCodeDatabase, message).WithDetails(details))
}

// TooManyRequests 发送请求过多响应
func (r *ResponseHandler) TooManyRequests(message string, details interface{}) {
	r.Error(errors.New(errors.ErrCodeTooManyRequests, message).WithDetails(details))
}

// ServiceUnavailable 发送服务不可用响应
func (r *ResponseHandler) ServiceUnavailable(message string, details interface{}) {
	r.Error(errors.New(errors.ErrCodeInternal, message).WithDetails(details))
}
