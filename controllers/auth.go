package controllers

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/services"
	"dashboard-starter/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginInput represents the login request payload
type LoginInput struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       uint      `json:"user_id"`
	UserType     string    `json:"user_type"`
}

// Login handles admin authentication
func Login(c *gin.Context) {
	var input LoginInput

	// Parse request body
	if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	// Validate input
	if err := utils.ValidateStruct(input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Find admin by email
	var admin models.Admin
	if err := db.DB.Where("email = ?", input.Email).First(&admin).Error; err != nil {
		// Use a generic error message to prevent user enumeration
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		// Use a generic error message to prevent user enumeration
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}

	// Update token version inside transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		admin.TokenVersion += 1
		admin.LastLogin = time.Now()
		return tx.Save(&admin).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Authentication failed: " + err.Error(),
		})
		return
	}

	// Generate JWT token
	token, exp, err := utils.GenerateToken(admin.ID, "admin", admin.TokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to generate token: " + err.Error(),
		})
		return
	}

	// Generate refresh token
	refreshTokenObj, err := services.CreateRefreshToken(admin.ID, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to generate refresh token: " + err.Error(),
		})
		return
	}

	// Return token
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: LoginResponse{
			Token:        token,
			RefreshToken: refreshTokenObj.Token,
			ExpiresAt:    exp,
			UserID:       admin.ID,
			UserType:     "admin",
		},
	})
}

// Logout handles admin logout
func Logout(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	// Find admin
	var admin models.Admin
	if err := db.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Admin not found",
		})
		return
	}

	// Invalidate token by incrementing token version
	err := db.Transaction(func(tx *gorm.DB) error {
		admin.TokenVersion += 1
		return tx.Save(&admin).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Logout failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    gin.H{"message": "Logged out successfully"},
	})
}

// RefreshInput represents the refresh token request payload
type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken handles access token refresh using a valid refresh token
func RefreshToken(c *gin.Context) {
	var input RefreshInput

	// Parse request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	// Validate refresh token
	refreshToken, err := services.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Invalid refresh token: " + err.Error(),
		})
		return
	}

	// Get token version for the user
	tokenVersion, err := services.GetUserTokenVersion(refreshToken.UserID, refreshToken.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to verify user: " + err.Error(),
		})
		return
	}

	// Generate new access token
	token, exp, err := utils.GenerateToken(refreshToken.UserID, refreshToken.UserType, tokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to generate token: " + err.Error(),
		})
		return
	}

	// Return new access token
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"token":      token,
			"expires_at": exp,
			"user_id":    refreshToken.UserID,
			"user_type":  refreshToken.UserType,
		},
	})
}

// GetProfile retrieves the admin's profile
func GetProfile(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var admin models.Admin
	if err := db.DB.Select("id, email, created_at, updated_at, last_login").First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Admin not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    admin,
	})
}
