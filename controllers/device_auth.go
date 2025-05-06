// controllers/device_auth.go
package controllers

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/services"
	"dashboard-starter/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeviceAuth authenticates an IoT device
func DeviceAuth(c *gin.Context) {
	var input models.DeviceAuthInput

	// Parse request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	// Find device by deviceID
	var device models.Device
	if err := db.DB.Where("device_id = ?", input.DeviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Invalid device ID or API key",
		})
		return
	}

	// Verify API key (simple comparison for this example)
	// In production, you might want to use a more secure comparison
	if device.ApiKey != input.ApiKey {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Invalid device ID or API key",
		})
		return
	}

	// Update device last seen status
	err := db.Transaction(func(tx *gorm.DB) error {
		device.LastSeen = time.Now()
		device.Status = "active"
		return tx.Save(&device).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Authentication failed: " + err.Error(),
		})
		return
	}

	// Generate access token (30 minutes as in your diagram)
	token, exp, err := utils.GenerateToken(device.ID, "device", device.TokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to generate token: " + err.Error(),
		})
		return
	}

	// Generate refresh token (1 year as in your diagram)
	refreshTokenObj, err := services.CreateRefreshToken(device.ID, "device")
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to generate refresh token: " + err.Error(),
		})
		return
	}

	// Return tokens
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: LoginResponse{
			Token:        token,
			RefreshToken: refreshTokenObj.Token,
			ExpiresAt:    exp,
			UserID:       device.ID,
			UserType:     "device",
		},
	})
}
