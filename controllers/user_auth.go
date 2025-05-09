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

// UserRegister handles the user registration
func UserRegister(c *gin.Context) {
	var input models.UserRegistrationInput

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

	// Check if email already exists
	var existingUser models.User
	if err := db.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		// Email exists
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Email is already registered",
		})
		return
	}

	// Generate password hash
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Registration failed: " + err.Error(),
		})
		return
	}

	// Create new user
	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     string(hashed),
		TokenVersion: 1,
		LastLogin:    time.Now(),
	}

	// Save user to database
	err = db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Registration failed: " + err.Error(),
		})
		return
	}

	// Generate tokens
	token, exp, err := utils.GenerateToken(user.ID, "user", user.TokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to generate token: " + err.Error(),
		})
		return
	}

	// Generate refresh token
	refreshTokenObj, err := services.CreateRefreshToken(user.ID, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to generate refresh token: " + err.Error(),
		})
		return
	}

	// Return token and user data (excluding password)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data: LoginResponse{
			Token:        token,
			RefreshToken: refreshTokenObj.Token,
			ExpiresAt:    exp,
			UserID:       user.ID,
			UserType:     "user",
		},
	})
}

// UserLogin handles the user login
func UserLogin(c *gin.Context) {
	var input models.UserLoginInput

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

	// Find user by email
	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// Use a generic error message to prevent user enumeration
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		// Use a generic error message to prevent user enumeration
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}

	// Update token version inside transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		user.TokenVersion += 1
		user.LastLogin = time.Now()
		return tx.Save(&user).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Authentication failed: " + err.Error(),
		})
		return
	}

	// Generate JWT token
	token, exp, err := utils.GenerateToken(user.ID, "user", user.TokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to generate token: " + err.Error(),
		})
		return
	}

	// Generate refresh token
	refreshTokenObj, err := services.CreateRefreshToken(user.ID, "user")
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
			UserID:       user.ID,
			UserType:     "user",
		},
	})
}

// UserLogout handles the user logout
func UserLogout(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Find user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Invalidate token by incrementing token version
	err := db.Transaction(func(tx *gorm.DB) error {
		user.TokenVersion += 1
		return tx.Save(&user).Error
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

// GetUserProfile retrieves the user's profile
func GetUserProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := db.DB.Select("id, name, email, created_at, updated_at, last_login").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

// ChangeUserPassword handles the user password change
func ChangeUserPassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input models.UserChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
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

	// Find user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Current password is incorrect",
		})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to update password: " + err.Error(),
		})
		return
	}

	// Update password and invalidate tokens
	err = db.Transaction(func(tx *gorm.DB) error {
		user.Password = string(hashedPassword)
		user.TokenVersion += 1 // Invalidate all tokens
		return tx.Save(&user).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to update password: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    gin.H{"message": "Password updated successfully"},
	})
}
