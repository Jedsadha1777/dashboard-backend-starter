package main

import (
	"dashboard-starter/config"
	"dashboard-starter/db"
	"dashboard-starter/models"
	"log"
)

func main() {
	log.Println("Starting user migration for AdminID...")

	// Load configuration
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Ensure database is closed when application exits
	sqlDB, _ := db.DB.DB()
	defer sqlDB.Close()

	// Find the first admin to assign to existing users
	var admin models.Admin
	if err := db.DB.First(&admin).Error; err != nil {
		log.Fatalf("Failed to find an admin: %v", err)
	}

	// Update all users without AdminID to use the first admin's ID
	result := db.DB.Model(&models.User{}).
		Where("admin_id = 0 OR admin_id IS NULL").
		Update("admin_id", admin.ID)

	if result.Error != nil {
		log.Fatalf("Failed to update users: %v", result.Error)
	}

	log.Printf("Successfully updated %d users with AdminID = %d", result.RowsAffected, admin.ID)
}
