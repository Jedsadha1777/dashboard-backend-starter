package db

import (
	"dashboard-starter/config"
	"dashboard-starter/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	DB.AutoMigrate(&models.User{})
}
