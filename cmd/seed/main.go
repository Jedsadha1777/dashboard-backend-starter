package main

import (
	"dashboard-starter/config"
	"dashboard-starter/db"
	"dashboard-starter/utils"
	"log"
)

func main() {
	log.Println("Starting database seeder...")

	// Load configuration
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize JWT (needed for some operations)
	if err := utils.InitJWT(); err != nil {
		log.Fatalf("Failed to initialize JWT: %v", err)
	}

	// Initialize database
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Ensure database is closed when application exits
	sqlDB, _ := db.DB.DB()
	defer sqlDB.Close()

	// Initialize validators
	utils.InitValidator()

	// Run seeders
	log.Println("Running admin seeder...")
	if err := db.SeedAdmin(); err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}

	log.Println("Seeding completed successfully")
}
