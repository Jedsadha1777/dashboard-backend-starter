package handlers

import (
	"context"
	"dashboard-starter/db"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	startTime time.Time
	version   string
}

func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// Liveness checks if the application is running
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "alive",
		"uptime":  time.Since(h.startTime).Seconds(),
		"version": h.version,
	})
}

// Readiness checks if the application is ready to serve traffic
func (h *HealthHandler) Readiness(c *gin.Context) {
	checks := make(map[string]interface{})
	allHealthy := true

	// Check database connection
	if db.DB != nil {
		sqlDB, err := db.DB.DB()
		if err != nil {
			checks["database"] = gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			}
			allHealthy = false
		} else {
			// Create context with timeout for database ping
			ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
			defer cancel()

			if err := sqlDB.PingContext(ctx); err != nil {
				checks["database"] = gin.H{
					"status": "unhealthy",
					"error":  err.Error(),
				}
				allHealthy = false
			} else {
				stats := sqlDB.Stats()
				checks["database"] = gin.H{
					"status": "healthy",
					"pool": gin.H{
						"open":       stats.OpenConnections,
						"in_use":     stats.InUse,
						"idle":       stats.Idle,
						"wait_count": stats.WaitCount,
						"wait_time":  stats.WaitDuration.Milliseconds(),
						"max_open":   stats.MaxOpenConnections,
					},
				}
			}
		}
	} else {
		checks["database"] = gin.H{
			"status": "unhealthy",
			"error":  "database not initialized",
		}
		allHealthy = false
	}

	// Memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	checks["memory"] = gin.H{
		"alloc_mb":       m.Alloc / 1024 / 1024,
		"total_alloc_mb": m.TotalAlloc / 1024 / 1024,
		"sys_mb":         m.Sys / 1024 / 1024,
		"num_gc":         m.NumGC,
		"goroutines":     runtime.NumGoroutine(),
	}

	// Overall status
	status := "ready"
	statusCode := http.StatusOK
	if !allHealthy {
		status = "not ready"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":    status,
		"checks":    checks,
		"uptime":    time.Since(h.startTime).Seconds(),
		"timestamp": time.Now().Unix(),
	})
}

// HealthSummary provides a simple health status
func (h *HealthHandler) HealthSummary(c *gin.Context) {
	// Quick health check with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbHealthy := false
	if db.DB != nil {
		if sqlDB, err := db.DB.DB(); err == nil {
			if err := sqlDB.PingContext(ctx); err == nil {
				dbHealthy = true
			}
		}
	}

	if dbHealthy {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"database": "connected",
			"uptime":   time.Since(h.startTime).Seconds(),
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "unhealthy",
			"database": "disconnected",
			"uptime":   time.Since(h.startTime).Seconds(),
		})
	}
}
