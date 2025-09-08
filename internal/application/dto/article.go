package dto

type ArticleInput struct {
	Title       string `json:"title" binding:"required" validate:"required,min=3,max=255"`
	Content     string `json:"content" binding:"required" validate:"required,min=10"`
	Slug        string `json:"slug" binding:"required" validate:"required,min=3,max=255"`
	Summary     string `json:"summary" validate:"max=500"`
	Status      string `json:"status" validate:"oneof=draft published archived"`
	PublishedAt string `json:"published_at"` // Optional, in ISO 8601 format
}
