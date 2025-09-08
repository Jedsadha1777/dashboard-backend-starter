package repository

import "dashboard-starter/internal/domain/article/entity"

type ArticleRepository interface {
	GetByID(id uint) (*entity.Article, error)
	GetBySlug(slug string) (*entity.Article, error)
	Create(article *entity.Article) error
	Update(article *entity.Article) error
	Delete(id uint) error
	List(offset, limit int, search, status string) ([]*entity.Article, int64, error)
	Count(conditions ...interface{}) (int64, error)
}
