package database

import (
	"dashboard-starter/internal/domain/auth/entity"
	"dashboard-starter/internal/domain/auth/repository"

	"gorm.io/gorm"
)

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) repository.AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) GetByID(id uint) (*entity.Admin, error) {
	var admin entity.Admin
	err := r.db.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) GetByEmail(email string) (*entity.Admin, error) {
	var admin entity.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) Create(admin *entity.Admin) error {
	return r.db.Create(admin).Error
}

func (r *adminRepository) Update(admin *entity.Admin) error {
	return r.db.Save(admin).Error
}

func (r *adminRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Admin{}, id).Error
}

func (r *adminRepository) IncrementTokenVersion(id uint) error {
	return r.db.Model(&entity.Admin{}).
		Where("id = ?", id).
		Update("token_version", gorm.Expr("token_version + 1")).Error
}
