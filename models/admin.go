package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model          //  struct ฝัง field `ID`, `CreatedAt`, `UpdatedAt` DeletedAt,
	Email        string `gorm:"unique"`
	Password     string
	TokenVersion int
}
