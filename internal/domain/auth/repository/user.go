package repository

import "dashboard-starter/internal/domain/user/entity"

type UserRepository interface {
	GetByID(id uint) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uint) error
	List(offset, limit int, search string) ([]*entity.User, int64, error)
	Count(conditions ...interface{}) (int64, error)
}
