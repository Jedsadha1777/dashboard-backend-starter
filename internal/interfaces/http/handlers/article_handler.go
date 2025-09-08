package handlers

import (
	"dashboard-starter/internal/application/dto"
	"dashboard-starter/internal/application/services"
	"dashboard-starter/internal/domain/shared/errors"
	"dashboard-starter/internal/interfaces/http/middleware"
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
		middleware.HandleError(c, errors.ErrUnauthorized)
		return
	}

	var input dto.ArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.HandleError(c, errors.ValidationError(map[string]string{
			"input": err.Error(),
		}))
		return
	}

	if err := utils.ValidateStruct(input); err != nil {
		middleware.HandleError(c, errors.ValidationError(map[string]string{
			"validation": err.Error(),
		}))
		return
	}

	article, err := h.articleService.CreateArticle(input, adminID.(uint))
	if err != nil {
		middleware.HandleError(c, errors.WrapError(err, errors.ErrCodeInternal, 500))
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
		middleware.HandleError(c, errors.DatabaseError(err))
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
		middleware.HandleError(c, errors.ErrBadRequest)
		return
	}

	article, err := h.articleService.GetArticle(uint(articleID))
	if err != nil {
		middleware.HandleError(c, errors.ErrNotFound)
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
		middleware.HandleError(c, errors.ErrUnauthorized)
		return
	}

	id := c.Param("id")
	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		middleware.HandleError(c, errors.ErrBadRequest)
		return
	}

	var input dto.ArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.HandleError(c, errors.ValidationError(map[string]string{
			"input": err.Error(),
		}))
		return
	}

	if err := utils.ValidateStruct(input); err != nil {
		middleware.HandleError(c, errors.ValidationError(map[string]string{
			"validation": err.Error(),
		}))
		return
	}

	article, err := h.articleService.UpdateArticle(uint(articleID), input, adminID.(uint))
	if err != nil {
		if err.Error() == "permission denied" {
			middleware.HandleError(c, errors.ErrForbidden)
		} else {
			middleware.HandleError(c, errors.WrapError(err, errors.ErrCodeInternal, 500))
		}
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
		middleware.HandleError(c, errors.ErrUnauthorized)
		return
	}

	id := c.Param("id")
	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		middleware.HandleError(c, errors.ErrBadRequest)
		return
	}

	if err := h.articleService.DeleteArticle(uint(articleID), adminID.(uint)); err != nil {
		if err.Error() == "permission denied" {
			middleware.HandleError(c, errors.ErrForbidden)
		} else {
			middleware.HandleError(c, errors.WrapError(err, errors.ErrCodeInternal, 500))
		}
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
		middleware.HandleError(c, errors.ErrUnauthorized)
		return
	}

	id := c.Param("id")
	articleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		middleware.HandleError(c, errors.ErrBadRequest)
		return
	}

	article, err := h.articleService.PublishArticle(uint(articleID), adminID.(uint))
	if err != nil {
		if err.Error() == "permission denied" {
			middleware.HandleError(c, errors.ErrForbidden)
		} else {
			middleware.HandleError(c, errors.WrapError(err, errors.ErrCodeInternal, 500))
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    article,
	})
}
