package services

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/utils"
	"errors"
	"strconv"

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
	// Convert id to uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	return s.repo.FindByID(uint(idUint))
}

// GetUsers retrieves users with pagination and search
func (s *UserService) GetUsers(params utils.PaginationParams) ([]models.User, *utils.PaginationResult, error) {
	var users []models.User

	// Create query
	query := db.DB.Model(&models.User{})

	// Apply search if provided
	if params.Search != "" {
		query = utils.ApplySearch(query, params.Search, "name", "email")
	}

	// Apply pagination
	result, err := utils.ApplyPagination(query, params, &users)
	if err != nil {
		return nil, nil, err
	}

	return users, result, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(input *models.UserInput, adminID uint) (*models.User, error) {
	// Check for duplicate email
	count, err := s.repo.Count("email = ?", input.Email)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("email already exists")
	}

	// Create new user
	user := &models.User{
		Name:    input.Name,
		Email:   input.Email,
		AdminID: adminID, // Set the admin ID who created this user
	}

	// Save to database
	err = db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Create(user)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id string, input *models.UserInput, adminID uint) (*models.User, error) {
	// Convert id to uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	// Find user
	user, err := s.repo.FindByID(uint(idUint))
	if err != nil {
		return nil, err
	}

	// Check if admin has permission to update this user
	if user.AdminID != adminID {
		return nil, errors.New("you don't have permission to update this user")
	}

	// Check for duplicate email (if changed)
	if input.Email != user.Email {
		count, err := s.repo.Count("email = ? AND id != ?", input.Email, user.ID)
		if err != nil {
			return nil, err
		}

		if count > 0 {
			return nil, errors.New("email already exists")
		}
	}

	// Update user data
	user.Name = input.Name
	user.Email = input.Email

	// Save changes
	err = db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Update(user)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string, adminID uint) error {
	// Convert id to uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errors.New("invalid ID format")
	}

	// Check if user exists
	user, err := s.repo.FindByID(uint(idUint))
	if err != nil {
		return err
	}

	// Check if admin has permission to delete this user
	if user.AdminID != adminID {
		return errors.New("you don't have permission to delete this user")
	}

	// Delete user
	return db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Delete(uint(idUint))
	})
}
