// db/repository.go
package db

import (
	"gorm.io/gorm"
)

// Repository เป็น interface ที่ระบุ operations พื้นฐานสำหรับแต่ละ entity
type Repository[T any] interface {
	FindByID(id interface{}) (*T, error)
	FindAll(conditions ...interface{}) ([]T, error)
	FindOne(conditions ...interface{}) (*T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id interface{}) error
	Count(conditions ...interface{}) (int64, error)
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

// FindByID ค้นหา entity ด้วย ID
func (r *GormRepository[T]) FindByID(id interface{}) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
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
		query = query.Where(conditions[0], conditions[1:]...)
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
		query = query.Where(conditions[0], conditions[1:]...)
	}

	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
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
	return r.db.Delete(&entity, id).Error
}

// Count นับจำนวน entities ที่ตรงกับเงื่อนไข
func (r *GormRepository[T]) Count(conditions ...interface{}) (int64, error) {
	var count int64
	var entity T
	query := r.db.Model(&entity)

	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}

	err := query.Count(&count).Error
	return count, err
}
