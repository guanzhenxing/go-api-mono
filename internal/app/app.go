package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"go-api-mono/internal/app/user/controller"
	"go-api-mono/internal/app/user/service"
	"go-api-mono/internal/pkg/cache"
	"go-api-mono/internal/pkg/config"
	"go-api-mono/internal/pkg/database"
	"go-api-mono/internal/pkg/health"
	"go-api-mono/internal/pkg/logger"
	"go-api-mono/internal/pkg/middleware"
	"go-api-mono/internal/pkg/validator"
)

// LoginRequest 定义登录请求的结构
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 定义登录响应的结构
type LoginResponse struct {
	Token string `json:"token"`
}

// App 代表应用程序的核心结构
// 它持有所有关键组件的引用，如配置、日志、数据库等
type App struct {
	config     *config.Config  // 应用配置
	logger     *logger.Logger  // 日志组件
	db         database.DB     // 数据库连接
	cache      cache.Cache     // 缓存客户端
	health     *health.Checker // 健康检查器
	httpServer *http.Server    // HTTP服务器
}

// New 创建一个新的应用实例
// 它会初始化所有必要的组件，如果任何组件初始化失败，将返回错误
func New(cfg *config.Config) (*App, error) {
	// 初始化日志组件
	log, err := logger.New(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 初始化数据库连接
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
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 初始化Redis缓存
	redisCache, err := cache.NewRedisCache(cache.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %w", err)
	}

	// 初始化验证器
	if err := validator.Register(); err != nil {
		return nil, fmt.Errorf("failed to register validator: %w", err)
	}

	// 初始化健康检查器
	healthChecker := health.NewChecker(db, redisCache, cfg.App.Version)

	return &App{
		config: cfg,
		logger: log,
		db:     db,
		cache:  redisCache,
		health: healthChecker,
	}, nil
}

// Start 启动应用程序
// 它会启动HTTP服务器并开始处理请求
func (a *App) Start() error {
	// 设置Gin的运行模式
	if a.config.App.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if a.config.App.Mode == "testing" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化HTTP服务器
	router := gin.New()

	// 注册中间件
	router.Use(
		middleware.Recovery(a.logger),            // 错误恢复中间件
		middleware.Logger(a.logger),              // 日志中间件
		middleware.Cors(),                        // CORS中间件
		middleware.RequestID(),                   // 请求ID中间件
		middleware.RateLimit(a.config.RateLimit), // 限流中间件
	)

	// 注册路由
	a.setupHTTPServer(router)

	// 配置HTTP服务器
	a.httpServer = &http.Server{
		Addr:    a.config.Server.Addr,
		Handler: router,
	}

	// 启动服务器
	a.logger.Info("Starting server", zap.String("addr", a.config.Server.Addr))
	return a.httpServer.ListenAndServe()
}

// Stop 优雅地停止应用程序
// 它会关闭所有活跃的连接并等待它们完成
func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("Shutting down server...")

	// 关闭Redis连接
	if err := a.cache.Close(); err != nil {
		a.logger.Error("Failed to close redis connection", zap.Error(err))
	}

	// 关闭HTTP服务器
	return a.httpServer.Shutdown(ctx)
}

// setupHTTPServer 设置HTTP服务器的路由
// 它注册所有的API端点和处理器
func (a *App) setupHTTPServer(router *gin.Engine) {
	// 注册健康检查路由
	healthHandler := health.NewHandler(a.health)
	healthHandler.Register(router)

	// API路由组
	v1 := router.Group("/api/v1")

	// 用户服务
	userService := service.NewUserService(a.db)

	// 认证路由
	v1.POST("/auth/register", controller.NewUserController(userService).Create)
	v1.POST("/auth/login", a.handleLogin(userService))

	// 用户路由
	userController := controller.NewUserController(userService)
	userController.Register(v1)
}

// handleLogin 处理用户登录请求
func (a *App) handleLogin(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			if validationErrors := validator.FormatError(err); len(validationErrors) > 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "Validation failed",
					"errors":  validationErrors,
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		// 验证用户凭据
		user, err := userService.Authenticate(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid credentials",
			})
			return
		}

		// 生成JWT令牌
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})

		// 签名令牌
		tokenString, err := token.SignedString([]byte(a.config.JWT.SigningKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Failed to generate token",
			})
			return
		}

		// 返回令牌
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Login successful",
			"data": LoginResponse{
				Token: tokenString,
			},
		})
	}
}
