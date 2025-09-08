package services

import (
	"dashboard-starter/internal/application/dto"
	"dashboard-starter/internal/domain/user/entity"
	"dashboard-starter/internal/domain/user/repository"
	"dashboard-starter/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserApplicationService struct {
	userRepo    repository.UserRepository
	authService *AuthService
}

func NewUserApplicationService(userRepo repository.UserRepository, authService *AuthService) *UserApplicationService {
	return &UserApplicationService{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (s *UserApplicationService) RegisterUser(input dto.UserRegistrationInput) (*dto.LoginResponse, error) {
	// Validate password strength
	if isStrong, msg := utils.IsStrongPassword(input.Password); !isStrong {
		return nil, errors.New("password not strong enough: " + msg)
	}

	// Check if email exists
	if _, err := s.userRepo.GetByEmail(input.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entity.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     string(hashedPassword),
		TokenVersion: 1,
		LastLogin:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	token, exp, err := utils.GenerateToken(user.ID, "user", user.TokenVersion)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.authService.CreateRefreshToken(user.ID, "user")

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken.Token,
		ExpiresAt:    exp,
		UserID:       user.ID,
		UserType:     "user",
	}, nil
}

func (s *UserApplicationService) LoginUser(input dto.UserLoginInput) (*dto.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Update last login and token version
	user.LastLogin = time.Now()
	user.TokenVersion++
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Generate tokens
	token, exp, err := utils.GenerateToken(user.ID, "user", user.TokenVersion)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.authService.CreateRefreshToken(user.ID, "user")

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken.Token,
		ExpiresAt:    exp,
		UserID:       user.ID,
		UserType:     "user",
	}, nil
}

func (s *UserApplicationService) GetUserProfile(userID uint) (*entity.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *UserApplicationService) UpdateUserProfile(userID uint, input dto.UserInput) (*entity.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Check for duplicate email if changed
	if input.Email != user.Email {
		if count, _ := s.userRepo.Count("email = ? AND id != ?", input.Email, user.ID); count > 0 {
			return nil, errors.New("email already exists")
		}
	}

	user.Name = input.Name
	user.Email = input.Email

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserApplicationService) ChangePassword(userID uint, input dto.UserChangePasswordInput) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Validate new password strength
	if isStrong, msg := utils.IsStrongPassword(input.NewPassword); !isStrong {
		return errors.New("new password not strong enough: " + msg)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.TokenVersion++ // Invalidate existing tokens

	return s.userRepo.Update(user)
}

func (s *UserApplicationService) ListUsers(page, limit int, search string) ([]*entity.User, int64, error) {
	offset := (page - 1) * limit
	return s.userRepo.List(offset, limit, search)
}

func (s *UserApplicationService) CreateUserByAdmin(input dto.UserInput, adminID uint) (*entity.User, string, error) {
	// Check for duplicate email
	if _, err := s.userRepo.GetByEmail(input.Email); err == nil {
		return nil, "", errors.New("email already exists")
	}

	// Generate temporary password
	tempPassword := utils.GenerateRandomPassword(12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &entity.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     string(hashedPassword),
		TokenVersion: 1,
		AdminID:      adminID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	return user, tempPassword, nil
}
