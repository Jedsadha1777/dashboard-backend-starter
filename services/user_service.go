package services

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

func GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers retrieves users with pagination and search
func GetUsers(search string, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var count int64

	query := db.DB.Model(&models.User{})

	// Apply search filter if provided
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", searchTerm, searchTerm)
	}

	// Count total matching records (before pagination)
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset based on page and limit
	offset := (page - 1) * limit

	// Execute query with pagination and sorting
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

// CreateUser creates a new user
func CreateUser(input *models.UserInput) (*models.User, error) {
	// Check for duplicate email
	var count int64
	if err := db.DB.Model(&models.User{}).Where("email = ?", input.Email).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("email already exists")
	}

	// Create user inside a transaction
	var user models.User
	err := db.Transaction(func(tx *gorm.DB) error {
		user = models.User{
			Name:  input.Name,
			Email: input.Email,
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user
func UpdateUser(id string, input *models.UserInput) (*models.User, error) {
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
func DeleteUser(id string) error {
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
