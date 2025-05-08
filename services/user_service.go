package services

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"dashboard-starter/utils"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type UserService struct {
	repo *db.GormRepository[models.User]
}

func NewUserService() *UserService {
	return &UserService{
		repo: db.NewRepository[models.User](),
	}
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	// แปลง id เป็น uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	return s.repo.FindByID(uint(idUint))
}

// GetUsers retrieves users with pagination and search
func (s *UserService) GetUsers(params utils.PaginationParams) ([]models.User, *utils.PaginationResult, error) {
	var users []models.User

	// สร้าง query
	query := db.DB.Model(&models.User{})

	// ใช้ search ถ้ามี
	if params.Search != "" {
		query = utils.ApplySearch(query, params.Search, "name", "email")
	}

	// ใช้ pagination
	result, err := utils.ApplyPagination(query, params, &users)
	if err != nil {
		return nil, nil, err
	}

	return users, result, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(input *models.UserInput) (*models.User, error) {
	// ตรวจสอบ email ซ้ำ
	count, err := s.repo.Count("email = ?", input.Email)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("email already exists")
	}

	// สร้าง user ใหม่
	user := &models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	// บันทึกลงฐานข้อมูล
	err = db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Create(user)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id string, input *models.UserInput) (*models.User, error) {
	// แปลง id เป็น uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	// ค้นหา user
	user, err := s.repo.FindByID(uint(idUint))
	if err != nil {
		return nil, err
	}

	// ตรวจสอบ email ซ้ำ (ถ้ามีการเปลี่ยน email)
	if input.Email != user.Email {
		count, err := s.repo.Count("email = ? AND id != ?", input.Email, user.ID)
		if err != nil {
			return nil, err
		}

		if count > 0 {
			return nil, errors.New("email already exists")
		}
	}

	// อัปเดตข้อมูล user
	user.Name = input.Name
	user.Email = input.Email

	// บันทึกการเปลี่ยนแปลง
	err = db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Update(user)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string) error {
	// แปลง id เป็น uint
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errors.New("invalid ID format")
	}

	// ตรวจสอบว่า user มีอยู่จริง
	_, err = s.repo.FindByID(uint(idUint))
	if err != nil {
		return err
	}

	// ลบ user
	return db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Delete(uint(idUint))
	})
}
