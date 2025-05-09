package services

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/utils"
	"errors"
	"strconv"

	"golang.org/x/crypto/bcrypt"
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

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id string) (*models.User, error) {
	// Convert id to uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	return s.repo.FindByID(uint(idUint))
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.FindOne("email = ?", email)
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

// CreateUser creates a new user (by admin)
func (s *UserService) CreateUser(input *models.UserInput, adminID uint) (*models.User, error) {
	// Check for duplicate email
	count, err := s.repo.Count("email = ?", input.Email)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("email already exists")
	}

	// Generate a default password
	defaultPassword := utils.GenerateRandomPassword(12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &models.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     string(hashedPassword),
		TokenVersion: 1,
		AdminID:      adminID, // Set the admin ID who created this user
	}

	// Save to database
	err = db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Create(user)
	})

	if err != nil {
		return nil, err
	}

	// Attach the temporary password to return to admin
	user.Password = defaultPassword

	return user, nil
}

// RegisterUser registers a new user (self-registration)
func (s *UserService) RegisterUser(input *models.UserRegistrationInput) (*models.User, error) {
	// Check for duplicate email
	count, err := s.repo.Count("email = ?", input.Email)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &models.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     string(hashedPassword),
		TokenVersion: 1,
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

	// If user was created by an admin, check if the current admin has permission
	if user.AdminID != 0 && user.AdminID != adminID {
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

// UpdateUserProfile updates a user's own profile
func (s *UserService) UpdateUserProfile(userID uint, input *models.UserInput) (*models.User, error) {
	// Find user
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
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

// ChangeUserPassword changes a user's password
func (s *UserService) ChangeUserPassword(userID uint, currentPassword, newPassword string) error {
	// Find user
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password and token version
	return db.Transaction(func(tx *gorm.DB) error {
		user.Password = string(hashedPassword)
		user.TokenVersion += 1 // Invalidate existing tokens
		return s.repo.Update(user)
	})
}

// ResetUserPassword resets a user's password (admin only)
func (s *UserService) ResetUserPassword(id string, adminID uint) (string, error) {
	// Convert id to uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return "", errors.New("invalid ID format")
	}

	// Find user
	user, err := s.repo.FindByID(uint(idUint))
	if err != nil {
		return "", err
	}

	// If user was created by an admin, check if the current admin has permission
	if user.AdminID != 0 && user.AdminID != adminID {
		return "", errors.New("you don't have permission to reset this user's password")
	}

	// Generate new password
	newPassword := utils.GenerateRandomPassword(12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Update password and token version
	err = db.Transaction(func(tx *gorm.DB) error {
		user.Password = string(hashedPassword)
		user.TokenVersion += 1 // Invalidate existing tokens
		return s.repo.Update(user)
	})

	if err != nil {
		return "", err
	}

	return newPassword, nil
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

	// If user was created by an admin, check if the current admin has permission
	if user.AdminID != 0 && user.AdminID != adminID {
		return errors.New("you don't have permission to delete this user")
	}

	// Delete user and their refresh tokens in a transaction
	return db.Transaction(func(tx *gorm.DB) error {
		// First revoke all refresh tokens
		if err := tx.Where("user_id = ? AND user_type = ?", user.ID, "user").Delete(&models.RefreshToken{}).Error; err != nil {
			return err
		}

		// Then delete the user
		return s.repo.Delete(uint(idUint))
	})
}
