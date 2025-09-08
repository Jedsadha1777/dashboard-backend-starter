package database

import (
	"dashboard-starter/internal/domain/device/entity"
	"dashboard-starter/internal/domain/device/repository"

	"gorm.io/gorm"
)

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) repository.DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) GetByID(id uint) (*entity.Device, error) {
	var device entity.Device
	err := r.db.First(&device, id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) GetByDeviceID(deviceID string) (*entity.Device, error) {
	var device entity.Device
	err := r.db.Where("device_id = ?", deviceID).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) Create(device *entity.Device) error {
	return r.db.Create(device).Error
}

func (r *deviceRepository) Update(device *entity.Device) error {
	return r.db.Save(device).Error
}

func (r *deviceRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Device{}, id).Error
}

func (r *deviceRepository) List(offset, limit int, search string) ([]*entity.Device, int64, error) {
	var devices []*entity.Device
	var total int64

	query := r.db.Model(&entity.Device{})

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("device_id LIKE ? OR name LIKE ?", searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&devices).Error
	return devices, total, err
}

func (r *deviceRepository) Count(conditions ...interface{}) (int64, error) {
	var count int64
	query := r.db.Model(&entity.Device{})

	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}

	err := query.Count(&count).Error
	return count, err
}
