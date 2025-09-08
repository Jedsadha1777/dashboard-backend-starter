package database

import (
	"dashboard-starter/internal/domain/auth/entity"
	"dashboard-starter/internal/domain/auth/repository"
	"time"

	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) repository.TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(token *entity.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) GetByToken(token string) (*entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken
	err := r.db.Where("token = ? AND is_revoked = ? AND expires_at > ?",
		token, false, time.Now()).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *tokenRepository) RevokeToken(token string) error {
	return r.db.Model(&entity.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error
}

func (r *tokenRepository) RevokeAllUserTokens(userID uint, userType string) error {
	return r.db.Model(&entity.RefreshToken{}).
		Where("user_id = ? AND user_type = ? AND is_revoked = ?", userID, userType, false).
		Update("is_revoked", true).Error
}

func (r *tokenRepository) CleanupExpiredTokens() error {
	return r.db.Where("expires_at < ? OR is_revoked = ?", time.Now(), true).
		Delete(&entity.RefreshToken{}).Error
}

func (r *tokenRepository) GetUserTokens(userID uint, userType string) ([]*entity.RefreshToken, error) {
	var tokens []*entity.RefreshToken
	err := r.db.Where("user_id = ? AND user_type = ? AND is_revoked = ? AND expires_at > ?",
		userID, userType, false, time.Now()).Find(&tokens).Error
	return tokens, err
}

func (r *tokenRepository) DeleteExpiredBefore(before time.Time) error {
	return r.db.Where("expires_at < ?", before).Delete(&entity.RefreshToken{}).Error
}
