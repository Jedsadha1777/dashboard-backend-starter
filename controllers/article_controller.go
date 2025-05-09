package controllers

import (
	"dashboard-starter/models"
	"dashboard-starter/services"
	"dashboard-starter/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateArticle handles the request to create a new article
func CreateArticle(c *gin.Context) {
	// Get admin ID from context
	adminID, _ := c.Get("admin_id")

	var input models.ArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	// Validate input
	if err := utils.ValidateStruct(input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Generate slug if not provided
	if input.Slug == "" {
		slug, err := utils.GenerateUniqueSlug(input.Title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Success: false,
				Error:   "Failed to generate slug: " + err.Error(),
			})
			return
		}
		input.Slug = slug
	}

	// Set default status if empty
	if input.Status == "" {
		input.Status = "draft"
	}

	// Create service and call
	articleService := services.NewArticleService()
	article, err := articleService.CreateArticle(&input, adminID.(uint))

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to create article: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    article,
	})
}

// ... rest of the controller remains the same
