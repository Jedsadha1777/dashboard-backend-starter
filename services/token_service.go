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

// UserTypeProvider defines an interface for handling different user types
type UserTypeProvider interface {
	GetTokenVersion(userID uint) (int, error)
	GetModelName() string
}

// AdminProvider implements UserTypeProvider for admin users
type AdminProvider struct{}

func (p *AdminProvider) GetTokenVersion(userID uint) (int, error) {
	var admin models.Admin
	if err := db.DB.Select("token_version").First(&admin, userID).Error; err != nil {
		return 0, err
	}
	return admin.TokenVersion, nil
}

func (p *AdminProvider) GetModelName() string {
	return "admin"
}

// DeviceProvider implements UserTypeProvider for IoT devices
type DeviceProvider struct{}

func (p *DeviceProvider) GetTokenVersion(userID uint) (int, error) {
	var device models.Device
	if err := db.DB.Select("token_version").First(&device, userID).Error; err != nil {
		return 0, err
	}
	return device.TokenVersion, nil
}

func (p *DeviceProvider) GetModelName() string {
	return "device"
}

// UserProvider implements UserTypeProvider for regular users
type UserProvider struct{}

func (p *UserProvider) GetTokenVersion(userID uint) (int, error) {
	// Regular users don't have token versioning yet
	// This is a placeholder and can be implemented later
	return 1, nil
}

func (p *UserProvider) GetModelName() string {
	return "user"
}

// Registry of all user type providers
var userTypeProviders = map[string]UserTypeProvider{
	"admin":  &AdminProvider{},
	"user":   &UserProvider{},
	"device": &DeviceProvider{},
}

// RegisterUserTypeProvider adds a new user type provider to the registry
func RegisterUserTypeProvider(userType string, provider UserTypeProvider) {
	userTypeProviders[userType] = provider
}

// GetUserTypeProvider retrieves the provider for a specific user type
func GetUserTypeProvider(userType string) (UserTypeProvider, error) {
	provider, exists := userTypeProviders[userType]
	if !exists {
		return nil, errors.New("unsupported user type")
	}
	return provider, nil
}

// CreateRefreshToken creates and stores a new refresh token
func CreateRefreshToken(userID uint, userType string) (*models.RefreshToken, error) {
	// Check if the user type is supported
	if _, err := GetUserTypeProvider(userType); err != nil {
		return nil, err
	}

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

	// Check if the user type is supported
	if _, err := GetUserTypeProvider(userType); err != nil {
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
	provider, err := GetUserTypeProvider(userType)
	if err != nil {
		return 0, err
	}

	return provider.GetTokenVersion(userID)
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
	// Check if the user type is supported
	if _, err := GetUserTypeProvider(userType); err != nil {
		return err
	}

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

// RotateUserTokenVersion increases the token version for a specific user
// This can be used to invalidate all existing access tokens
func RotateUserTokenVersion(userID uint, userType string) error {
	provider, err := GetUserTypeProvider(userType)
	if err != nil {
		return err
	}

	switch provider.GetModelName() {
	case "admin":
		return db.Transaction(func(tx *gorm.DB) error {
			var admin models.Admin
			if err := tx.First(&admin, userID).Error; err != nil {
				return err
			}
			admin.TokenVersion += 1
			return tx.Save(&admin).Error
		})
	case "device":
		return db.Transaction(func(tx *gorm.DB) error {
			var device models.Device
			if err := tx.First(&device, userID).Error; err != nil {
				return err
			}
			device.TokenVersion += 1
			return tx.Save(&device).Error
		})
	case "user":
		// Users don't have token versioning yet
		return nil
	default:
		return errors.New("unsupported user type for token rotation")
	}
}
