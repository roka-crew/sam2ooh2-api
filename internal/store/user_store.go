package store

import (
	"context"

	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserStore struct {
	db *sqlite.Sqlite
}

func NewUserStore(
	db *sqlite.Sqlite,
) (*UserStore, error) {
	return &UserStore{db: db}, nil
}

func (u *UserStore) CreateUser(ctx context.Context, param domain.CreateUserParam) (domain.User, error) {
	user := domain.User{
		Nickname:  param.Nickname,
		Biography: param.Biography,
	}

	if err := gorm.G[domain.User](u.db.DB).Create(ctx, &user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u *UserStore) ListUsers(ctx context.Context, params domain.ListUsersParam) (domain.Users, error) {
	db := u.db.WithContext(ctx)

	if len(params.IDs) > 0 {
		db = db.Where("id IN ?", params.IDs)
	}

	if len(params.Biographies) > 0 {
		db = db.Where("biography IN ?", params.Biographies)
	}

	for _, sort := range params.Sorts {
		db = db.Order(clause.OrderByColumn{
			Column: clause.Column{Name: sort.By},
			Desc:   sort.Order.IsDESC(),
		})
	}

	if params.Offset > 0 {
		db = db.Offset(params.Offset)
	}

	if params.Limit > 0 {
		db = db.Limit(params.Limit)
	}

	var users domain.Users
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserStore) PatchUser(ctx context.Context, params domain.PatchUserParam) error {
	user := domain.User{
		Model: gorm.Model{ID: params.ID},
	}

	if params.Nickname != nil {
		user.Nickname = *params.Nickname
	}

	if params.Biography != nil {
		user.Biography = params.Biography
	}

	if err := u.db.WithContext(ctx).Updates(user).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserStore) PatchUsers(ctx context.Context, param domain.PatchUsersParam) error {
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

	if err := u.db.WithContext(ctx).Updates(users).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserStore) DeleteUser(ctx context.Context, param domain.DeleteUserParam) error {
	db := u.db.WithContext(ctx)

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
