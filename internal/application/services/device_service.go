package services

import (
	"crypto/rand"
	"dashboard-starter/internal/application/dto"
	"dashboard-starter/internal/domain/device/entity"
	"dashboard-starter/internal/domain/device/repository"
	"dashboard-starter/utils"
	"encoding/hex"
	"errors"
	"time"
)

type DeviceApplicationService struct {
	deviceRepo  repository.DeviceRepository
	authService *AuthService
}

func NewDeviceApplicationService(deviceRepo repository.DeviceRepository, authService *AuthService) *DeviceApplicationService {
	return &DeviceApplicationService{
		deviceRepo:  deviceRepo,
		authService: authService,
	}
}

func (s *DeviceApplicationService) AuthenticateDevice(input dto.DeviceAuthInput) (*dto.LoginResponse, error) {
	// Find device by device ID
	device, err := s.deviceRepo.GetByDeviceID(input.DeviceID)
	if err != nil {
		return nil, errors.New("invalid device ID or API key")
	}

	// Verify API key
	if device.ApiKey != input.ApiKey {
		return nil, errors.New("invalid device ID or API key")
	}

	// Update device status
	device.LastSeen = time.Now()
	device.Status = "active"
	if err := s.deviceRepo.Update(device); err != nil {
		return nil, err
	}

	// Generate tokens
	token, exp, err := utils.GenerateToken(device.ID, "device", device.TokenVersion)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.authService.CreateRefreshToken(device.ID, "device")

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken.Token,
		ExpiresAt:    exp,
		UserID:       device.ID,
		UserType:     "device",
	}, nil
}

func (s *DeviceApplicationService) RegisterDevice(deviceID, name string) (*entity.Device, string, error) {
	// Check if device ID already exists
	if count, _ := s.deviceRepo.Count("device_id = ?", deviceID); count > 0 {
		return nil, "", errors.New("device ID already exists")
	}

	// Generate API key
	apiKey, err := s.generateAPIKey(32)
	if err != nil {
		return nil, "", err
	}

	device := &entity.Device{
		DeviceID:     deviceID,
		Name:         name,
		ApiKey:       apiKey,
		TokenVersion: 1,
		Status:       "inactive",
		LastSeen:     time.Now(),
	}

	if err := s.deviceRepo.Create(device); err != nil {
		return nil, "", err
	}

	return device, apiKey, nil
}

func (s *DeviceApplicationService) generateAPIKey(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *DeviceApplicationService) ListDevices(page, limit int, search string) ([]*entity.Device, int64, error) {
	offset := (page - 1) * limit
	return s.deviceRepo.List(offset, limit, search)
}

func (s *DeviceApplicationService) GetDevice(id uint) (*entity.Device, error) {
	return s.deviceRepo.GetByID(id)
}

func (s *DeviceApplicationService) ResetDeviceAPIKey(id uint) (*entity.Device, string, error) {
	device, err := s.deviceRepo.GetByID(id)
	if err != nil {
		return nil, "", err
	}

	// Generate new API key
	newAPIKey, err := s.generateAPIKey(32)
	if err != nil {
		return nil, "", err
	}

	device.ApiKey = newAPIKey
	device.TokenVersion++ // Invalidate existing tokens

	if err := s.deviceRepo.Update(device); err != nil {
		return nil, "", err
	}

	return device, newAPIKey, nil
}
