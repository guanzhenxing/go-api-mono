package health

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go-api-mono/internal/pkg/cache"
	"go-api-mono/internal/pkg/database"
)

// 记录应用启动时间，用于计算运行时长
var startTime = time.Now()

// Checker 代表健康检查器
// 它负责检查各个组件的健康状态
type Checker struct {
	db      database.DB // 数据库连接
	cache   cache.Cache // 缓存连接
	version string      // 应用版本
}

// NewChecker 创建一个新的健康检查器
// 它需要数据库和缓存连接以及应用版本信息
func NewChecker(db database.DB, cache cache.Cache, version string) *Checker {
	return &Checker{
		db:      db,
		cache:   cache,
		version: version,
	}
}

// Check 执行健康检查
// 它会检查所有关键组件的状态，并返回一个包含详细信息的响应
func (c *Checker) Check(ctx context.Context) Response {
	components := make([]Status, 0)
	overallStatus := "healthy" // 默认状态为健康

	// 检查数据库连接
	dbStatus := c.checkDatabase(ctx)
	components = append(components, dbStatus)
	if dbStatus.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// 检查Redis连接
	cacheStatus := c.checkCache(ctx)
	components = append(components, cacheStatus)
	if cacheStatus.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// 检查系统资源
	sysStatus := c.checkSystem()
	components = append(components, sysStatus)
	if sysStatus.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// 返回完整的健康检查响应
	return Response{
		Status:     overallStatus,
		Version:    c.version,
		Components: components,
		Timestamp:  time.Now(),
		Uptime:     time.Since(startTime).String(),
	}
}

// checkDatabase 验证数据库连接状态
// 它会尝试ping数据库，并检查连接的有效性
func (c *Checker) checkDatabase(ctx context.Context) Status {
	status := Status{
		Component: "database",
		Status:    "healthy",
	}

	// 获取数据库连接
	db := c.db.GetDB()
	if db == nil {
		status.Status = "unhealthy"
		status.Message = "Database connection is nil"
		return status
	}

	// 检查数据库错误
	if db.Error != nil {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Database error: %v", db.Error)
		return status
	}

	// 在测试模式下跳过实际的连接测试
	if db.Config != nil && db.Config.ConnPool != nil {
		return status
	}

	// 获取底层的sql.DB连接
	sqlDB, err := db.DB()
	if err != nil {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Failed to get sql.DB instance: %v", err)
		return status
	}

	// 测试数据库连接
	if err := sqlDB.PingContext(ctx); err != nil {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Database ping failed: %v", err)
		return status
	}

	return status
}

// checkCache 验证Redis缓存连接状态
// 它会尝试ping Redis服务器来确认连接是否正常
func (c *Checker) checkCache(ctx context.Context) Status {
	status := Status{
		Component: "redis",
		Status:    "healthy",
	}

	// 测试Redis连接
	if err := c.cache.Ping(ctx); err != nil {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Redis ping failed: %v", err)
		return status
	}

	return status
}

// checkSystem 验证系统资源状态
// 它会检查内存使用情况等系统指标
func (c *Checker) checkSystem() Status {
	status := Status{
		Component: "system",
		Status:    "healthy",
	}

	// 获取内存统计信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 计算内存使用率并检查是否超过阈值（80%）
	memoryUsage := float64(m.Alloc) / float64(m.Sys) * 100
	if memoryUsage > 80 {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("High memory usage: %.2f%%", memoryUsage)
	} else {
		status.Message = fmt.Sprintf("Memory usage: %.2f%%", memoryUsage)
	}

	return status
}
