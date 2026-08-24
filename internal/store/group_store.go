package store

import (
	"context"

	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/payload"
	"github.com/roka-crew/sam2ooh2-api/internal/sqlite"
	"gorm.io/gorm/clause"
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

func (s *GroupStore) ListGroups(c context.Context, param payload.ListGroupsParam) (domain.Groups, error) {
	db := s.db.WithContext(c)

	if len(param.IDs) > 0 {
		db = db.Where("id IN ?", param.IDs)
	}

	if len(param.Authors) > 0 {
		db = db.Where("author IN ?", param.Authors)
	}

	if len(param.PageCounts) > 0 {
		db = db.Where("page_count IN ?", param.PageCounts)
	}

	if len(param.Publishers) > 0 {
		db = db.Where("publisher IN ?", param.Publishers)
	}

	if len(param.Descriptions) > 0 {
		db = db.Where("description IN ?", param.Descriptions)
	}

	for _, sort := range param.Sorts {
		db = db.Order(clause.OrderByColumn{
			Column: clause.Column{Name: sort.By},
			Desc:   sort.Order.IsDESC(),
		})
	}

	if param.Offset > 0 {
		db = db.Offset(param.Offset)
	}

	if param.Limit > 0 {
		db = db.Limit(param.Limit)
	}

	var groups domain.Groups
	if err := db.Find(&groups).Error; err != nil {
		return domain.Groups{}, err
	}

	return groups, nil
}

func (s *GroupStore) PatchGroup(ctx context.Context, param payload.PatchGroupParam) error {
	updates := map[string]any{}

	if param.Title != "" {
		updates["title"] = param.Title
	}
	if param.Author != "" {
		updates["author"] = param.Author
	}
	if param.PageCount > 0 {
		updates["page_count"] = param.PageCount
	}
	if param.Publisher != "" {
		updates["publisher"] = param.Publisher
	}
	if param.Description != "" {
		updates["description"] = param.Description
	}

	if len(updates) == 0 {
		return nil
	}

	return s.db.WithContext(ctx).Model(&domain.Group{}).Where("id = ?", param.ID).Updates(updates).Error
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
