package middleware

import (
	"dashboard-starter/internal/domain/shared/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware handles all errors in a centralized way
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Get request ID if exists
			requestID, _ := c.Get("request_id")

			// Check if it's an AppError
			if appErr, ok := err.(*errors.AppError); ok {
				appErr.Path = c.Request.URL.Path
				if reqID, ok := requestID.(string); ok {
					appErr.RequestID = reqID
				}

				c.JSON(appErr.StatusCode, appErr.ToResponse())
				return
			}

			// Handle standard errors
			response := map[string]interface{}{
				"success": false,
				"error":   err.Error(),
				"code":    errors.ErrCodeInternal,
				"path":    c.Request.URL.Path,
			}

			if reqID, ok := requestID.(string); ok {
				response["request_id"] = reqID
			}

			c.JSON(http.StatusInternalServerError, response)
		}
	}
}

// HandleError is a helper to handle errors in handlers
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Add error to context
	c.Error(err)

	// Abort further processing
	c.Abort()
}
