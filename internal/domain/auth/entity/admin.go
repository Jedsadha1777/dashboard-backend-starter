package entity

import (
	"time"

	"gorm.io/gorm"
)

// Admin represents the admin user model in the database
type Admin struct {
	gorm.Model
	Email        string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Password     string    `gorm:"size:255;not null" json:"-"` // Exclude from JSON response
	TokenVersion int       `gorm:"default:1" json:"-"`         // Exclude from JSON response
	LastLogin    time.Time `json:"last_login"`
}
