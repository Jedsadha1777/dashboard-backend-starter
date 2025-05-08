package db

import (
	"dashboard-starter/models"
	"log"
)

var (
	// modelRegistry เก็บรายการโมเดลทั้งหมดที่ต้องทำ migration
	modelRegistry []interface{}
)

// RegisterAllModels ลงทะเบียนโมเดลทั้งหมดที่ต้องทำ migration
func RegisterAllModels() {
	// ล้างรายการเดิม (ถ้ามี) และเพิ่มโมเดลทั้งหมด
	modelRegistry = []interface{}{
		&models.User{},
		&models.Admin{},
		&models.RefreshToken{},
		&models.Device{},
		&models.Article{},
		// เพิ่มโมเดลใหม่ตรงนี้:
		// &models.Product{},
		// &models.Category{},
		// ...
	}

	log.Printf("Registered %d models for migrations", len(modelRegistry))
}

// GetAllModels คืนค่ารายการโมเดลทั้งหมดสำหรับการทำ migrations
func GetAllModels() []interface{} {
	return modelRegistry
}
