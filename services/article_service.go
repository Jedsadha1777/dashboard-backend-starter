package services

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ArticleService struct {
	repo *db.GormRepository[models.Article]
}

func NewArticleService() *ArticleService {
	return &ArticleService{
		repo: db.NewRepository[models.Article](),
	}
}

func (s *ArticleService) GetByID(id string) (*models.Article, error) {
	// ในกรณีนี้ต้องใช้ query แบบพิเศษเพื่อ preload Admin
	var article models.Article
	err := db.DB.Preload("Admin").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetArticleBySlug retrieves an article by slug
func (s *ArticleService) GetArticleBySlug(slug string) (*models.Article, error) {
	var article models.Article
	if err := db.DB.Preload("Admin").Where("slug = ?", slug).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// CreateArticle creates a new article
func (s *ArticleService) CreateArticle(input *models.ArticleInput, adminID uint) (*models.Article, error) {
	// Check for duplicate slug
	var count int64
	if err := db.DB.Model(&models.Article{}).Where("slug = ?", input.Slug).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("slug already exists")
	}

	// Parse published date if provided
	var publishedAt *time.Time
	if input.PublishedAt != "" && input.Status == "published" {
		parsedTime, err := time.Parse(time.RFC3339, input.PublishedAt)
		if err != nil {
			return nil, errors.New("invalid published_at format. Use ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)")
		}
		publishedAt = &parsedTime
	} else if input.Status == "published" {
		// If status is published but no date provided, use current time
		now := time.Now()
		publishedAt = &now
	}

	// Create article inside a transaction
	var article models.Article
	err := db.Transaction(func(tx *gorm.DB) error {
		article = models.Article{
			Title:       input.Title,
			Content:     input.Content,
			Slug:        input.Slug,
			Summary:     input.Summary,
			Status:      input.Status,
			PublishedAt: publishedAt,
			AdminID:     adminID,
		}

		if err := tx.Create(&article).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &article, nil
}

// GetArticles retrieves articles with pagination and search
func (s *ArticleService) GetArticles(params utils.PaginationParams) ([]models.Article, *utils.PaginationResult, error) {
	var articles []models.Article

	// สร้าง query
	query := db.DB.Model(&models.Article{}).Preload("Admin")

	// ใช้ search ถ้ามี
	if params.Search != "" {
		query = utils.ApplySearch(query, params.Search, "title", "content", "slug")
	}

	// กรองตาม status ถ้ามี
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// ใช้ pagination
	result, err := utils.ApplyPagination(query, params, &articles)
	if err != nil {
		return nil, nil, err
	}

	return articles, result, nil
}

// UpdateArticle updates an existing article
func (s *ArticleService) UpdateArticle(id string, input *models.ArticleInput, adminID uint) (*models.Article, error) {
	// Find the article
	var article models.Article
	if err := db.DB.Where("id = ?", id).First(&article).Error; err != nil {
		return nil, err
	}

	// Check if admin owns this article
	if article.AdminID != adminID {
		return nil, errors.New("you don't have permission to update this article")
	}

	// Check for duplicate slug if changing
	if input.Slug != article.Slug {
		var count int64
		if err := db.DB.Model(&models.Article{}).Where("slug = ? AND id != ?", input.Slug, id).Count(&count).Error; err != nil {
			return nil, err
		}

		if count > 0 {
			return nil, errors.New("slug already exists")
		}
	}

	// Parse published date if provided
	var publishedAt *time.Time
	if input.PublishedAt != "" && input.Status == "published" {
		parsedTime, err := time.Parse(time.RFC3339, input.PublishedAt)
		if err != nil {
			return nil, errors.New("invalid published_at format. Use ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)")
		}
		publishedAt = &parsedTime
	} else if input.Status == "published" && (article.PublishedAt == nil || article.Status != "published") {
		// If status is being changed to published but no date provided, use current time
		now := time.Now()
		publishedAt = &now
	} else {
		// Keep the existing published date
		publishedAt = article.PublishedAt
	}

	// Update article inside a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		article.Title = input.Title
		article.Content = input.Content
		article.Slug = input.Slug
		article.Summary = input.Summary
		article.Status = input.Status
		article.PublishedAt = publishedAt

		if err := tx.Save(&article).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload article with admin data
	if err := db.DB.Preload("Admin").Where("id = ?", id).First(&article).Error; err != nil {
		return nil, err
	}

	return &article, nil
}

// DeleteArticle deletes an article
func (s *ArticleService) DeleteArticle(id string, adminID uint) error {
	// Check if article exists and belongs to admin
	var article models.Article
	if err := db.DB.Where("id = ?", id).First(&article).Error; err != nil {
		return err
	}

	// Check if admin owns this article
	if article.AdminID != adminID {
		return errors.New("you don't have permission to delete this article")
	}

	// Delete article inside a transaction
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&article).Error
	})
}

// PublishArticle sets an article to published status
func (s *ArticleService) PublishArticle(id string, adminID uint) (*models.Article, error) {
	// Find the article
	var article models.Article
	if err := db.DB.Where("id = ?", id).First(&article).Error; err != nil {
		return nil, err
	}

	// Check if admin owns this article
	if article.AdminID != adminID {
		return nil, errors.New("you don't have permission to publish this article")
	}

	// Update article status inside a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		article.Status = "published"
		article.PublishedAt = &now

		if err := tx.Save(&article).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &article, nil
}
