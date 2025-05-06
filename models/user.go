package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:255;not null"`
	Email     string         `json:"email" gorm:"size:255;not null;uniqueIndex"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserInput represents the input data for creating or updating a user
type UserInput struct {
	Name  string `json:"name" binding:"required" validate:"required,min=2,max=255"`
	Email string `json:"email" binding:"required" validate:"required,email,max=255"`
}
