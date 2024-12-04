package core

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-api-mono/internal/pkg/logger"
)

// Middleware 定义了中间件函数类型
type Middleware func(HandlerFunc) HandlerFunc

// ServerOptions 定义服务器选项
type ServerOptions struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Logger       *logger.Logger
}

// Server 封装HTTP服务器
type Server struct {
	server      *http.Server
	logger      *logger.Logger
	mux         *http.ServeMux
	middlewares []Middleware
}

// Group 路由组
type Group struct {
	prefix      string
	server      *Server
	parent      *Group
	middlewares []Middleware
}

// NewServer 创建一个新的服务器实例
func NewServer(opts ServerOptions) *Server {
	mux := http.NewServeMux()
	srv := &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", opts.Port),
			Handler:      mux,
			ReadTimeout:  opts.ReadTimeout,
			WriteTimeout: opts.WriteTimeout,
		},
		logger:      opts.Logger,
		mux:         mux,
		middlewares: make([]Middleware, 0),
	}
	return srv
}

// Use 添加全局中间件
func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

// Group 创建一个新的路由组
func (s *Server) Group(prefix string) *Group {
	return &Group{
		prefix:      prefix,
		server:      s,
		middlewares: make([]Middleware, 0),
	}
}

// Group 创建子路由组
func (g *Group) Group(prefix string) *Group {
	return &Group{
		prefix:      g.prefix + prefix,
		server:      g.server,
		parent:      g,
		middlewares: make([]Middleware, 0),
	}
}

// Use 为路由组添加中间件
func (g *Group) Use(middlewares ...Middleware) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// Handle 注册路由处理器
func (g *Group) Handle(method, pattern string, handler HandlerFunc) {
	// 收集所有中间件
	var allMiddlewares []Middleware
	allMiddlewares = append(allMiddlewares, g.server.middlewares...)

	// 收集父组的中间件
	parent := g
	for parent != nil {
		allMiddlewares = append(allMiddlewares, parent.middlewares...)
		parent = parent.parent
	}

	// 创建中间件链
	finalHandler := handler
	for i := len(allMiddlewares) - 1; i >= 0; i-- {
		finalHandler = allMiddlewares[i](finalHandler)
	}

	// 移除前导和尾部斜杠以避免重定向
	pattern = strings.Trim(pattern, "/")
	prefix := strings.Trim(g.prefix, "/")

	// 构建完整的路由模式
	var fullPattern string
	if pattern == "" {
		fullPattern = fmt.Sprintf("%s %s", method, "/"+prefix)
	} else {
		fullPattern = fmt.Sprintf("%s %s", method, "/"+prefix+"/"+pattern)
	}

	g.server.mux.Handle(fullPattern, WrapHandler(finalHandler, g.server.logger))
}

// HandleFunc 注册路由处理器
func (s *Server) HandleFunc(pattern string, handler HandlerFunc) {
	s.mux.Handle(pattern, WrapHandler(handler, s.logger))
}

// Start 启动服务器
func (s *Server) Start() error {
	s.logger.Info(fmt.Sprintf("Server is starting on %s", s.server.Addr))
	return s.server.ListenAndServe()
}

// Stop 优雅关闭服务器
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Server is shutting down...")
	return s.server.Shutdown(ctx)
}

// Router 返回路由器实例
func (s *Server) Router() *Server {
	return s
}
