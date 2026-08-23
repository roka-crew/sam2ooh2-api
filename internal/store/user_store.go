package store

import (
	"context"

	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/payload"
	"github.com/roka-crew/sam2ooh2-api/internal/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserStore struct {
	db *sqlite.Sqlite
}

func NewUserStore(
	db *sqlite.Sqlite,
) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) CreateUser(c context.Context, param payload.CreateUserParam) (domain.User, error) {
	user := domain.User{
		Nickname:  param.Nickname,
		Biography: param.Biography,
	}

	if err := gorm.G[domain.User](s.db.DB).Create(c, &user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s *UserStore) ListUsers(c context.Context, param payload.ListUsersParam) (domain.Users, error) {
	db := s.db.WithContext(c)

	if len(param.IDs) > 0 {
		db = db.Where("id IN ?", param.IDs)
	}

	if len(param.Nicknames) > 0 {
		db = db.Where("nickname IN ?", param.Nicknames)
	}

	if len(param.Biographies) > 0 {
		db = db.Where("biography IN ?", param.Biographies)
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

	var users domain.Users
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStore) PatchUser(c context.Context, param payload.PatchUserParam) error {
	updates := map[string]any{}

	if param.Nickname != nil {
		updates["nickname"] = *param.Nickname
	}
	if param.Biography != nil {
		updates["biography"] = *param.Biography 
	}

	if len(updates) == 0 {
		return nil 
	}

	return s.db.WithContext(c).Model(&domain.User{}).Where("id = ?", param.ID).Updates(updates).Error
}

func (s *UserStore) PatchUsers(c context.Context, param payload.PatchUsersParam) error {
	err := s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		for _, patchUserParam := range param {
			updates := map[string]any{}

			if patchUserParam.Nickname != nil {
				updates["nickname"] = *patchUserParam.Nickname
			}
			if patchUserParam.Biography != nil {
				updates["biography"] = *patchUserParam.Biography
			}

			if len(updates) == 0 {
				continue
			}

			if err := tx.Model(&domain.User{}).Where("id = ?", patchUserParam.ID).Updates(updates).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (s *UserStore) DeleteUser(c context.Context, param payload.DeleteUserParam) error {
	db := s.db.WithContext(c)

	if param.ID > 0 {
		db = db.Where("id = ?", param.ID)
	}

	if param.Nickname != "" {
		db = db.Where("nickname = ?", param.Nickname)
	}

	if param.IsHardDelete {
		db = db.Unscoped()
	}

	if err := db.Delete(new(domain.User)).Error; err != nil {
		return err
	}

	return nil
}
