package middleware

import (
	"dashboard-starter/config"
	"dashboard-starter/db"
	authEntity "dashboard-starter/internal/domain/auth/entity"
	deviceEntity "dashboard-starter/internal/domain/device/entity"
	userEntity "dashboard-starter/internal/domain/user/entity"
	"dashboard-starter/utils"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPLimiter struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

var (
	ipLimiters        = make(map[string]*IPLimiter)
	limitersMutex     sync.Mutex
	cleanupInterval   = 5 * time.Minute
	inactiveThreshold = 20 * time.Minute
)

func init() {
	if config.Config.RateLimit.CleanupMinutes <= 0 {
		cleanupInterval = 5 * time.Minute
	} else {
		cleanupInterval = time.Duration(config.Config.RateLimit.CleanupMinutes) * time.Minute
	}

	if config.Config.RateLimit.InactiveMinutes <= 0 {
		inactiveThreshold = 20 * time.Minute
	} else {
		inactiveThreshold = time.Duration(config.Config.RateLimit.InactiveMinutes) * time.Minute
	}
	go cleanupIPLimiters()
}

func cleanupIPLimiters() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Error during IP limiter cleanup: %v", r)
				}
			}()

			threshold := time.Now().Add(-1 * inactiveThreshold)

			limitersMutex.Lock()
			defer limitersMutex.Unlock()

			beforeCleanup := len(ipLimiters)

			var keysToRemove []string

			for ip, limiter := range ipLimiters {
				if limiter.lastAccess.Before(threshold) {
					keysToRemove = append(keysToRemove, ip)
				}
			}

			for _, ip := range keysToRemove {
				delete(ipLimiters, ip)
			}

			afterCleanup := len(ipLimiters)
			if beforeCleanup != afterCleanup {
				log.Printf("IP rate limiter cleanup: removed %d inactive limiters, %d remaining", beforeCleanup-afterCleanup, afterCleanup)
			}
		}()
	}
}

func getIPLimiter(ip string) *IPLimiter {
	limitersMutex.Lock()
	defer limitersMutex.Unlock()

	ipLimiter, exists := ipLimiters[ip]
	if !exists {
		requestsPerMinute := config.Config.RateLimit.RequestsPerMinute
		if requestsPerMinute <= 0 {
			requestsPerMinute = 60
		}
		limiter := rate.NewLimiter(rate.Limit(requestsPerMinute)/60, requestsPerMinute)
		ipLimiter = &IPLimiter{
			limiter:    limiter,
			lastAccess: time.Now(),
		}
		ipLimiters[ip] = ipLimiter
	} else {
		ipLimiter.lastAccess = time.Now()
	}

	return ipLimiter
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authentication required. Please provide a valid Bearer token",
			})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		userID, userType, tokenVer, err := utils.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token: " + err.Error(),
			})
			return
		}

		if userType == "admin" {
			var admin authEntity.Admin
			if err := db.DB.First(&admin, userID).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Admin account not found",
				})
				return
			}

			if tokenVer != admin.TokenVersion {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Token has been revoked. Please login again",
				})
				return
			}
		} else if userType == "user" {
			var user userEntity.User
			if err := db.DB.First(&user, userID).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "User account not found",
				})
				return
			}

			if tokenVer != user.TokenVersion {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Token has been revoked. Please login again",
				})
				return
			}
		} else if userType == "device" {
			var device deviceEntity.Device
			if err := db.DB.First(&device, userID).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Device not found",
				})
				return
			}

			if tokenVer != device.TokenVersion {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Token has been revoked. Please register device again",
				})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unknown user type",
			})
			return
		}

		c.Set("user_id", userID)
		c.Set("user_type", userType)

		if userType == "admin" {
			c.Set("admin_id", userID)
		}

		c.Next()
	}
}

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
			ip := c.ClientIP()
			ipLimiter := getIPLimiter(ip)

			if !ipLimiter.limiter.Allow() {
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
