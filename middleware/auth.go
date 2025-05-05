package middleware

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddlewear() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		adminID, tokenVer, err := utils.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		var admin models.Admin
		if err := db.DB.First(&admin, adminID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Admin not found"})
			return
		}

		// ตรวจว่า token_version ใน DB ตรงกับ token
		if tokenVer != admin.TokenVersion {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Set("admin_id", adminID)
		c.Next()
	}
}
