package services

import (
	"dashboard-starter/db"
	"dashboard-starter/models"
	"strings"
)

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := db.DB.Find(&users)
	return users, result.Error
}

func GetUsers(search string, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var count int64

	query := db.DB.Model(&models.User{})

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", like, like)
	}

	// นับจำนวนทั้งหมดก่อน (ก่อน limit)
	query.Count(&count)

	offset := (page - 1) * limit
	result := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users)

	return users, count, result.Error
}

func CreateUser(user *models.User) error {
	return db.DB.Create(user).Error
}

func UpdateUser(id string, updated *models.User) (*models.User, error) {
	var user models.User

	if err := db.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	user.Name = updated.Name
	user.Email = updated.Email
	db.DB.Save(&user)
	return &user, nil
}

func DeleteUser(id string) error {
	return db.DB.Delete(&models.User{}, id).Error
}
