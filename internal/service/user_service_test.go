package service

import (
	"errors"
	"testing"

	"github.com/roka-crew/sam2ooh2-api/internal/apperr"
	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/payload"
	"github.com/roka-crew/sam2ooh2-api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	biography := "hello world"
	errDB := errors.New("db connection error")
	errUnknown := errors.New("unexpected db error")

	tests := []struct {
		name string

		request   payload.CreateUserRequest
		setupMock func(store *mocks.MockUserStore)

		expectedResponse payload.CreateUserResponse
		expectedError    error
	}{
		// ------------------------------------------------------------
		// 성공 케이스
		// ------------------------------------------------------------
		{
			name: "(1) 성공 - biography 있음",
			request: payload.CreateUserRequest{
				Nickname:  "roka",
				Biography: &biography,
			},
			setupMock: func(store *mocks.MockUserStore) {
				// 닉네임 중복 없음
				store.On("ListUsers", mock.Anything, payload.ListUsersParam{
					Nicknames: []string{"roka"},
					Limit:     1,
				}).Return(domain.Users{}, nil).Once()

				// 생성 성공
				store.On("CreateUser", mock.Anything, payload.CreateUserParam{
					Nickname:  "roka",
					Biography: &biography,
				}).Return(domain.User{
					Model:     gorm.Model{ID: 1},
					Nickname:  "roka",
					Biography: &biography,
				}, nil).Once()
			},
			expectedResponse: payload.CreateUserResponse{
				UserID:    1,
				Nickname:  "roka",
				Biography: biography,
			},
			expectedError: nil,
		},
		{
			name: "(2) 성공 - biography 없음(nil)",
			request: payload.CreateUserRequest{
				Nickname: "roka2",
			},
			setupMock: func(store *mocks.MockUserStore) {
				store.On("ListUsers", mock.Anything, payload.ListUsersParam{
					Nicknames: []string{"roka2"},
					Limit:     1,
				}).Return(domain.Users{}, nil).Once()

				store.On("CreateUser", mock.Anything, payload.CreateUserParam{
					Nickname: "roka2",
				}).Return(domain.User{
					Model:    gorm.Model{ID: 2},
					Nickname: "roka2",
				}, nil).Once()
			},
			expectedResponse: payload.CreateUserResponse{
				UserID:    2,
				Nickname:  "roka2",
				Biography: "",
			},
			expectedError: nil,
		},

		// ------------------------------------------------------------
		// 실패 케이스
		// ------------------------------------------------------------
		{
			name: "(3) 실패 - 닉네임 중복 (ListUsers에서 이미 존재)",
			request: payload.CreateUserRequest{
				Nickname: "duplicate",
			},
			setupMock: func(store *mocks.MockUserStore) {
				store.On("ListUsers", mock.Anything, payload.ListUsersParam{
					Nicknames: []string{"duplicate"},
					Limit:     1,
				}).Return(domain.Users{
					{Model: gorm.Model{ID: 99}, Nickname: "duplicate"},
				}, nil).Once()

				// CreateUser는 호출되지 않아야 하므로 On 설정을 하지 않음
			},
			expectedResponse: payload.CreateUserResponse{},
			expectedError:    apperr.ErrAlreadyExist,
		},
		{
			name: "(4) 실패 - ListUsers 조회 중 에러",
			request: payload.CreateUserRequest{
				Nickname: "roka3",
			},
			setupMock: func(store *mocks.MockUserStore) {
				store.On("ListUsers", mock.Anything, payload.ListUsersParam{
					Nicknames: []string{"roka3"},
					Limit:     1,
				}).Return(nil, errDB).Once()
			},
			expectedResponse: payload.CreateUserResponse{},
			expectedError:    errDB,
		},
		{
			name: "(5) 실패 - CreateUser에서 gorm.ErrDuplicatedKey 반환 (동시성 race)",
			request: payload.CreateUserRequest{
				Nickname: "race",
			},
			setupMock: func(store *mocks.MockUserStore) {
				store.On("ListUsers", mock.Anything, payload.ListUsersParam{
					Nicknames: []string{"race"},
					Limit:     1,
				}).Return(domain.Users{}, nil).Once()

				store.On("CreateUser", mock.Anything, payload.CreateUserParam{
					Nickname: "race",
				}).Return(domain.User{}, gorm.ErrDuplicatedKey).Once()
			},
			expectedResponse: payload.CreateUserResponse{},
			expectedError:    apperr.ErrAlreadyExist,
		},
		{
			name: "(6) 실패 - CreateUser에서 알 수 없는 에러 반환",
			request: payload.CreateUserRequest{
				Nickname: "roka4",
			},
			setupMock: func(store *mocks.MockUserStore) {
				store.On("ListUsers", mock.Anything, payload.ListUsersParam{
					Nicknames: []string{"roka4"},
					Limit:     1,
				}).Return(domain.Users{}, nil).Once()

				store.On("CreateUser", mock.Anything, payload.CreateUserParam{
					Nickname: "roka4",
				}).Return(domain.User{}, errUnknown).Once()
			},
			expectedResponse: payload.CreateUserResponse{},
			expectedError:    errUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockStore := mocks.NewMockUserStore(t)
			tt.setupMock(mockStore)
			userService := NewUserService(mockStore)

			// when
			response, err := userService.CreateUser(t.Context(), tt.request)

			// then
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResponse, response)
		})
	}
}
