package store

import (
	"context"

	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/sqlite"
	"gorm.io/gorm"
)

type UserStore struct {
	rdb *sqlite.Sqlite
}

func NewUserStore(
	rdb *sqlite.Sqlite,
) (*UserStore, error) {
	return &UserStore{rdb: rdb}, nil
}

func (u *UserStore) CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	user := domain.User{
		Nickname:  params.Nickname,
		Biography: params.Biography,
	}

	if err := gorm.G[domain.User](u.rdb.DB).Create(ctx, &user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}
