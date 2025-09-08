package services

import (
	"dashboard-starter/internal/application/dto"
	"dashboard-starter/internal/domain/article/entity"
	"dashboard-starter/internal/domain/article/repository"
	"dashboard-starter/utils"
	"errors"
	"time"
)

type ArticleApplicationService struct {
	articleRepo repository.ArticleRepository
}

func NewArticleApplicationService(articleRepo repository.ArticleRepository) *ArticleApplicationService {
	return &ArticleApplicationService{
		articleRepo: articleRepo,
	}
}

func (s *ArticleApplicationService) CreateArticle(input dto.ArticleInput, adminID uint) (*entity.Article, error) {
	// Generate slug if not provided
	if input.Slug == "" {
		baseSlug, err := utils.GenerateSlug(input.Title, 100)
		if err != nil {
			return nil, err
		}
		// TODO: Use EnsureUniqueSlug with article repository
		input.Slug = baseSlug
	}

	// Set default status
	if input.Status == "" {
		input.Status = "draft"
	}

	// Parse published date
	var publishedAt *time.Time
	if input.PublishedAt != "" && input.Status == "published" {
		parsedTime, err := time.Parse(time.RFC3339, input.PublishedAt)
		if err != nil {
			return nil, errors.New("invalid published_at format")
		}
		publishedAt = &parsedTime
	} else if input.Status == "published" {
		now := time.Now()
		publishedAt = &now
	}

	article := &entity.Article{
		Title:       input.Title,
		Content:     input.Content,
		Slug:        input.Slug,
		Summary:     input.Summary,
		Status:      input.Status,
		PublishedAt: publishedAt,
		AdminID:     adminID,
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *ArticleApplicationService) GetArticle(id uint) (*entity.Article, error) {
	return s.articleRepo.GetByID(id)
}

func (s *ArticleApplicationService) GetArticleBySlug(slug string) (*entity.Article, error) {
	return s.articleRepo.GetBySlug(slug)
}

func (s *ArticleApplicationService) UpdateArticle(id uint, input dto.ArticleInput, adminID uint) (*entity.Article, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if article.AdminID != adminID {
		return nil, errors.New("permission denied")
	}

	// Update fields
	article.Title = input.Title
	article.Content = input.Content
	article.Slug = input.Slug
	article.Summary = input.Summary
	article.Status = input.Status

	// Handle published date
	if input.Status == "published" && article.PublishedAt == nil {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *ArticleApplicationService) DeleteArticle(id uint, adminID uint) error {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check ownership
	if article.AdminID != adminID {
		return errors.New("permission denied")
	}

	return s.articleRepo.Delete(id)
}

func (s *ArticleApplicationService) ListArticles(page, limit int, search, status string) ([]*entity.Article, int64, error) {
	offset := (page - 1) * limit
	return s.articleRepo.List(offset, limit, search, status)
}

func (s *ArticleApplicationService) PublishArticle(id uint, adminID uint) (*entity.Article, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if article.AdminID != adminID {
		return nil, errors.New("permission denied")
	}

	article.Status = "published"
	now := time.Now()
	article.PublishedAt = &now

	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}

	return article, nil
}
