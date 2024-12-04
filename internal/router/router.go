package router

import (
	"gin-api-mono/internal/api/admin"
	"gin-api-mono/internal/pkg/core"
	"gin-api-mono/internal/repository/mysql"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func NewHTTPMux(logger *zap.Logger, db mysql.Repo) (core.Mux, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	if db == nil {
		return nil, errors.New("db required")
	}

	mux, err := core.New(logger,
		core.WithEnableCors(),
		core.WithEnableSwagger(),
		core.WithEnablePProf(),
	)

	if err != nil {
		panic(err)
	}

	// 定义自动生成的路由组前缀为 /api
	generatedRouterGroup := mux.Group("/api")

	// 注册路由
	admin.RegisterGeneratedAdminRoutes(logger, db, generatedRouterGroup)

	return mux, nil
}
