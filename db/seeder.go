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

	// After creating/verifying admin, seed test data
	err := SeedTestData(adminID)
	if err != nil {
		return err
	}

	// Seed test regular users (self-registered)
	return SeedTestUsers()
}

// SeedTestData seeds the database with test data (for development only)
func SeedTestData(adminID uint) error {
	// Create test users with AdminID (admin-created users)
	testUsers := []models.User{
		{Name: "John Doe", Email: "john.doe@example.com", AdminID: adminID},
		{Name: "Jane Smith", Email: "jane.smith@example.com", AdminID: adminID},
		{Name: "Bob Johnson", Email: "bob.johnson@example.com", AdminID: adminID},
		{Name: "Alice Williams", Email: "alice.williams@example.com", AdminID: adminID},
		{Name: "Charlie Brown", Email: "charlie.brown@example.com", AdminID: adminID},
	}

	// Check if test data already exists
	var count int64
	if err := DB.Model(&models.User{}).Where("admin_id = ?", adminID).Count(&count).Error; err != nil {
		return err
	}

	// Only seed if no admin-created users exist
	if count == 0 {
		log.Println("Seeding admin-created test users...")

		// Generate default password
		defaultPassword, err := bcrypt.GenerateFromPassword([]byte("User@123!"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Use transaction for data consistency
		err = Transaction(func(tx *gorm.DB) error {
			for _, user := range testUsers {
				// Set password and token version
				user.Password = string(defaultPassword)
				user.TokenVersion = 1
				user.LastLogin = time.Now()

				if err := tx.Create(&user).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			log.Printf("Failed to seed admin-created test users: %v", err)
			return err
		}

		log.Println("Successfully seeded admin-created test users")
		log.Println("Default password for all admin-created test users: User@123!")
	}

	return nil
}

// SeedTestUsers creates self-registered test users
func SeedTestUsers() error {
	// Create self-registered users (no AdminID)
	selfRegisteredUsers := []models.User{
		{Name: "Sam Wilson", Email: "sam.wilson@example.com"},
		{Name: "Maria Rodriguez", Email: "maria.rodriguez@example.com"},
		{Name: "David Kim", Email: "david.kim@example.com"},
	}

	// Check if self-registered test users already exist
	var count int64
	if err := DB.Model(&models.User{}).Where("admin_id = ? OR admin_id IS NULL", 0).Count(&count).Error; err != nil {
		return err
	}

	// Only seed if no self-registered users exist
	if count == 0 {
		log.Println("Seeding self-registered test users...")

		// Generate default password
		defaultPassword, err := bcrypt.GenerateFromPassword([]byte("SelfUser@123!"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Use transaction for data consistency
		err = Transaction(func(tx *gorm.DB) error {
			for _, user := range selfRegisteredUsers {
				// Set password and token version
				user.Password = string(defaultPassword)
				user.TokenVersion = 1
				user.LastLogin = time.Now()
				user.AdminID = 0 // Explicitly set to 0 to indicate self-registration

				if err := tx.Create(&user).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			log.Printf("Failed to seed self-registered test users: %v", err)
			return err
		}

		log.Println("Successfully seeded self-registered test users")
		log.Println("Default password for all self-registered test users: SelfUser@123!")
	}

	return nil
}
