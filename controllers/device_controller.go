package controllers

import (
	"crypto/rand"
	"dashboard-starter/db"
	"dashboard-starter/models"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ฟังก์ชันช่วยสร้าง API key แบบสุ่ม
func generateRandomAPIKey(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateDevice สร้างอุปกรณ์ใหม่ในระบบ (เฉพาะ admin เท่านั้น)
func CreateDevice(c *gin.Context) {
	// ตรวจสอบว่าเป็น admin โดย AdminRequired middleware แล้ว
	adminID, _ := c.Get("admin_id")

	var input struct {
		DeviceID string `json:"device_id" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ข้อมูลไม่ถูกต้อง: " + err.Error(),
		})
		return
	}

	// ตรวจสอบว่า device_id ซ้ำหรือไม่
	var count int64
	if err := db.DB.Model(&models.Device{}).Where("device_id = ?", input.DeviceID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "ไม่สามารถตรวจสอบข้อมูลได้: " + err.Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Device ID นี้มีในระบบแล้ว กรุณาใช้ ID อื่น",
		})
		return
	}

	// สร้าง API key แบบสุ่ม
	apiKey, err := generateRandomAPIKey(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "ไม่สามารถสร้าง API key ได้: " + err.Error(),
		})
		return
	}

	// สร้างอุปกรณ์ใหม่
	device := models.Device{
		DeviceID:     input.DeviceID,
		Name:         input.Name,
		ApiKey:       apiKey,
		TokenVersion: 1,
		Status:       "inactive",
		LastSeen:     time.Now(),
	}

	if err := db.DB.Create(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "ไม่สามารถสร้างอุปกรณ์ได้: " + err.Error(),
		})
		return
	}

	// ส่งข้อมูลกลับพร้อม API key (แสดงครั้งเดียวตอนสร้าง)
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data: gin.H{
			"id":         device.ID,
			"device_id":  device.DeviceID,
			"name":       device.Name,
			"api_key":    apiKey, // แสดง API key ให้ admin เห็นครั้งเดียว
			"status":     device.Status,
			"created_by": adminID,
		},
	})
}

// ListDevices แสดงรายการอุปกรณ์ทั้งหมด
func ListDevices(c *gin.Context) {
	// ตรวจสอบว่าเป็น admin โดย AdminRequired middleware แล้ว

	// ดึงค่า pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// ค้นหาข้อมูล
	var devices []models.Device
	var total int64

	query := db.DB.Model(&models.Device{})

	// ตรวจสอบหากมีการค้นหา
	search := c.Query("search")
	if search != "" {
		query = query.Where("device_id LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// นับจำนวนทั้งหมด
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "ไม่สามารถนับจำนวนอุปกรณ์ได้: " + err.Error(),
		})
		return
	}

	// ดึงข้อมูลตาม limit และ offset
	if err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "ไม่สามารถดึงข้อมูลอุปกรณ์ได้: " + err.Error(),
		})
		return
	}

	// คำนวณจำนวนหน้าทั้งหมด
	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    devices,
		Meta: gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// GetDevice ดึงข้อมูลอุปกรณ์ตาม ID
func GetDevice(c *gin.Context) {
	// ตรวจสอบว่าเป็น admin โดย AdminRequired middleware แล้ว

	id := c.Param("id")

	var device models.Device
	if err := db.DB.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "ไม่พบอุปกรณ์",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    device,
	})
}

// UpdateDevice อัปเดตข้อมูลอุปกรณ์
func UpdateDevice(c *gin.Context) {
	// ตรวจสอบว่าเป็น admin โดย AdminRequired middleware แล้ว

	id := c.Param("id")

	var device models.Device
	if err := db.DB.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "ไม่พบอุปกรณ์",
		})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ข้อมูลไม่ถูกต้อง: " + err.Error(),
		})
		return
	}

	// อัปเดตข้อมูลที่อนุญาตให้แก้ไขได้
	err := db.Transaction(func(tx *gorm.DB) error {
		device.Name = input.Name
		return tx.Save(&device).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "ไม่สามารถอัปเดตอุปกรณ์ได้: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    device,
	})
}

// ResetDeviceApiKey รีเซ็ต API key ของอุปกรณ์
func ResetDeviceApiKey(c *gin.Context) {
	// ตรวจสอบว่าเป็น admin โดย AdminRequired middleware แล้ว

	id := c.Param("id")

	var device models.Device
	if err := db.DB.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "ไม่พบอุปกรณ์",
		})
		return
	}

	// สร้าง API key ใหม่แบบสุ่ม
	newApiKey, err := generateRandomAPIKey(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "ไม่สามารถสร้าง API key ได้: " + err.Error(),
		})
		return
	}

	// อัปเดตค่า
	err = db.Transaction(func(tx *gorm.DB) error {
		device.ApiKey = newApiKey
		device.TokenVersion += 1 // เพิ่มเวอร์ชันเพื่อทำให้ token เก่าหมดอายุ
		return tx.Save(&device).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "ไม่สามารถรีเซ็ต API key ได้: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"id":        device.ID,
			"device_id": device.DeviceID,
			"name":      device.Name,
			"api_key":   newApiKey, // แสดง API key ใหม่
			"message":   "API key ถูกรีเซ็ตเรียบร้อยแล้ว ต้องลงทะเบียนอุปกรณ์ใหม่",
		},
	})
}

// DeleteDevice ลบอุปกรณ์ออกจากระบบ
func DeleteDevice(c *gin.Context) {
	// ตรวจสอบว่าเป็น admin โดย AdminRequired middleware แล้ว

	id := c.Param("id")

	var device models.Device
	if err := db.DB.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "ไม่พบอุปกรณ์",
		})
		return
	}

	// ลบอุปกรณ์
	err := db.Transaction(func(tx *gorm.DB) error {
		// ลบ refresh token ที่เกี่ยวข้องทั้งหมด
		if err := tx.Where("user_id = ? AND user_type = ?", device.ID, "device").Delete(&models.RefreshToken{}).Error; err != nil {
			return err
		}

		// ลบอุปกรณ์
		return tx.Delete(&device).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "ไม่สามารถลบอุปกรณ์ได้: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"message": "ลบอุปกรณ์เรียบร้อยแล้ว",
		},
	})
}
