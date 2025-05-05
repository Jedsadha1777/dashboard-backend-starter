package main

import (
	"dashboard-starter/config"
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/routes"
	"dashboard-starter/utils"
	"time"
)

func init() {
	time.Local = time.UTC
}

func main() {
	config.Init()
	utils.InitJWT()

	db.Init()
	db.DB.AutoMigrate(&models.Admin{})
	db.SeedAdmin()
	r := routes.SetupRouter()
	r.Run(":8080")
}
