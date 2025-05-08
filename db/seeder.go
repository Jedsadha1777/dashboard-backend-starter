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

	var adminID uint = 1 // Default admin ID to assign to seeded users

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

			// Store the new admin's ID for use with test users
			adminID = admin.ID
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
	} else {
		// If admin already exists, get the ID of the first admin
		var admin models.Admin
		if err := DB.First(&admin).Error; err != nil {
			log.Printf("Failed to get existing admin ID: %v", err)
			// Continue with default admin ID
		} else {
			adminID = admin.ID
		}
	}

	// After creating/verifying admin, seed test users
	return SeedTestData(adminID)
}

// SeedTestData seeds the database with test data (for development only)
func SeedTestData(adminID uint) error {
	// Create test users with AdminID
	testUsers := []models.User{
		{Name: "John Doe", Email: "john.doe@example.com", AdminID: adminID},
		{Name: "Jane Smith", Email: "jane.smith@example.com", AdminID: adminID},
		{Name: "Bob Johnson", Email: "bob.johnson@example.com", AdminID: adminID},
		{Name: "Alice Williams", Email: "alice.williams@example.com", AdminID: adminID},
		{Name: "Charlie Brown", Email: "charlie.brown@example.com", AdminID: adminID},
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
