package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-api-mono/internal/app/user/controller"
	"go-api-mono/internal/app/user/service"
	"go-api-mono/internal/pkg/config"
	"go-api-mono/internal/pkg/database"
	"go-api-mono/internal/pkg/health"
	"go-api-mono/internal/pkg/logger"
	"go-api-mono/internal/pkg/middleware"
	"go-api-mono/internal/pkg/validator"
)

// App represents the application
type App struct {
	config     *config.Config
	logger     *logger.Logger
	db         database.DB
	health     *health.Checker
	httpServer *http.Server
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	// Initialize logger
	log, err := logger.New(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize database
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

	// Initialize validator
	if err := validator.Register(); err != nil {
		return nil, fmt.Errorf("failed to register validator: %w", err)
	}

	// Initialize health checker
	healthChecker := health.NewChecker(db, cfg.App.Version)

	return &App{
		config: cfg,
		logger: log,
		db:     db,
		health: healthChecker,
	}, nil
}

// Start starts the application
func (a *App) Start() error {
	// Set Gin mode
	if a.config.App.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if a.config.App.Mode == "testing" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize HTTP server
	router := gin.New()

	// Register middleware
	router.Use(
		middleware.Recovery(a.logger),
		middleware.Logger(a.logger),
		middleware.Cors(),
		middleware.RequestID(),
		middleware.RateLimit(a.config.RateLimit),
	)

	// Register routes
	a.setupHTTPServer(router)

	a.httpServer = &http.Server{
		Addr:    a.config.Server.Addr,
		Handler: router,
	}

	a.logger.Info("Starting server", zap.String("addr", a.config.Server.Addr))
	return a.httpServer.ListenAndServe()
}

// Stop stops the application
func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("Shutting down server...")
	return a.httpServer.Shutdown(ctx)
}

// setupHTTPServer sets up the HTTP server routes
func (a *App) setupHTTPServer(router *gin.Engine) {
	// Register health check
	healthHandler := health.NewHandler(a.health)
	healthHandler.Register(router)

	// API routes
	v1 := router.Group("/api/v1")

	// Auth routes
	v1.POST("/auth/register", controller.NewUserController(
		service.NewUserService(a.db),
	).Create)

	v1.POST("/auth/login", func(c *gin.Context) {
		// TODO: Implement login
		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
		})
	})

	// User routes
	userService := service.NewUserService(a.db)
	userController := controller.NewUserController(userService)
	userController.Register(v1)
}
