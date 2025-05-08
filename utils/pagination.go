// utils/pagination.go
package utils

import (
	"strings"

	"gorm.io/gorm"
)

// PaginationParams พารามิเตอร์สำหรับ pagination
type PaginationParams struct {
	Page    int    `json:"page" form:"page"`
	Limit   int    `json:"limit" form:"limit"`
	OrderBy string `json:"order_by" form:"order_by"`
	Search  string `json:"search" form:"search"`
	Status  string `json:"status" form:"status"`
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
		Page:    1,
		Limit:   10,
		OrderBy: "created_at desc",
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

	// Clone the query to count total records
	var count int64
	countQuery := db
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := (count + int64(params.Limit) - 1) / int64(params.Limit)

	// Apply pagination and get results
	err := db.Limit(params.Limit).Offset(params.GetOffset()).
		Order(params.OrderBy).
		Find(result).Error

	if err != nil {
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

	searchTerm := "%" + strings.ToLower(search) + "%"

	query := db.Session(&gorm.Session{})
	for i, field := range fields {
		if i == 0 {
			query = query.Where("LOWER("+field+") LIKE ?", searchTerm)
		} else {
			query = query.Or("LOWER("+field+") LIKE ?", searchTerm)
		}
	}

	return query
}
