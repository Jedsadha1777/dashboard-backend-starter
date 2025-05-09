package services

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/utils"
	"errors"
	"strconv"
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
	// แปลง id เป็น uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	// ใช้ FindWithPreload เพื่อดึงข้อมูลพร้อม Admin
	return s.repo.FindWithPreload([]string{"Admin"}, uint(idUint))
}

// GetArticleBySlug retrieves an article by slug
func (s *ArticleService) GetArticleBySlug(slug string) (*models.Article, error) {
	// ใช้ FindOneWithPreload เพื่อค้นหาตาม slug พร้อม preload Admin
	return s.repo.FindOneWithPreload([]string{"Admin"}, "slug = ?", slug)
}

// CreateArticle creates a new article
func (s *ArticleService) CreateArticle(input *models.ArticleInput, adminID uint) (*models.Article, error) {
	// ตรวจสอบ slug ซ้ำ
	count, err := s.repo.Count("slug = ?", input.Slug)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		// If slug already exists, generate a unique one
		uniqueSlug, err := utils.EnsureUniqueSlug(db.DB, input.Slug, "articles", "slug")
		if err != nil {
			return nil, err
		}
		input.Slug = uniqueSlug
	}

	// แปลงวันที่ published
	var publishedAt *time.Time
	if input.PublishedAt != "" && input.Status == "published" {
		parsedTime, err := time.Parse(time.RFC3339, input.PublishedAt)
		if err != nil {
			return nil, errors.New("invalid published_at format. Use ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)")
		}
		publishedAt = &parsedTime
	} else if input.Status == "published" {
		now := time.Now()
		publishedAt = &now
	}

	// สร้าง article
	article := &models.Article{
		Title:       input.Title,
		Content:     input.Content,
		Slug:        input.Slug,
		Summary:     input.Summary,
		Status:      input.Status,
		PublishedAt: publishedAt,
		AdminID:     adminID,
	}

	// สร้าง article ในฐานข้อมูล
	err = db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Create(article)
	})

	if err != nil {
		return nil, err
	}

	// ดึงข้อมูลที่สมบูรณ์พร้อม Admin
	return s.repo.FindWithPreload([]string{"Admin"}, article.ID)
}

// GetArticles retrieves articles with pagination and search
func (s *ArticleService) GetArticles(params utils.PaginationParams) ([]models.Article, *utils.PaginationResult, error) {
	var articles []models.Article

	// เพิ่ม preload Admin
	params.Preloads = []string{"Admin"}

	// สร้าง query
	query := db.DB.Model(&models.Article{})

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
	// แปลง id เป็น uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	// ค้นหา article
	article, err := s.repo.FindByID(uint(idUint))
	if err != nil {
		return nil, err
	}

	// ตรวจสอบว่า admin เป็นเจ้าของบทความนี้
	if article.AdminID != adminID {
		return nil, errors.New("you don't have permission to update this article")
	}

	// ตรวจสอบ slug ซ้ำถ้ามีการเปลี่ยนแปลง
	if input.Slug != article.Slug {
		count, err := s.repo.Count("slug = ? AND id != ?", input.Slug, article.ID)
		if err != nil {
			return nil, err
		}

		if count > 0 {
			// If slug already exists, generate a unique one
			uniqueSlug, err := utils.EnsureUniqueSlug(db.DB, input.Slug, "articles", "slug", article.ID)
			if err != nil {
				return nil, err
			}
			input.Slug = uniqueSlug
		}
	}

	// แปลงวันที่เผยแพร่ถ้ามีการระบุ
	var publishedAt *time.Time
	if input.PublishedAt != "" && input.Status == "published" {
		parsedTime, err := time.Parse(time.RFC3339, input.PublishedAt)
		if err != nil {
			return nil, errors.New("invalid published_at format. Use ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)")
		}
		publishedAt = &parsedTime
	} else if input.Status == "published" && (article.PublishedAt == nil || article.Status != "published") {
		// ถ้าสถานะเปลี่ยนเป็น published แต่ไม่ได้ระบุวันที่ ให้ใช้เวลาปัจจุบัน
		now := time.Now()
		publishedAt = &now
	} else {
		// คงวันที่เผยแพร่เดิม
		publishedAt = article.PublishedAt
	}

	// อัปเดตข้อมูลบทความ
	article.Title = input.Title
	article.Content = input.Content
	article.Slug = input.Slug
	article.Summary = input.Summary
	article.Status = input.Status
	article.PublishedAt = publishedAt

	// อัปเดตด้วย transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Update(article)
	})

	if err != nil {
		return nil, err
	}

	// ดึงข้อมูลที่อัปเดตแล้วพร้อม Admin
	return s.repo.FindWithPreload([]string{"Admin"}, article.ID)
}

// DeleteArticle deletes an article
func (s *ArticleService) DeleteArticle(id string, adminID uint) error {
	// แปลง id เป็น uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errors.New("invalid ID format")
	}

	// ค้นหาบทความ
	article, err := s.repo.FindByID(uint(idUint))
	if err != nil {
		return err
	}

	// ตรวจสอบว่า admin เป็นเจ้าของบทความนี้
	if article.AdminID != adminID {
		return errors.New("you don't have permission to delete this article")
	}

	// ลบบทความด้วย transaction
	return db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Delete(article.ID)
	})
}

// PublishArticle sets an article to published status
func (s *ArticleService) PublishArticle(id string, adminID uint) (*models.Article, error) {
	// แปลง id เป็น uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	// ค้นหาบทความ
	article, err := s.repo.FindByID(uint(idUint))
	if err != nil {
		return nil, err
	}

	// ตรวจสอบว่า admin เป็นเจ้าของบทความนี้
	if article.AdminID != adminID {
		return nil, errors.New("you don't have permission to publish this article")
	}

	// อัปเดตสถานะบทความเป็น published
	now := time.Now()
	article.Status = "published"
	article.PublishedAt = &now

	// บันทึกการเปลี่ยนแปลงด้วย transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Update(article)
	})

	if err != nil {
		return nil, err
	}

	// ดึงข้อมูลที่อัปเดตแล้วพร้อม Admin
	return s.repo.FindWithPreload([]string{"Admin"}, article.ID)
}
