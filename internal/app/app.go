package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-api-mono/internal/app/user/controller"
	"go-api-mono/internal/app/user/repository"
	"go-api-mono/internal/app/user/service"
	"go-api-mono/internal/pkg/auth"
	"go-api-mono/internal/pkg/config"
	"go-api-mono/internal/pkg/core"
	"go-api-mono/internal/pkg/database"
	"go-api-mono/internal/pkg/http/middleware"
	"go-api-mono/internal/pkg/logger"
	"go-api-mono/internal/pkg/security"

	"go.uber.org/zap"
)

// Options 定义应用程序选项
type Options struct {
	ConfigFile string
	LogLevel   string
	DevMode    bool
}

// App 应用程序结构
type App struct {
	config  *config.Config
	logger  *logger.Logger
	mux     *http.ServeMux
	db      *database.DB
	server  *http.Server
	jwt     *auth.JWT
	limiter *security.IPRateLimiter
}

// New 创建新的应用程序实例
func New(opts Options) (*App, error) {
	app := &App{}

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	app.config = cfg

	// 创建日志记录器
	log, err := logger.New(logger.LogConfig{
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
	app.logger = log

	// 创建数据库连接
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
	app.db = db

	// 创建JWT实例
	jwt := auth.New(auth.Config{
		SigningKey:     cfg.JWT.SigningKey,
		ExpirationTime: cfg.JWT.ExpirationTime,
		SigningMethod:  cfg.JWT.SigningMethod,
		TokenPrefix:    cfg.JWT.TokenPrefix,
	})
	app.jwt = jwt

	// 创建速率限制器
	limiter := security.NewIPRateLimiter(float64(cfg.RateLimit.Requests), float64(cfg.RateLimit.Burst))
	app.limiter = limiter

	// 创建HTTP服务器
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.ShutdownTimeout,
	}
	app.server = server
	app.mux = mux

	// 初始化路由
	if err := app.initRoutes(); err != nil {
		return nil, fmt.Errorf("failed to initialize routes: %w", err)
	}

	return app, nil
}

// initRoutes 初始化路由
func (a *App) initRoutes() error {
	// 创建用户仓储
	userRepo := repository.NewUserRepository(a.db)

	// 创建用户服务
	userService := service.NewUserService(userRepo)

	// 创建用户控制器
	userController := controller.NewUserController(userService, a.jwt)

	// 创建中间件链
	chain := middleware.Chain(
		middleware.Recovery(a.logger),
		middleware.Logger(a.logger),
		middleware.RequestID(),
		middleware.CORS(),
	)

	// 创建需要认证的中间件链
	protectedChain := middleware.Chain(
		middleware.Recovery(a.logger),
		middleware.Logger(a.logger),
		middleware.RequestID(),
		middleware.CORS(),
		middleware.JWT(a.jwt, middleware.DefaultJWTOptions),
		middleware.RateLimit(a.limiter, middleware.DefaultRateLimitOptions),
	)

	// 注册路由
	a.mux.HandleFunc("POST /api/v1/auth/login", a.wrapHandler(userController.Login, chain))
	a.mux.HandleFunc("POST /api/v1/auth/register", a.wrapHandler(userController.Register, chain))
	a.mux.HandleFunc("GET /api/v1/users", a.wrapHandler(userController.List, protectedChain))
	a.mux.HandleFunc("GET /api/v1/users/{id}", a.wrapHandler(userController.Get, protectedChain))
	a.mux.HandleFunc("PUT /api/v1/users/{id}", a.wrapHandler(userController.Update, protectedChain))
	a.mux.HandleFunc("DELETE /api/v1/users/{id}", a.wrapHandler(userController.Delete, protectedChain))

	return nil
}

// wrapHandler 包装控制器处理函数
func (a *App) wrapHandler(h func(*core.Context), m middleware.Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := core.NewContext(r, w, a.logger)
		m(func(w http.ResponseWriter, r *http.Request) {
			h(ctx)
		})(w, r)
	}
}

// Run 运行应用程序
func Run(opts Options) error {
	app, err := New(opts)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	return app.Run()
}

// Run 运行应用程序
func (a *App) Run() error {
	// 创建错误通道
	errChan := make(chan error, 1)

	// 启动HTTP服务器
	go func() {
		a.logger.Info("Starting server", zap.String("address", a.server.Addr))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start server: %w", err)
		}
	}()

	// 等待信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case <-quit:
		a.logger.Info("Shutting down server...")

		// 创建关闭上下文
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// 关闭HTTP服务器
		if err := a.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		// 关闭数据库连接
		if err := a.db.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}

		a.logger.Info("Server stopped")
	}

	return nil
}

// Stop 停止应用程序
func (a *App) Stop() error {
	// 创建关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout)
	defer cancel()

	// 关闭HTTP服务器
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	// 关闭数据库连接
	if err := a.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// InitConfig 初始化配置
func (app *App) InitConfig() error {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	app.config = cfg
	return nil
}
