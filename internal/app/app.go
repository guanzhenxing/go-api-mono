package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-api-mono/internal/app/user/controller"
	"go-api-mono/internal/app/user/service"
	"go-api-mono/internal/pkg/config"
	"go-api-mono/internal/pkg/database"
	"go-api-mono/internal/pkg/logger"
	"go-api-mono/internal/pkg/middleware"
)

// App represents the application
type App struct {
	config     *config.Config
	logger     *logger.Logger
	db         database.DB
	engine     *gin.Engine
	httpServer *http.Server
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	// 初始化日志
	l, err := logger.New(logger.LogConfig{
		Level:      cfg.Log.Level,
		Filename:   cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// 初始化数据库
	db, err := database.New(database.DatabaseConfig{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		Username:     cfg.Database.Username,
		Password:     cfg.Database.Password,
		Database:     cfg.Database.Database,
		MaxOpenConns: cfg.Database.MaxOpenConns,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		MaxLifetime:  cfg.Database.MaxLifetime,
		Debug:        cfg.Database.Debug,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	app := &App{
		config: cfg,
		logger: l,
		db:     db,
	}

	// 初始化 HTTP 服务器
	if err := app.setupHTTPServer(); err != nil {
		return nil, err
	}

	return app, nil
}

// setupHTTPServer initializes the HTTP server and routes
func (a *App) setupHTTPServer() error {
	// 初始化 Gin
	if a.config.App.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	a.engine = gin.New()

	// 添加中间件
	a.engine.Use(
		middleware.Recovery(a.logger),
		middleware.Logger(a.logger),
		middleware.Cors(),
		middleware.RequestID(),
		middleware.RateLimit(middleware.RateLimitConfig{
			Requests: a.config.RateLimit.Requests,
			Burst:    a.config.RateLimit.Burst,
		}),
	)

	// 创建用户服务和控制器
	userService := service.NewUserService(a.db)
	userController := controller.NewUserController(userService)

	// 设置路由
	v1 := a.engine.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"time":   time.Now().Format(time.RFC3339),
			})
		})

		// 认证路由
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userController.Create)
			auth.POST("/login", func(c *gin.Context) {
				var loginReq struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				}
				if err := c.ShouldBindJSON(&loginReq); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				// TODO: 实现实际的登录逻辑
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": "OK",
					"data": gin.H{
						"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwicm9sZSI6InVzZXIiLCJleHAiOjE3MzMzODk1MTgsIm5iZiI6MTczMzMwMzExOCwiaWF0IjoxNzMzMzAzMTE4fQ.rlWO9ijjmeTxxvnc6A07iep-Z0ZNVmMpU_N3zrQxjrc",
					},
				})
			})
		}

		// 注册用户路由
		userController.Register(v1)
	}

	// 创建 HTTP 服务器
	a.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", a.config.Server.Port),
		Handler:        a.engine,
		ReadTimeout:    a.config.Server.ReadTimeout,
		WriteTimeout:   a.config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return nil
}

// Start starts the application
func (a *App) Start() error {
	a.logger.Info("Starting server", zap.String("addr", a.httpServer.Addr))
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

// Stop gracefully stops the application
func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("Shutting down server...")
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	return nil
}
