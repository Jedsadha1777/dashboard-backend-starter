package controllers

import (
	"dashboard-starter/models"
	"dashboard-starter/services"
	"dashboard-starter/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateArticle handles the request to create a new article
func CreateArticle(c *gin.Context) {
	// Get admin ID from context
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Unauthorized: admin authentication required",
		})
		return
	}

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
		// Create a slug from title
		slug := strings.ToLower(input.Title)
		// Replace non-alphanumeric with dash
		slug = strings.ReplaceAll(slug, " ", "-")
		// Remove non-alphanumeric
		input.Slug = slug
	}

	// Set default status if empty
	if input.Status == "" {
		input.Status = "draft"
	}

	article, err := services.CreateArticle(&input, adminID.(uint))
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

// ListArticles handles the request to list articles with pagination and filters
func ListArticles(c *gin.Context) {
	// Get admin ID from context
	_, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Unauthorized: admin authentication required",
		})
		return
	}

	search := c.Query("search")
	status := c.Query("status")

	// Get page and limit with defaults
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	articles, total, err := services.GetArticles(search, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve articles: " + err.Error(),
		})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    articles,
		Meta: gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// GetArticle handles the request to get an article by ID
func GetArticle(c *gin.Context) {
	// Get admin ID from context
	_, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Unauthorized: admin authentication required",
		})
		return
	}

	id := c.Param("id")

	article, err := services.GetArticleByID(id)
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
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Unauthorized: admin authentication required",
		})
		return
	}

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

	article, err := services.UpdateArticle(id, &input, adminID.(uint))
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
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Unauthorized: admin authentication required",
		})
		return
	}

	id := c.Param("id")

	err := services.DeleteArticle(id, adminID.(uint))
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
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Unauthorized: admin authentication required",
		})
		return
	}

	id := c.Param("id")

	article, err := services.PublishArticle(id, adminID.(uint))
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
