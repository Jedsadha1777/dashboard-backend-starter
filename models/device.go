// models/device.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Device represents an IoT device in the system
type Device struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	DeviceID     string         `json:"device_id" gorm:"size:100;not null;uniqueIndex"`
	Name         string         `json:"name" gorm:"size:255;not null"`
	ApiKey       string         `json:"-" gorm:"size:255;not null"` // Secret, not exposed in JSON
	TokenVersion int            `json:"-" gorm:"default:1"`         // For token invalidation
	LastSeen     time.Time      `json:"last_seen"`
	Status       string         `json:"status" gorm:"size:50;default:'inactive'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// DeviceAuthInput represents device authentication request
type DeviceAuthInput struct {
	DeviceID string `json:"device_id" binding:"required"`
	ApiKey   string `json:"api_key" binding:"required"`
}
