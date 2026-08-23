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

type CreateUserParam struct {
	Nickname  string
	Biography *string
}

type ListUsersParam struct {
	// conditions
	IDs         []uint
	Nicknames   []string
	Biographies []string

	// sort
	Sorts query.Sorts

	// optinos
	Limit  int
	Offset int
}

type PatchUsersParam []PatchUserParam

type PatchUserParam struct {
	// conditions
	ID uint

	// fields
	Nickname  *string
	Biography *string
}

type DeleteUserParam struct {
	// conditions
	ID       uint
	Nickname string

	// option
	IsHardDelete bool
}
