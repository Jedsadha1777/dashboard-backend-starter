package repository

import "dashboard-starter/internal/domain/device/entity"

type DeviceRepository interface {
	GetByID(id uint) (*entity.Device, error)
	GetByDeviceID(deviceID string) (*entity.Device, error)
	Create(device *entity.Device) error
	Update(device *entity.Device) error
	Delete(id uint) error
	List(offset, limit int, search string) ([]*entity.Device, int64, error)
	Count(conditions ...interface{}) (int64, error)
}
