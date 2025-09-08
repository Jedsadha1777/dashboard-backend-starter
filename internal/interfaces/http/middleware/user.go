package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserRequired ensures the authenticated user is a regular user
func UserRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "user" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized: user authentication required",
			})
			return
		}
		c.Next()
	}
}

// SelfOrAdminRequired ensures the authenticated user is either the requested user or an admin
func SelfOrAdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, userExists := c.Get("user_id")
		userType, typeExists := c.Get("user_type")
		requestedID := c.Param("id")

		// Allow if user is admin
		if typeExists && userType == "admin" {
			c.Next()
			return
		}

		// Allow if user is accessing their own resource
		if userExists && userType == "user" && requestedID != "" {
			// แปลง requestedID จาก string เป็น uint เพื่อเปรียบเทียบกับ userID
			requestedIDUint, err := strconv.ParseUint(requestedID, 10, 32)
			if err == nil && userID.(uint) == uint(requestedIDUint) {
				c.Next()
				return
			}
		}

		// Deny access if neither condition is met
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Forbidden: you don't have permission to access this resource",
		})
	}
}
