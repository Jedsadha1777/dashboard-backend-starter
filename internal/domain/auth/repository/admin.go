package repository

import "dashboard-starter/internal/domain/auth/entity"

type AdminRepository interface {
	GetByID(id uint) (*entity.Admin, error)
	GetByEmail(email string) (*entity.Admin, error)
	Create(admin *entity.Admin) error
	Update(admin *entity.Admin) error
	Delete(id uint) error
	IncrementTokenVersion(id uint) error
}
