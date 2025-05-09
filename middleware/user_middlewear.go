package middleware

import (
	"net/http"

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
			// Convert to string for comparison
			if userID.(uint) == gin.H{"id": requestedID}["id"] {
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