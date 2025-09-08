package repository

import (
	"dashboard-starter/internal/domain/auth/entity"
	"time"
)

type TokenRepository interface {
	Create(token *entity.RefreshToken) error
	GetByToken(token string) (*entity.RefreshToken, error)
	RevokeToken(token string) error
	RevokeAllUserTokens(userID uint, userType string) error
	CleanupExpiredTokens() error
	GetUserTokens(userID uint, userType string) ([]*entity.RefreshToken, error)
	DeleteExpiredBefore(before time.Time) error
}
