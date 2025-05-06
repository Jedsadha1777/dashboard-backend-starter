// models/token.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken represents a refresh token for authentication
type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Token     string         `json:"token" gorm:"size:255;not null;uniqueIndex"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	UserType  string         `json:"user_type" gorm:"size:50;not null"` // "admin", "user", "device", etc.
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	IsRevoked bool           `json:"is_revoked" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
