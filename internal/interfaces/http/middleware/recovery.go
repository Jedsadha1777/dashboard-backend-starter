package middleware

import (
	"dashboard-starter/internal/domain/shared/errors"
	"dashboard-starter/internal/infrastructure/logger"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var recoveryLogger *logger.ZapLogger

// InitRecoveryLogger initializes logger for recovery middleware
func InitRecoveryLogger(l *logger.ZapLogger) {
	recoveryLogger = l
}

var (
	panicCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_panic_recovered_total",
			Help: "Total number of recovered panics",
		},
		[]string{"path", "method"},
	)
)

// CustomRecovery handles panics with proper logging
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := debug.Stack()

				// Get request ID
				requestID := GetRequestID(c)

				// Increment panic counter
				panicCounter.WithLabelValues(c.Request.URL.Path, c.Request.Method).Inc()

				// Log the panic
				if recoveryLogger != nil {
					recoveryLogger.Error("Panic recovered",
						zap.Any("error", err),
						zap.String("request_id", requestID),
						zap.String("path", c.Request.URL.Path),
						zap.String("method", c.Request.Method),
						zap.String("client_ip", c.ClientIP()),
						zap.String("user_agent", c.Request.UserAgent()),
						zap.ByteString("stack", stack),
					)
				} else {
					// Fallback logging if logger not initialized
					fmt.Printf("[PANIC] %v\nPath: %s\nMethod: %s\nStack:\n%s\n",
						err,
						c.Request.URL.Path,
						c.Request.Method,
						string(stack),
					)
				}

				// Return error response
				appErr := errors.NewAppError(
					errors.ErrCodeInternal,
					http.StatusInternalServerError,
					"An unexpected error occurred",
				).WithPath(c.Request.URL.Path).WithRequestID(requestID)

				c.AbortWithStatusJSON(appErr.StatusCode, appErr.ToResponse())
			}
		}()

		c.Next()
	}
}
