// utils/pagination.go
package utils

import (
	"strings"

	"gorm.io/gorm"
)

// PaginationParams พารามิเตอร์สำหรับ pagination
type PaginationParams struct {
	Page     int      `json:"page" form:"page"`
	Limit    int      `json:"limit" form:"limit"`
	OrderBy  string   `json:"order_by" form:"order_by"`
	Search   string   `json:"search" form:"search"`
	Status   string   `json:"status" form:"status"`
	Preloads []string `json:"-" form:"-"`
}

// PaginationResult ผลลัพธ์ของ pagination
type PaginationResult struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}

// NewPaginationParams สร้าง pagination params ใหม่ด้วยค่าเริ่มต้น
func NewPaginationParams() PaginationParams {
	return PaginationParams{
		Page:     1,
		Limit:    10,
		OrderBy:  "created_at desc",
		Preloads: []string{},
	}
}

// Normalize ปรับค่า pagination ให้อยู่ในช่วงที่ถูกต้อง
func (p *PaginationParams) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 || p.Limit > 100 {
		p.Limit = 10
	}
	if p.OrderBy == "" {
		p.OrderBy = "created_at desc"
	}
}

// GetOffset คำนวณ offset สำหรับ SQL query
func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// ApplyPagination ใช้ pagination กับ GORM query
func ApplyPagination[T any](
	db *gorm.DB,
	params PaginationParams,
	result *[]T,
) (*PaginationResult, error) {
	params.Normalize()

	// ใช้ preload ถ้ามี
	query := db
	for _, preload := range params.Preloads {
		query = query.Preload(preload)
	}

	// คำนวณจำนวนรายการทั้งหมด
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, err
	}

	// คำนวณจำนวนหน้าทั้งหมด
	totalPages := (count + int64(params.Limit) - 1) / int64(params.Limit)

	// ตรวจสอบและทำความสะอาด OrderBy เพื่อป้องกัน SQL Injection
	orderBy := sanitizeOrderBy(params.OrderBy)

	// ใช้ pagination และการจัดเรียง
	if err := query.Limit(params.Limit).Offset(params.GetOffset()).
		Order(orderBy).
		Find(result).Error; err != nil {
		return nil, err
	}

	return &PaginationResult{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      count,
		TotalPages: totalPages,
	}, nil
}

// ApplySearch เพิ่มเงื่อนไข search ให้กับ query
func ApplySearch(db *gorm.DB, search string, fields ...string) *gorm.DB {
	if search == "" || len(fields) == 0 {
		return db
	}

	// ทำความสะอาด search term เพื่อป้องกัน SQL Injection
	searchTerm := sanitizeSearchTerm(search)

	query := db.Session(&gorm.Session{})
	args := []interface{}{}

	// สร้างเงื่อนไข search ปลอดภัย
	conditions := []string{}

	for _, field := range fields {
		// ตรวจสอบชื่อคอลัมน์ให้ปลอดภัย
		sanitizedField := sanitizeColumnName(field)
		conditions = append(conditions, "LOWER("+sanitizedField+") LIKE ?")
		args = append(args, "%"+strings.ToLower(searchTerm)+"%")
	}

	if len(conditions) > 0 {
		whereClause := "(" + strings.Join(conditions, " OR ") + ")"
		query = query.Where(whereClause, args...)
	}

	return query
}

// sanitizeOrderBy ทำความสะอาด OrderBy เพื่อป้องกัน SQL Injection
func sanitizeOrderBy(orderBy string) string {
	// เพิ่มรายชื่อคอลัมน์ที่อนุญาตให้ใช้ในการเรียงลำดับ
	allowedColumns := map[string]bool{
		"id":           true,
		"created_at":   true,
		"updated_at":   true,
		"published_at": true,
		"title":        true,
		"name":         true,
		"email":        true,
		"status":       true,
		"last_seen":    true,
		"last_login":   true,
	}

	orderByParts := strings.Split(orderBy, " ")
	if len(orderByParts) < 1 || len(orderByParts) > 2 {
		return "created_at desc" // ค่าเริ่มต้นที่ปลอดภัย
	}

	column := strings.ToLower(orderByParts[0])
	if !allowedColumns[column] {
		return "created_at desc" // ค่าเริ่มต้นที่ปลอดภัย
	}

	direction := "asc"
	if len(orderByParts) == 2 {
		if strings.ToLower(orderByParts[1]) == "desc" {
			direction = "desc"
		}
	}

	return column + " " + direction
}

// sanitizeColumnName ตรวจสอบและทำความสะอาดชื่อคอลัมน์เพื่อป้องกัน SQL Injection
func sanitizeColumnName(column string) string {
	// เพิ่มรายชื่อคอลัมน์ที่อนุญาตให้ใช้ในการค้นหา
	allowedColumns := map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
		"title":      true,
		"content":    true,
		"slug":       true,
		"summary":    true,
		"status":     true,
		"name":       true,
		"email":      true,
		"device_id":  true,
		"last_seen":  true,
		"last_login": true,
	}

	if !allowedColumns[column] {
		return "id" // ค่าเริ่มต้นที่ปลอดภัย
	}

	return column
}

// sanitizeSearchTerm ทำความสะอาด search term เพื่อป้องกัน SQL Injection
func sanitizeSearchTerm(term string) string {
	// ตัดอักขระพิเศษที่อาจใช้ในการทำ SQL Injection ออก
	dangerousChars := []string{"'", "\"", ";", "--", "/*", "*/", "\\"}
	result := term

	for _, char := range dangerousChars {
		result = strings.ReplaceAll(result, char, "")
	}

	// จำกัดความยาวเพื่อป้องกันการโจมตี
	if len(result) > 100 {
		result = result[:100]
	}

	return result
}
