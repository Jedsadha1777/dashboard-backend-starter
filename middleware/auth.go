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
	// Create a new rate limiter that allows 5 requests per minute
	loginLimiter = rate.NewLimiter(rate.Every(time.Minute), 5)
)

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
		adminID, tokenVer, err := utils.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token: " + err.Error(),
			})
			return
		}

		// Find admin in database
		var admin models.Admin
		if err := db.DB.First(&admin, adminID).Error; err != nil {
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

		// Set admin ID in context for future handlers
		c.Set("admin_id", adminID)
		c.Next()
	}
}

// RateLimitMiddleware limits the number of requests from a single IP
func RateLimitMiddleware() gin.HandlerFunc {
	//  rate limiter จากค่า config
	requestsPerMinute := config.Config.RateLimit.RequestsPerMinute
	if requestsPerMinute <= 0 {
		requestsPerMinute = 60 // ค่าเริ่มต้น 60 ครั้งต่อนาที
	}

	limiter := rate.NewLimiter(rate.Every(time.Minute), requestsPerMinute)

	return func(c *gin.Context) {
		// ตรวจสอบว่าเป็น endpoint ที่ต้องการจำกัดหรือไม่
		if c.FullPath() == "/api/v1/auth/login" {
			if !limiter.Allow() {
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
