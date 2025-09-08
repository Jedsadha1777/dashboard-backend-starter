package dto

import "time"

// Auth DTOs
type AdminInput struct {
	Email    string `json:"email" binding:"required" validate:"required,email,max=255"`
	Password string `json:"password" binding:"required" validate:"required,min=8,max=72"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
}

type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       uint      `json:"user_id"`
	UserType     string    `json:"user_type"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required" validate:"required,min=8,max=72,nefield=CurrentPassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword" validate:"required,eqfield=NewPassword"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type DeviceAuthInput struct {
	DeviceID string `json:"device_id" binding:"required"`
	ApiKey   string `json:"api_key" binding:"required"`
}
