package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nickname  string  `gorm:"column:nickname;type:varchar(255);unique"` // min(2), max(12)
	Biography *string `gorm:"column:biography;type:varchar(255);"`      // min(0), max(14)
}

type CreateUserParams struct {
	Nickname  string
	Biography *string
}
