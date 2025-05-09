package models

import (
	"time"

	"gorm.io/gorm"
)

// Admin represents the admin user model in the database
type Admin struct {
	gorm.Model             // Embeds ID, CreatedAt, UpdatedAt, DeletedAt
	Email        string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Password     string    `gorm:"size:255;not null" json:"-"` // Exclude from JSON response
	TokenVersion int       `gorm:"default:1" json:"-"`         // Exclude from JSON response
	LastLogin    time.Time `json:"last_login"`
}

// AdminInput represents the input for creating or updating an admin
type AdminInput struct {
	Email    string `json:"email" binding:"required" validate:"required,email,max=255"`
	Password string `json:"password" binding:"required" validate:"required,min=8,max=72"`
}

// ChangePasswordInput represents the input for changing admin password
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required" validate:"required,min=8,max=72,nefield=CurrentPassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword" validate:"required,eqfield=NewPassword"`
}
