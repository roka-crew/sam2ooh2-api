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

	// User < N:M > Group : 사용자가 가입한 그룹들
	Groups []Group `gorm:"many2many:users_groups;"`

	// User < 1:N > Chapter : 사용자가 생성한 챕터(목표)들 (추가된 부분!)
	Chapters []Chapter `gorm:"foreignKey:UserID;"`

	// User < 1:N > Topic : 사용자가 작성한 토론 주제들 (추가된 부분!)
	Topics []Topic `gorm:"foreignKey:UserID;"`
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
