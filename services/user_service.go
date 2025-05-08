package services

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/utils"
	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	repo *db.GormRepository[models.User]
}

func NewUserService() *UserService {
	return &UserService{
		repo: db.NewRepository[models.User](),
	}
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	return s.repo.FindByID(id)
}

// GetUsers retrieves users with pagination and search
func (s *UserService) GetUsers(params utils.PaginationParams) ([]models.User, *utils.PaginationResult, error) {
	var users []models.User

	// สร้าง query
	query := db.DB.Model(&models.User{})

	// ใช้ search ถ้ามี
	if params.Search != "" {
		query = utils.ApplySearch(query, params.Search, "name", "email")
	}

	// ใช้ pagination
	result, err := utils.ApplyPagination(query, params, &users)
	if err != nil {
		return nil, nil, err
	}

	return users, result, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(input *models.UserInput) (*models.User, error) {
	// ตรวจสอบ email ซ้ำ
	count, err := s.repo.Count("email = ?", input.Email)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("email already exists")
	}

	// สร้าง user ใหม่
	user := models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Create(&user)
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id string, input *models.UserInput) (*models.User, error) {
	// Find the user
	var user models.User
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	// Check for duplicate email (if email is being changed)
	if input.Email != user.Email {
		var count int64
		if err := db.DB.Model(&models.User{}).Where("email = ? AND id != ?", input.Email, id).Count(&count).Error; err != nil {
			return nil, err
		}

		if count > 0 {
			return nil, errors.New("email already exists")
		}
	}

	// Update user inside a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		user.Name = input.Name
		user.Email = input.Email

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string) error {
	// Check if user exists
	var user models.User
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	// Delete user inside a transaction
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&user).Error
	})
}
