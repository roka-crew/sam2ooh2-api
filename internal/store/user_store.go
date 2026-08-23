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

func (u *UserStore) CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	user := domain.User{
		Nickname:  params.Nickname,
		Biography: params.Biography,
	}

	if err := gorm.G[domain.User](u.db.DB).Create(ctx, &user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u *UserStore) ListUsers(ctx context.Context, params domain.ListUsersParams) (domain.Users, error) {
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
