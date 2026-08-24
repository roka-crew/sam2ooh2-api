package store

import (
	"context"

	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/payload"
	"github.com/roka-crew/sam2ooh2-api/internal/sqlite"
)

type GroupStore struct {
	db *sqlite.Sqlite
}

func NewGroupStore(
	db *sqlite.Sqlite,
) domain.GroupStore {
	return &GroupStore{
		db: db,
	}
}

func (s *GroupStore) CreateGroup(c context.Context, param payload.CreateGroupParam) (domain.Group, error) {
	group := domain.Group{
		Title:       param.Title,
		Author:      param.Author,
		PageCount:   param.PageCount,
		Publisher:   param.Publisher,
		Description: param.Description,
	}

	if err := s.db.WithContext(c).Create(&group).Error; err != nil {
		return domain.Group{}, err
	}

	return group, nil
}

func (s *GroupStore) ListGroups(ctx context.Context, param payload.ListGroupsParam) (domain.Groups, error) {
	return domain.Groups{}, nil
}
func (s *GroupStore) PatchGroup(ctx context.Context, param payload.PatchGroupParam) error {
	return nil
}
func (s *GroupStore) DeleteGroup(ctx context.Context, param payload.DeleteGroupParam) error {
	return nil
}

// N:M 멤버십 조작 API
func (s *GroupStore) AddUserToGroup(ctx context.Context, param payload.AddUserToGroupParam) error {
	return nil
}
func (s *GroupStore) RemoveUserFromGroup(ctx context.Context, param payload.RemoveUserFromGroup) error {
	return nil
}
func (s *GroupStore) CountGroupMembers(ctx context.Context, param payload.CountGroupUsers) (int, error) {
	return 0, nil
}
func (s *GroupStore) IsUserInGroup(ctx context.Context, param payload.IsUserInGroup) (bool, error) {
	return false, nil
}
