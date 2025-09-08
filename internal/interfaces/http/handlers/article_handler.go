package handlers

import (
	"dashboard-starter/internal/application/dto"
	"dashboard-starter/internal/application/services"
	"dashboard-starter/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService *services.ArticleApplicationService
}

func NewArticleHandler(articleService *services.ArticleApplicationService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Admin ID not found",
		})
		return
	}

	var input dto.ArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	if err := utils.ValidateStruct(input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	article, err := h.articleService.CreateArticle(input, adminID.(uint))
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

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	articles, total, err := h.articleService.ListArticles(page, limit, search, status)
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

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id := c.Param("id")
	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid article ID",
		})
		return
	}

	article, err := h.articleService.GetArticle(uint(articleID))
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

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Admin ID not found",
		})
		return
	}

	id := c.Param("id")
	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid article ID",
		})
		return
	}

	var input dto.ArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	if err := utils.ValidateStruct(input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	article, err := h.articleService.UpdateArticle(uint(articleID), input, adminID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "permission denied" {
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

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Admin ID not found",
		})
		return
	}

	id := c.Param("id")
	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid article ID",
		})
		return
	}

	if err := h.articleService.DeleteArticle(uint(articleID), adminID.(uint)); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "permission denied" {
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

func (h *ArticleHandler) PublishArticle(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Admin ID not found",
		})
		return
	}

	id := c.Param("id")
	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid article ID",
		})
		return
	}

	article, err := h.articleService.PublishArticle(uint(articleID), adminID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "permission denied" {
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
