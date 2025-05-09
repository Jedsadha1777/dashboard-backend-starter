package models

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

// UserInput represents the input data for creating or updating a user
type UserInput struct {
	Name  string `json:"name" binding:"required" validate:"required,min=2,max=255"`
	Email string `json:"email" binding:"required" validate:"required,email,max=255"`
}

// UserRegistrationInput represents the input for user registration
type UserRegistrationInput struct {
	Name            string `json:"name" binding:"required" validate:"required,min=2,max=255"`
	Email           string `json:"email" binding:"required" validate:"required,email,max=255"`
	Password        string `json:"password" binding:"required" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required,eqfield=Password"`
}

// UserLoginInput represents the login request payload for users
type UserLoginInput struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
}

// ChangePasswordInput represents the input for changing user password
type UserChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required" validate:"required,min=8,max=72,nefield=CurrentPassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required,eqfield=NewPassword"`
}
