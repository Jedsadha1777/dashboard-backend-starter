package entity

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:255;not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Slug        string         `json:"slug" gorm:"size:255;not null;uniqueIndex"`
	Summary     string         `json:"summary" gorm:"size:500"`
	Status      string         `json:"status" gorm:"size:20;default:'draft'"`
	PublishedAt *time.Time     `json:"published_at"`
	AdminID     uint           `json:"admin_id" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
