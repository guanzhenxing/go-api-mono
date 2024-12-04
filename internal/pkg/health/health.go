package health

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go-api-mono/internal/pkg/database"
)

var startTime = time.Now()

// Checker represents a health checker
type Checker struct {
	db      database.DB
	version string
}

// NewChecker creates a new health checker
func NewChecker(db database.DB, version string) *Checker {
	return &Checker{
		db:      db,
		version: version,
	}
}

// Check performs the health check
func (c *Checker) Check(ctx context.Context) Response {
	components := make([]Status, 0)
	overallStatus := "healthy"

	// Check database connection
	dbStatus := c.checkDatabase(ctx)
	components = append(components, dbStatus)
	if dbStatus.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// Check system resources
	sysStatus := c.checkSystem()
	components = append(components, sysStatus)
	if sysStatus.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	return Response{
		Status:     overallStatus,
		Version:    c.version,
		Components: components,
		Timestamp:  time.Now(),
		Uptime:     time.Since(startTime).String(),
	}
}

// checkDatabase verifies database connectivity
func (c *Checker) checkDatabase(ctx context.Context) Status {
	status := Status{
		Component: "database",
		Status:    "healthy",
	}

	// Try to ping database
	db := c.db.GetDB()
	sqlDB, err := db.DB()
	if err != nil {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Failed to get sql.DB instance: %v", err)
		return status
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Database ping failed: %v", err)
		return status
	}

	return status
}

// checkSystem verifies system resources
func (c *Checker) checkSystem() Status {
	status := Status{
		Component: "system",
		Status:    "healthy",
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Check if memory usage is within acceptable limits (e.g., < 80% of total memory)
	memoryUsage := float64(m.Alloc) / float64(m.Sys) * 100
	if memoryUsage > 80 {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("High memory usage: %.2f%%", memoryUsage)
	} else {
		status.Message = fmt.Sprintf("Memory usage: %.2f%%", memoryUsage)
	}

	return status
}
