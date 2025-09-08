package middleware

import (
	"dashboard-starter/internal/domain/shared/errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestSizeLimitMiddleware limits the size of incoming requests
func RequestSizeLimitMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Content-Length header first
		if c.Request.ContentLength > maxBytes {
			appErr := errors.NewAppError(
				errors.ErrCodeBadRequest,
				http.StatusRequestEntityTooLarge,
				fmt.Sprintf("Request body too large. Max size: %d bytes", maxBytes),
			).WithPath(c.Request.URL.Path)

			c.AbortWithStatusJSON(appErr.StatusCode, appErr.ToResponse())
			return
		}

		// Limit the reader
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

		c.Next()

		// Check if there was an error reading the body
		if c.Errors.Last() != nil {
			if c.Errors.Last().Error() == "http: request body too large" {
				appErr := errors.NewAppError(
					errors.ErrCodeBadRequest,
					http.StatusRequestEntityTooLarge,
					"Request body exceeded maximum size",
				).WithPath(c.Request.URL.Path)

				c.AbortWithStatusJSON(appErr.StatusCode, appErr.ToResponse())
				return
			}
		}
	}
}

// FileUploadLimitMiddleware specific limit for file upload endpoints
func FileUploadLimitMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
