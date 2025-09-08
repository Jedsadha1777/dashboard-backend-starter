package database

import (
	"dashboard-starter/internal/domain/article/entity"
	"dashboard-starter/internal/domain/article/repository"

	"gorm.io/gorm"
)

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) repository.ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) GetByID(id uint) (*entity.Article, error) {
	var article entity.Article
	err := r.db.First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) GetBySlug(slug string) (*entity.Article, error) {
	var article entity.Article
	err := r.db.Where("slug = ?", slug).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) Create(article *entity.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepository) Update(article *entity.Article) error {
	return r.db.Save(article).Error
}

func (r *articleRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Article{}, id).Error
}

func (r *articleRepository) List(offset, limit int, search, status string) ([]*entity.Article, int64, error) {
	var articles []*entity.Article
	var total int64

	query := r.db.Model(&entity.Article{})

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("title LIKE ? OR content LIKE ? OR slug LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Apply status filter
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&articles).Error
	return articles, total, err
}

func (r *articleRepository) Count(conditions ...interface{}) (int64, error) {
	var count int64
	query := r.db.Model(&entity.Article{})

	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}

	err := query.Count(&count).Error
	return count, err
}
