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

func (s *UserStore) PatchUser(c context.Context, params payload.PatchUserParam) error {
	user := domain.User{
		Model: gorm.Model{ID: params.ID},
	}

	if params.Nickname != nil {
		user.Nickname = *params.Nickname
	}

	if params.Biography != nil {
		user.Biography = params.Biography
	}

	if err := s.db.WithContext(c).Updates(user).Error; err != nil {
		return err
	}

	return nil
}

func (s *UserStore) PatchUsers(c context.Context, param payload.PatchUsersParam) error {
	users := domain.Users{}

	for _, patchUserParam := range param {
		user := domain.User{
			Model: gorm.Model{ID: patchUserParam.ID},
		}

		if patchUserParam.Nickname != nil {
			user.Nickname = *patchUserParam.Nickname
		}

		if patchUserParam.Biography != nil {
			user.Biography = patchUserParam.Biography
		}

		users = append(users, user)
	}

	if err := s.db.WithContext(c).Updates(users).Error; err != nil {
		return err
	}

	return nil
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
