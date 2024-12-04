package main

import (
	"context"
	"fmt"
	"net/http"

	"gin-api-mono/configs"
	"gin-api-mono/internal/pkg/env"
	"gin-api-mono/internal/pkg/logger"
	"gin-api-mono/internal/pkg/shutdown"
	"gin-api-mono/internal/pkg/timeutil"
	"gin-api-mono/internal/repository/mysql"
	"gin-api-mono/internal/router"

	"go.uber.org/zap"
)

// @title gin-api-mono 接口文档
// @version v0.0.1

// @BasePath /
func main() {
	// 初始化 access logger
	accessLogger, err := logger.NewJSONLogger(
		logger.WithOutputInConsole(),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileP(configs.ProjectAccessLogFile),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		_ = accessLogger.Sync()
	}()

	// 初始化 MySQL
	dbRepo, err := mysql.New()
	if err != nil {
		panic(err)
	}

	// 初始化 HTTP 服务
	mux, err := router.NewHTTPMux(accessLogger, dbRepo)
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    configs.ProjectPort,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			accessLogger.Fatal("http server startup err", zap.Error(err))
		}
	}()

	// 优雅关闭
	shutdown.Close(
		func() {
			// 关闭 http server
			if err := server.Shutdown(context.TODO()); err != nil {
				accessLogger.Error("server shutdown err", zap.Error(err))
			}

			// 关闭 db master (支持读写)
			if err := dbRepo.DbWClose(); err != nil {
				accessLogger.Error("dbw close err", zap.Error(err))
			}

			// 关闭 db slave (仅支持读)
			if err := dbRepo.DbRClose(); err != nil {
				accessLogger.Error("dbr close err", zap.Error(err))
			}
		},
	)
}
