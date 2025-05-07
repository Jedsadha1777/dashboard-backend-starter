package models

import (
	"time"

	"gorm.io/gorm"
)

// Article represents a content article in the system
type Article struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:255;not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Slug        string         `json:"slug" gorm:"size:255;not null;uniqueIndex"`
	Summary     string         `json:"summary" gorm:"size:500"`
	Status      string         `json:"status" gorm:"size:20;default:'draft'"` // draft, published, archived
	PublishedAt *time.Time     `json:"published_at"`
	AdminID     uint           `json:"admin_id" gorm:"not null"`
	Admin       Admin          `json:"admin" gorm:"foreignKey:AdminID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// ArticleInput represents the input data for creating or updating an article
type ArticleInput struct {
	Title       string `json:"title" binding:"required" validate:"required,min=3,max=255"`
	Content     string `json:"content" binding:"required" validate:"required,min=10"`
	Slug        string `json:"slug" binding:"required" validate:"required,min=3,max=255"`
	Summary     string `json:"summary" validate:"max=500"`
	Status      string `json:"status" validate:"oneof=draft published archived"`
	PublishedAt string `json:"published_at"` // Optional, in ISO 8601 format
}
