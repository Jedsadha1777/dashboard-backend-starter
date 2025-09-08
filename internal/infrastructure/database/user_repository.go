package database

import (
	"dashboard-starter/internal/domain/user/entity"
	"dashboard-starter/internal/domain/user/repository"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}

func (r *userRepository) List(offset, limit int, search string) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	query := r.db.Model(&entity.User{})

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&users).Error
	return users, total, err
}

func (r *userRepository) Count(conditions ...interface{}) (int64, error) {
	var count int64
	query := r.db.Model(&entity.User{})

	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}

	err := query.Count(&count).Error
	return count, err
}
