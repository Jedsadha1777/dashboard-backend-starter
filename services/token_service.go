// services/token_service.go
package services

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CreateRefreshToken creates and stores a new refresh token
func CreateRefreshToken(userID uint, userType string) (*models.RefreshToken, error) {
	// Generate a JWT refresh token
	tokenString, expiresAt, err := utils.GenerateRefreshToken(userID, userType)
	if err != nil {
		return nil, err
	}

	// Create token inside a transaction
	var refreshToken models.RefreshToken
	err = db.Transaction(func(tx *gorm.DB) error {
		// Create new refresh token
		refreshToken = models.RefreshToken{
			Token:     tokenString,
			UserID:    userID,
			UserType:  userType,
			ExpiresAt: expiresAt,
			IsRevoked: false,
		}

		if err := tx.Create(&refreshToken).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// ValidateRefreshToken checks if a refresh token is valid and not revoked
func ValidateRefreshToken(tokenString string) (*models.RefreshToken, error) {
	// Parse the token to get user information
	userID, userType, err := utils.ParseRefreshToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Find the token in database
	var refreshToken models.RefreshToken
	if err := db.DB.Where("token = ? AND user_id = ? AND user_type = ? AND is_revoked = ? AND expires_at > ?",
		tokenString, userID, userType, false, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// GetUserTokenVersion retrieves the token version for a user
func GetUserTokenVersion(userID uint, userType string) (int, error) {
	switch userType {
	case "admin":
		var admin models.Admin
		if err := db.DB.Select("token_version").First(&admin, userID).Error; err != nil {
			return 0, err
		}
		return admin.TokenVersion, nil
	case "user":
		// If regular users also have token versions, implement here
		// For now, assume users don't have token versions
		return 1, nil
		// Add more cases for other user types
	case "device":
		// เพิ่มการรองรับสำหรับ device
		var device models.Device
		if err := db.DB.Select("token_version").First(&device, userID).Error; err != nil {
			return 0, err
		}
		return device.TokenVersion, nil
	default:
		return 0, errors.New("unsupported user type")
	}
}

// RevokeRefreshToken marks a refresh token as revoked
func RevokeRefreshToken(tokenString string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&models.RefreshToken{}).
			Where("token = ?", tokenString).
			Update("is_revoked", true).Error
	})
}

// RevokeAllRefreshTokens revokes all refresh tokens for a user
func RevokeAllRefreshTokens(userID uint, userType string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&models.RefreshToken{}).
			Where("user_id = ? AND user_type = ? AND is_revoked = ?", userID, userType, false).
			Update("is_revoked", true).Error
	})
}

// CleanupExpiredTokens removes all expired refresh tokens
func CleanupExpiredTokens() error {
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("expires_at < ? OR is_revoked = ?", time.Now(), true).
			Delete(&models.RefreshToken{}).Error
	})
}
