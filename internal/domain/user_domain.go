package domain

import (
	"context"

	"github.com/roka-crew/sam2ooh2-api/internal/payload"
	"gorm.io/gorm"
)

type Users []User

type User struct {
	gorm.Model
	Nickname  string  `gorm:"column:nickname;type:varchar(255);unique"` // min(2), max(12)
	Biography *string `gorm:"column:biography;type:varchar(255);"`      // min(0), max(14)
}

type UserStore interface {
	CreateUser(context.Context, payload.CreateUserParam) (User, error)
	ListUsers(context.Context, payload.ListUsersParam) (Users, error)
	PatchUser(context.Context, payload.PatchUserParam) error
	PatchUsers(context.Context, payload.PatchUsersParam) error
	DeleteUser(context.Context, payload.DeleteUserParam) error
}

type UserService interface {
	CreateUser(context.Context, payload.CreateUserRequest) (payload.CreateUserResponse, error)
}
