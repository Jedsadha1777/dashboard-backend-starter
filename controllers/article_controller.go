package controllers

import (
	"dashboard-starter/db"
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
		// First generate a base slug
		baseSlug, err := utils.GenerateSlug(input.Title, 100)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Success: false,
				Error:   "Failed to generate slug: " + err.Error(),
			})
			return
		}

		// Then ensure it's unique in the database
		uniqueSlug, err := utils.EnsureUniqueSlug(db.DB, baseSlug, "articles", "slug")
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Success: false,
				Error:   "Failed to generate unique slug: " + err.Error(),
			})
			return
		}

		input.Slug = uniqueSlug
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

// ListArticles handles the request to list articles with pagination and search
func ListArticles(c *gin.Context) {
	// Convert query parameters to PaginationParams
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		// Use default values
		params = utils.NewPaginationParams()
	}

	// Call service
	articleService := services.NewArticleService()
	articles, pagination, err := articleService.GetArticles(params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve articles: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    articles,
		Meta:    pagination,
	})
}

// GetArticle handles the request to get an article by ID
func GetArticle(c *gin.Context) {
	id := c.Param("id")

	// Create service and call
	articleService := services.NewArticleService()
	article, err := articleService.GetByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Article not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    article,
	})
}

// UpdateArticle handles the request to update an article
func UpdateArticle(c *gin.Context) {
	// Get admin ID from context
	adminID, _ := c.Get("admin_id")

	id := c.Param("id")

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

	// Update article
	articleService := services.NewArticleService()
	article, err := articleService.UpdateArticle(id, &input, adminID.(uint))

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "you don't have permission to update this article" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, Response{
			Success: false,
			Error:   "Failed to update article: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    article,
	})
}

// DeleteArticle handles the request to delete an article
func DeleteArticle(c *gin.Context) {
	// Get admin ID from context
	adminID, _ := c.Get("admin_id")

	id := c.Param("id")

	// Delete article
	articleService := services.NewArticleService()
	err := articleService.DeleteArticle(id, adminID.(uint))

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "you don't have permission to delete this article" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, Response{
			Success: false,
			Error:   "Failed to delete article: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    gin.H{"message": "Article deleted successfully"},
	})
}

// PublishArticle handles the request to publish an article
func PublishArticle(c *gin.Context) {
	// Get admin ID from context
	adminID, _ := c.Get("admin_id")

	id := c.Param("id")

	// Publish article
	articleService := services.NewArticleService()
	article, err := articleService.PublishArticle(id, adminID.(uint))

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "you don't have permission to publish this article" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, Response{
			Success: false,
			Error:   "Failed to publish article: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    article,
	})
}
