package db

import (
	"dashboard-starter/models"
	"time"

	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAdmin creates a default admin user if none exists
func SeedAdmin() error {
	// Check if any admin exists
	var count int64
	if err := DB.Model(&models.Admin{}).Count(&count).Error; err != nil {
		return err
	}

	// If no admin exists, create one
	if count == 0 {
		// Generate password hash with higher cost for security
		hashed, err := bcrypt.GenerateFromPassword([]byte("Admin@123!"), bcrypt.DefaultCost+1)
		if err != nil {
			return err
		}

		// Create admin in a transaction
		err = Transaction(func(tx *gorm.DB) error {
			admin := models.Admin{
				Email:        "admin@example.com",
				Password:     string(hashed),
				TokenVersion: 1,
				LastLogin:    time.Now(),
			}

			if err := tx.Create(&admin).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Printf("Failed to seed admin: %v", err)
			return err
		}

		log.Println("Successfully created default admin account")
		log.Println("Email: admin@example.com")
		log.Println("Password: Admin@123!")
		log.Println("IMPORTANT: Please change this password after first login")
	}

	return nil
}

// SeedTestData seeds the database with test data (for development only)
func SeedTestData() error {
	// Create test users
	testUsers := []models.User{
		{Name: "John Doe", Email: "john.doe@example.com"},
		{Name: "Jane Smith", Email: "jane.smith@example.com"},
		{Name: "Bob Johnson", Email: "bob.johnson@example.com"},
		{Name: "Alice Williams", Email: "alice.williams@example.com"},
		{Name: "Charlie Brown", Email: "charlie.brown@example.com"},
	}

	// Check if test data already exists
	var count int64
	if err := DB.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	// Only seed if no users exist
	if count == 0 {
		log.Println("Seeding test data...")

		// Use transaction for data consistency
		err := Transaction(func(tx *gorm.DB) error {
			for _, user := range testUsers {
				if err := tx.Create(&user).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			log.Printf("Failed to seed test data: %v", err)
			return err
		}

		log.Println("Successfully seeded test data")
	}

	return nil
}
