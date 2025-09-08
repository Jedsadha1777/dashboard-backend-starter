package db

import (
	articleEntity "dashboard-starter/internal/domain/article/entity"
	authEntity "dashboard-starter/internal/domain/auth/entity"
	deviceEntity "dashboard-starter/internal/domain/device/entity"
	userEntity "dashboard-starter/internal/domain/user/entity"

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
		&userEntity.User{},
		&authEntity.Admin{},
		&authEntity.RefreshToken{},
		&deviceEntity.Device{},
		&articleEntity.Article{},
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
