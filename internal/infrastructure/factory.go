package infrastructure

import (
	"dashboard-starter/internal/application/ports"

	"dashboard-starter/internal/infrastructure/cache"
	"dashboard-starter/internal/infrastructure/database"
	"dashboard-starter/internal/infrastructure/email"
	"dashboard-starter/internal/infrastructure/logger"

	"gorm.io/gorm"
)

type InfrastructureFactory struct {
	db *gorm.DB
}

func NewInfrastructureFactory(db *gorm.DB) *InfrastructureFactory {
	return &InfrastructureFactory{db: db}
}

func (f *InfrastructureFactory) CreateRepositoryPorts() *ports.RepositoryPorts {
	return &ports.RepositoryPorts{
		AdminRepo:   database.NewAdminRepository(f.db),
		TokenRepo:   database.NewTokenRepository(f.db),
		UserRepo:    database.NewUserRepository(f.db),
		DeviceRepo:  database.NewDeviceRepository(f.db),
		ArticleRepo: database.NewArticleRepository(f.db),
	}
}

func (f *InfrastructureFactory) CreateCacheService(config cache.Config) (ports.CacheService, error) {
	return cache.NewRedisCache(config)
}

func (f *InfrastructureFactory) CreateEmailService(config email.Config) ports.EmailService {
	return email.NewSMTPEmailService(config)
}

func (f *InfrastructureFactory) CreateLoggerService(env string) (*logger.ZapLogger, error) {
	zapLogger, err := logger.NewZapLogger(env)
	if err != nil {
		return nil, err
	}
	return zapLogger, nil

}
