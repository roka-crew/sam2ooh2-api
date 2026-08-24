package domain

import (
	"context"

	"github.com/roka-crew/sam2ooh2-api/internal/payload"
	"gorm.io/gorm"
)

type Groups []Group

type Group struct {
	gorm.Model
	Title       string `gorm:"type:varchar(200);index"`
	Author      string `gorm:"type:varchar(100)"`
	PageCount   int    `gorm:"type:integer;check:page_count > 0"`
	Publisher   string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:varchar(1024)"`

	// Group < N:M > User
	Users []User `gorm:"many2many:users_groups;"`

	// Group < 1:N > Chapter : 그룹 내에 생성된 챕터들
	Chapters []Chapter `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE;"`
}

type GroupStore interface {
	CreateGroup(ctx context.Context, param payload.CreateGroupParam) (Group, error)
	ListGroups(ctx context.Context, param payload.ListGroupsParam) (Groups, error)
	PatchGroup(ctx context.Context, param payload.PatchGroupParam) error
	DeleteGroup(ctx context.Context, param payload.DeleteGroupParam) error

	// N:M 멤버십 조작 API
	AddUserToGroup(ctx context.Context, param payload.AddUserToGroupParam) error
	RemoveUserFromGroup(ctx context.Context, param payload.RemoveUserFromGroup) error
	CountGroupMembers(ctx context.Context, param payload.CountGroupUsers) (int, error)
	IsUserInGroup(ctx context.Context, param payload.IsUserInGroup) (bool, error)
}
