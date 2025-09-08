package ports

import (
	articleRepo "dashboard-starter/internal/domain/article/repository"
	"dashboard-starter/internal/domain/auth/repository"
	deviceRepo "dashboard-starter/internal/domain/device/repository"
	userRepo "dashboard-starter/internal/domain/user/repository"
)

// RepositoryPorts aggregates all repository interfaces
type RepositoryPorts struct {
	AdminRepo   repository.AdminRepository
	TokenRepo   repository.TokenRepository
	UserRepo    userRepo.UserRepository
	DeviceRepo  deviceRepo.DeviceRepository
	ArticleRepo articleRepo.ArticleRepository
}
