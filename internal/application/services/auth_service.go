package services

import (
	"dashboard-starter/internal/application/dto"
	"dashboard-starter/internal/domain/auth/entity"
	"dashboard-starter/internal/domain/auth/repository"
	"dashboard-starter/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	adminRepo repository.AdminRepository
	tokenRepo repository.TokenRepository
}

func NewAuthService(adminRepo repository.AdminRepository, tokenRepo repository.TokenRepository) *AuthService {
	return &AuthService{
		adminRepo: adminRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *AuthService) LoginAdmin(input dto.LoginInput) (*dto.LoginResponse, error) {
	admin, err := s.adminRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	admin.LastLogin = time.Now()
	admin.TokenVersion++
	if err := s.adminRepo.Update(admin); err != nil {
		return nil, err
	}

	token, exp, err := utils.GenerateToken(admin.ID, "admin", admin.TokenVersion)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.CreateRefreshToken(admin.ID, "admin")
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken.Token,
		ExpiresAt:    exp,
		UserID:       admin.ID,
		UserType:     "admin",
	}, nil
}

func (s *AuthService) GetAdminProfile(adminID uint) (*entity.Admin, error) {
	admin, err := s.adminRepo.GetByID(adminID)
	if err != nil {
		return nil, errors.New("admin not found")
	}

	admin.Password = ""
	return admin, nil
}

func (s *AuthService) LogoutAdmin(adminID uint) error {
	return s.adminRepo.IncrementTokenVersion(adminID)
}

func (s *AuthService) RefreshToken(input dto.RefreshTokenInput) (*dto.LoginResponse, error) {
	refreshToken, err := s.tokenRepo.GetByToken(input.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if refreshToken.IsRevoked || refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired or revoked")
	}

	var tokenVersion int
	if refreshToken.UserType == "admin" {
		admin, err := s.adminRepo.GetByID(refreshToken.UserID)
		if err != nil {
			return nil, err
		}
		tokenVersion = admin.TokenVersion
	}

	token, exp, err := utils.GenerateToken(refreshToken.UserID, refreshToken.UserType, tokenVersion)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:     token,
		ExpiresAt: exp,
		UserID:    refreshToken.UserID,
		UserType:  refreshToken.UserType,
	}, nil
}

func (s *AuthService) CreateRefreshToken(userID uint, userType string) (*entity.RefreshToken, error) {
	tokenString, expiresAt, err := utils.GenerateRefreshToken(userID, userType)
	if err != nil {
		return nil, err
	}

	refreshToken := &entity.RefreshToken{
		Token:     tokenString,
		UserID:    userID,
		UserType:  userType,
		ExpiresAt: expiresAt,
		IsRevoked: false,
	}

	if err := s.tokenRepo.Create(refreshToken); err != nil {
		return nil, err
	}

	return refreshToken, nil
}
