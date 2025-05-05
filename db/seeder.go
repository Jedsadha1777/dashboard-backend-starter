package db

import (
	"dashboard-starter/models"

	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	var count int64
	DB.Model(&models.Admin{}).Count(&count)
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.Admin{
			Email:        "admin@mail.com",
			Password:     string(hashed),
			TokenVersion: 1,
		}
		if err := DB.Create(&admin).Error; err != nil {
			log.Println("Failed to seed admin:", err)
		} else {
			log.Println("Seeded default admin.")
		}
	}
}
