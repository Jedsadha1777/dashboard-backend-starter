package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"size:255;not null"`
	Email        string         `json:"email" gorm:"size:255;not null;uniqueIndex"`
	Password     string         `json:"-" gorm:"size:255;not null"` // Exclude from JSON response
	TokenVersion int            `json:"-" gorm:"default:1"`         // For token invalidation
	LastLogin    time.Time      `json:"last_login"`
	AdminID      uint           `json:"admin_id"` // Can be optional, null if self-registered
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
