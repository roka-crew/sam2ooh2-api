package service

import (
	"context"
	"errors"

	"github.com/roka-crew/sam2ooh2-api/internal/apperr"
	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/payload"
	"gorm.io/gorm"
)

type UserService struct {
	userStore domain.UserStore
}

func NewUserService(
	userStore domain.UserStore,
) domain.UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (s *UserService) CreateUser(c context.Context, request payload.CreateUserRequest) (payload.CreateUserResponse, error) {
	// (1) 사용자의 Nickname 중복 검사
	listUsersByNickname, err := s.userStore.ListUsers(c, payload.ListUsersParam{
		Nicknames: []string{request.Nickname},
		Limit:     1,
	})
	if err != nil {
		return payload.CreateUserResponse{}, err
	}
	if len(listUsersByNickname) > 0 {
		return payload.CreateUserResponse{}, apperr.ErrAlreadyExist
	}

	// (2) 사용자 생성
	createdUser, err := s.userStore.CreateUser(c, payload.CreateUserParam{
		Nickname:  request.Nickname,
		Biography: request.Biography,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return payload.CreateUserResponse{}, apperr.ErrAlreadyExist
		}

		return payload.CreateUserResponse{}, err
	}

	// (3) 반환
	var biography string
	if createdUser.Biography != nil {
		biography = *createdUser.Biography
	}
	return payload.CreateUserResponse{
		UserID:    createdUser.ID,
		Nickname:  createdUser.Nickname,
		Biography: biography,
	}, nil
}
