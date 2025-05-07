package middleware

import (
	"dashboard-starter/config"
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/utils"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	// Map เก็บ rate limiters แยกตาม IP
	ipLimiters    = make(map[string]*rate.Limiter)
	limitersMutex sync.RWMutex
)

// init initializes the package-level variables
func init() {
	go cleanupIPLimiters()
}

// ฟังก์ชันล้าง map เป็นระยะเพื่อประหยัดหน่วยความจำ
func cleanupIPLimiters() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		limitersMutex.Lock()
		// สร้าง map ใหม่ ล้าง IP เก่าทั้งหมด
		ipLimiters = make(map[string]*rate.Limiter)
		limitersMutex.Unlock()

		log.Println("IP rate limiter map has been cleaned up")
	}
}

// ดึง limiter สำหรับ IP ที่กำหนด
func getIPLimiter(ip string) *rate.Limiter {
	limitersMutex.RLock()
	limiter, exists := ipLimiters[ip]
	limitersMutex.RUnlock()

	if !exists {
		limitersMutex.Lock()
		// ตรวจสอบอีกครั้งเพื่อป้องกัน race condition
		limiter, exists = ipLimiters[ip]
		if !exists {
			requestsPerMinute := config.Config.RateLimit.RequestsPerMinute
			if requestsPerMinute <= 0 {
				requestsPerMinute = 60 // ค่าเริ่มต้น
			}
			limiter = rate.NewLimiter(rate.Every(time.Minute), requestsPerMinute)
			ipLimiters[ip] = limiter
		}
		limitersMutex.Unlock()
	}

	return limiter
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
	return func(c *gin.Context) {
		currentPath := c.FullPath()
		shouldLimit := false

		for _, path := range config.Config.RateLimit.LimitedPaths {
			if currentPath == path {
				shouldLimit = true
				break
			}
		}

		if shouldLimit {
			// เปลี่ยนเป็นใช้ IP-based limiter
			ip := c.ClientIP()
			requestsPerMinute := config.Config.RateLimit.RequestsPerMinute
			if requestsPerMinute <= 0 {
				requestsPerMinute = 60
			}

			limiter := getIPLimiter(ip, requestsPerMinute)

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
