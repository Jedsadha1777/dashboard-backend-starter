// db/repository.go
package db

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Repository เป็น interface ที่ระบุ operations พื้นฐานสำหรับแต่ละ entity
type Repository[T any] interface {
	FindByID(id interface{}) (*T, error)
	FindAll(conditions ...interface{}) ([]T, error)
	FindOne(conditions ...interface{}) (*T, error)
	FindWithPreload(preloads []string, id interface{}) (*T, error)
	FindOneWithPreload(preloads []string, conditions ...interface{}) (*T, error)
	FindAllWithPreload(preloads []string, conditions ...interface{}) ([]T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id interface{}) error
	Count(conditions ...interface{}) (int64, error)
	Paginate(query *gorm.DB, page, limit int, result *[]T) (int64, error)
}

// GormRepository implementation ของ Repository ด้วย GORM
type GormRepository[T any] struct {
	db *gorm.DB
}

// NewRepository สร้าง repository ใหม่
func NewRepository[T any]() *GormRepository[T] {
	return &GormRepository[T]{
		db: DB,
	}
}

// sanitizeID ทำความสะอาด id เพื่อป้องกัน SQL Injection
func sanitizeID(id interface{}) (interface{}, error) {
	switch v := id.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v, nil
	case string:
		// พยายามแปลงเป็นตัวเลข
		if intVal, err := strconv.ParseUint(v, 10, 64); err == nil {
			return intVal, nil
		}
		return nil, errors.New("invalid ID format")
	default:
		return nil, errors.New("unsupported ID type")
	}
}

// sanitizeConditions ทำความสะอาด conditions เพื่อป้องกัน SQL Injection
func sanitizeConditions(conditions ...interface{}) ([]interface{}, error) {
	if len(conditions) == 0 {
		return conditions, nil
	}

	// ตรวจสอบ condition แรกที่ควรเป็น string
	if whereClause, ok := conditions[0].(string); ok {
		// ตรวจสอบ SQL Injection ใน where clause
		if strings.ContainsAny(whereClause, ";--'\"\\") {
			return nil, errors.New("invalid characters in where clause")
		}

		// ตรวจสอบว่าเป็น parameterized query หรือไม่
		if !strings.Contains(whereClause, "?") && len(conditions) > 1 {
			return nil, errors.New("non-parameterized queries are not allowed")
		}
	}

	return conditions, nil
}

// FindByID ค้นหา entity ด้วย ID
func (r *GormRepository[T]) FindByID(id interface{}) (*T, error) {
	var entity T

	// ทำความสะอาด id
	cleanID, err := sanitizeID(id)
	if err != nil {
		return nil, err
	}

	err = r.db.First(&entity, cleanID).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll ค้นหา entities ทั้งหมดที่ตรงกับเงื่อนไข
func (r *GormRepository[T]) FindAll(conditions ...interface{}) ([]T, error) {
	var entities []T
	query := r.db

	if len(conditions) > 0 {
		// ทำความสะอาด conditions
		cleanConditions, err := sanitizeConditions(conditions...)
		if err != nil {
			return nil, err
		}

		query = query.Where(cleanConditions[0], cleanConditions[1:]...)
	}

	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// FindOne ค้นหา entity เดียวที่ตรงกับเงื่อนไข
func (r *GormRepository[T]) FindOne(conditions ...interface{}) (*T, error) {
	var entity T
	query := r.db

	if len(conditions) > 0 {
		// ทำความสะอาด conditions
		cleanConditions, err := sanitizeConditions(conditions...)
		if err != nil {
			return nil, err
		}

		query = query.Where(cleanConditions[0], cleanConditions[1:]...)
	}

	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// sanitizePreloads ตรวจสอบและทำความสะอาด preloads
func sanitizePreloads(preloads []string) ([]string, error) {
	// รายชื่อ associations ที่อนุญาตให้ preload
	allowedPreloads := map[string]bool{
		"Admin":         true,
		"User":          true,
		"Device":        true,
		"Articles":      true,
		"RefreshTokens": true,
	}

	var cleanPreloads []string

	for _, preload := range preloads {
		// ตรวจสอบชื่อ association
		if !allowedPreloads[preload] {
			return nil, fmt.Errorf("preload '%s' is not allowed", preload)
		}

		cleanPreloads = append(cleanPreloads, preload)
	}

	return cleanPreloads, nil
}

// FindWithPreload ค้นหา entity พร้อม preload ความสัมพันธ์ที่ระบุ
func (r *GormRepository[T]) FindWithPreload(preloads []string, id interface{}) (*T, error) {
	var entity T
	query := r.db

	// ทำความสะอาด id
	cleanID, err := sanitizeID(id)
	if err != nil {
		return nil, err
	}

	// ทำความสะอาด preloads
	cleanPreloads, err := sanitizePreloads(preloads)
	if err != nil {
		return nil, err
	}

	for _, preload := range cleanPreloads {
		query = query.Preload(preload)
	}

	err = query.First(&entity, cleanID).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindOneWithPreload ค้นหา entity หนึ่งรายการพร้อม preload ความสัมพันธ์
func (r *GormRepository[T]) FindOneWithPreload(preloads []string, conditions ...interface{}) (*T, error) {
	var entity T
	query := r.db

	// ทำความสะอาด preloads
	cleanPreloads, err := sanitizePreloads(preloads)
	if err != nil {
		return nil, err
	}

	for _, preload := range cleanPreloads {
		query = query.Preload(preload)
	}

	if len(conditions) > 0 {
		// ทำความสะอาด conditions
		cleanConditions, err := sanitizeConditions(conditions...)
		if err != nil {
			return nil, err
		}

		query = query.Where(cleanConditions[0], cleanConditions[1:]...)
	}

	err = query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAllWithPreload ค้นหา entities ทั้งหมดพร้อม preload ความสัมพันธ์
func (r *GormRepository[T]) FindAllWithPreload(preloads []string, conditions ...interface{}) ([]T, error) {
	var entities []T
	query := r.db

	// ทำความสะอาด preloads
	cleanPreloads, err := sanitizePreloads(preloads)
	if err != nil {
		return nil, err
	}

	for _, preload := range cleanPreloads {
		query = query.Preload(preload)
	}

	if len(conditions) > 0 {
		// ทำความสะอาด conditions
		cleanConditions, err := sanitizeConditions(conditions...)
		if err != nil {
			return nil, err
		}

		query = query.Where(cleanConditions[0], cleanConditions[1:]...)
	}

	err = query.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// Create สร้าง entity ใหม่
func (r *GormRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// Update อัปเดต entity ที่มีอยู่
func (r *GormRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete ลบ entity ด้วย ID
func (r *GormRepository[T]) Delete(id interface{}) error {
	var entity T

	// ทำความสะอาด id
	cleanID, err := sanitizeID(id)
	if err != nil {
		return err
	}

	return r.db.Delete(&entity, cleanID).Error
}

// Count นับจำนวน entities ที่ตรงกับเงื่อนไข
func (r *GormRepository[T]) Count(conditions ...interface{}) (int64, error) {
	var count int64
	var entity T
	query := r.db.Model(&entity)

	if len(conditions) > 0 {
		// ทำความสะอาด conditions
		cleanConditions, err := sanitizeConditions(conditions...)
		if err != nil {
			return 0, err
		}

		query = query.Where(cleanConditions[0], cleanConditions[1:]...)
	}

	err := query.Count(&count).Error
	return count, err
}

// Paginate ทำ pagination กับ query ที่ระบุ
func (r *GormRepository[T]) Paginate(query *gorm.DB, page, limit int, result *[]T) (int64, error) {
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	// ตรวจสอบและปรับค่า pagination
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	err = query.Limit(limit).Offset(offset).Find(result).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
