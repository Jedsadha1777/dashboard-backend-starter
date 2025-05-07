package middleware

import (
	"dashboard-starter/config"
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	// Global rate limiter that will be initialized in init()
	loginLimiter *rate.Limiter
)

// init initializes the package-level variables
func init() {
	// Default to 5 requests per minute if config is not yet loaded
	// This will be overridden when the app starts and config is loaded
	loginLimiter = rate.NewLimiter(rate.Every(time.Minute), 5)
}

// UpdateRateLimiters updates rate limiters based on current config
// Should be called after config is fully loaded
func UpdateRateLimiters() {
	requestsPerMinute := config.Config.RateLimit.RequestsPerMinute
	if requestsPerMinute <= 0 {
		requestsPerMinute = 60 // Default to 60 requests per minute
	}

	// Update the global limiter
	loginLimiter = rate.NewLimiter(rate.Every(time.Minute), requestsPerMinute)
}

// AuthMiddleware checks if the request has a valid JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")

		// Check if the header exists and has the Bearer prefix
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authentication required. Please provide a valid Bearer token",
			})
			return
		}

		// Extract token from header
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		userID, userType, tokenVer, err := utils.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token: " + err.Error(),
			})
			return
		}

		// Verify token version based on user type
		if userType == "admin" {
			var admin models.Admin
			if err := db.DB.First(&admin, userID).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Admin account not found",
				})
				return
			}

			// Verify token version matches the one in database
			if tokenVer != admin.TokenVersion {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Token has been revoked. Please login again",
				})
				return
			}
		} else if userType == "user" {
			// Handle regular users if they have token versioning
			// For now, we're assuming they don't
		} else if userType == "device" {
			var device models.Device
			if err := db.DB.First(&device, userID).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Device not found",
				})
				return
			}

			// Verify token version matches the one in database
			if tokenVer != device.TokenVersion {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Token has been revoked. Please register device again",
				})
				return
			}
		}
		// Add more user types as needed

		// Set user info in context for future handlers
		c.Set("user_id", userID)
		c.Set("user_type", userType)

		// For backward compatibility
		if userType == "admin" {
			c.Set("admin_id", userID)
		}

		c.Next()
	}
}

// AdminRequired ensures the authenticated user is an admin
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized: admin authentication required",
			})
			return
		}
		c.Next()
	}
}

// RateLimitMiddleware limits the number of requests from a single IP
func RateLimitMiddleware() gin.HandlerFunc {
	// Use the global rate limiter which is configured from env variables
	// This ensures we're using the same configured limiter everywhere
	return func(c *gin.Context) {
		// Check if current path should be rate-limited
		currentPath := c.FullPath()
		shouldLimit := false

		for _, path := range config.Config.RateLimit.LimitedPaths {
			if currentPath == path {
				shouldLimit = true
				break
			}
		}

		if shouldLimit {
			if !loginLimiter.Allow() {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"success": false,
					"error":   "Rate limit exceeded. Please try again later",
				})
				return
			}
		}
		c.Next()
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
