package domain

import (
	"github.com/roka-crew/sam2ooh2-api/pkg/query"
	"gorm.io/gorm"
)

type Users []User

type User struct {
	gorm.Model
	Nickname  string  `gorm:"column:nickname;type:varchar(255);unique"` // min(2), max(12)
	Biography *string `gorm:"column:biography;type:varchar(255);"`      // min(0), max(14)
}

type CreateUserParams struct {
	Nickname  string
	Biography *string
}

type ListUsersParams struct {
	// conditions
	IDs         []uint
	Nicknames   []string
	Biographies []string

	// sort
	Sorts query.Sorts

	// pagiantion
	Page query.Page
}
