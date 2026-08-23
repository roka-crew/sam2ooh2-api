package store

import (
	"testing"

	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/payload"
	"github.com/roka-crew/sam2ooh2-api/internal/sqlite"
	"github.com/roka-crew/sam2ooh2-api/internal/testutil/testdb"
	"github.com/roka-crew/sam2ooh2-api/pkg/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	biography := "hello world"

	tests := []struct {
		name string

		param     payload.CreateUserParam
		seedUsers []domain.User // 테스트 실행 전 미리 심어둘 데이터

		expectedError func(t *testing.T, err error)
		checkResult   func(t *testing.T, created domain.User)
	}{
		{
			// (1) 성공 - biography 있음
			name: "(1) 성공 - biography 있음",
			param: payload.CreateUserParam{
				Nickname:  "roka",
				Biography: &biography,
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, created domain.User) {
				assert.NotZero(t, created.ID)
				assert.Equal(t, "roka", created.Nickname)
				require.NotNil(t, created.Biography)
				assert.Equal(t, biography, *created.Biography)
				assert.False(t, created.CreatedAt.IsZero())
			},
		},
		{
			// (2) 성공 - biography 없음(nil)
			name: "(2) 성공 - biography 없음(nil)",
			param: payload.CreateUserParam{
				Nickname: "roka2",
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, created domain.User) {
				assert.NotZero(t, created.ID)
				assert.Equal(t, "roka2", created.Nickname)
				assert.Nil(t, created.Biography)
			},
		},
		{
			// (3) 실패 - 닉네임 unique 제약 위반
			name: "(3) 실패 - 닉네임 unique 제약 위반",
			seedUsers: []domain.User{
				{Nickname: "duplicate"},
			},
			param: payload.CreateUserParam{
				Nickname: "duplicate",
			},
			expectedError: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, gorm.ErrDuplicatedKey)
			},
			checkResult: func(t *testing.T, created domain.User) {
				assert.Zero(t, created.ID) // 생성 실패했으니 빈 값이어야 함
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db := testdb.NewTestDBSQLite(t)
			userStore := NewUserStore(db)

			for _, seed := range tt.seedUsers {
				require.NoError(t, db.Create(&seed).Error)
			}

			// when
			created, err := userStore.CreateUser(t.Context(), tt.param)

			// then
			tt.expectedError(t, err)
			tt.checkResult(t, created)
		})
	}
}

func TestListUsers(t *testing.T) {
	biography := "hi"

	tests := []struct {
		name string

		seedUsers []domain.User
		param     payload.ListUsersParam

		expectedError func(t *testing.T, err error)
		checkResult   func(t *testing.T, users domain.Users)
	}{
		{
			// (1) 닉네임으로 필터링
			name: "(1) 닉네임으로 필터링",
			seedUsers: []domain.User{
				{Nickname: "a"},
				{Nickname: "b"},
			},
			param: payload.ListUsersParam{
				Nicknames: []string{"a"},
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, users domain.Users) {
				require.Len(t, users, 1)
				assert.Equal(t, "a", users[0].Nickname)
			},
		},
		{
			// (2) biography로 필터링
			name: "(2) biography로 필터링",
			seedUsers: []domain.User{
				{Nickname: "a", Biography: &biography},
				{Nickname: "b"},
			},
			param: payload.ListUsersParam{
				Biographies: []string{biography},
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, users domain.Users) {
				require.Len(t, users, 1)
				assert.Equal(t, "a", users[0].Nickname)
			},
		},
		{
			// (3) limit/offset 적용
			name: "(3) limit/offset 적용",
			seedUsers: []domain.User{
				{Nickname: "a"},
				{Nickname: "b"},
				{Nickname: "c"},
			},
			param: payload.ListUsersParam{
				Offset: 1,
				Limit:  1,
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, users domain.Users) {
				// 시드 순서대로 id가 1,2,3이므로 offset=1,limit=1 -> 2번째 row(b) 하나
				require.Len(t, users, 1)
				assert.Equal(t, "b", users[0].Nickname)
			},
		},
		{
			// (4) 정렬 - nickname DESC
			name: "(4) 정렬 - nickname DESC",
			seedUsers: []domain.User{
				{Nickname: "a"},
				{Nickname: "c"},
				{Nickname: "b"},
			},
			param: payload.ListUsersParam{
				Sorts: query.Sorts{
					{By: "nickname", Order: query.OrderDESC},
				},
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, users domain.Users) {
				require.Len(t, users, 3)
				assert.Equal(t, []string{"c", "b", "a"}, []string{
					users[0].Nickname, users[1].Nickname, users[2].Nickname,
				})
			},
		},
		{
			// (5) 조건에 맞는 데이터 없음
			name: "(5) 조건에 맞는 데이터 없음",
			seedUsers: []domain.User{
				{Nickname: "a"},
			},
			param: payload.ListUsersParam{
				Nicknames: []string{"nonexistent"},
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, users domain.Users) {
				assert.Empty(t, users)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db := testdb.NewTestDBSQLite(t)
			userStore := NewUserStore(db)

			for _, seed := range tt.seedUsers {
				require.NoError(t, db.Create(&seed).Error)
			}

			// when
			users, err := userStore.ListUsers(t.Context(), tt.param)

			// then
			tt.expectedError(t, err)
			tt.checkResult(t, users)
		})
	}
}
func TestPatchUser(t *testing.T) {
	newNickname := "updated"
	newBiography := "updated bio"

	tests := []struct {
		name string

		seedUser domain.User
		param    func(seededID uint) payload.PatchUserParam // seed 생성 후 나온 ID를 param에 반영하기 위해 함수로

		expectedError func(t *testing.T, err error)
		checkResult   func(t *testing.T, db *sqlite.Sqlite, seededID uint)
	}{
		{
			// (1) 닉네임만 변경
			name:     "(1) 닉네임만 변경",
			seedUser: domain.User{Nickname: "before"},
			param: func(seededID uint) payload.PatchUserParam {
				return payload.PatchUserParam{
					ID:       seededID,
					Nickname: &newNickname,
				}
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, db *sqlite.Sqlite, seededID uint) {
				var got domain.User
				require.NoError(t, db.First(&got, seededID).Error)
				assert.Equal(t, "updated", got.Nickname)
			},
		},
		{
			// (2) biography만 변경 (nil -> 값 있음)
			name:     "(2) biography만 변경",
			seedUser: domain.User{Nickname: "roka"},
			param: func(seededID uint) payload.PatchUserParam {
				return payload.PatchUserParam{
					ID:        seededID,
					Biography: &newBiography,
				}
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, db *sqlite.Sqlite, seededID uint) {
				var got domain.User
				require.NoError(t, db.First(&got, seededID).Error)
				assert.Equal(t, "roka", got.Nickname) // 변경 안 됨
				require.NotNil(t, got.Biography)
				assert.Equal(t, newBiography, *got.Biography)
			},
		},
		{
			// (3) 변경할 필드가 없으면 아무 것도 하지 않음(no-op)
			name:     "(3) 변경 필드 없음 - no-op",
			seedUser: domain.User{Nickname: "unchanged"},
			param: func(seededID uint) payload.PatchUserParam {
				return payload.PatchUserParam{ID: seededID}
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, db *sqlite.Sqlite, seededID uint) {
				var got domain.User
				require.NoError(t, db.First(&got, seededID).Error)
				assert.Equal(t, "unchanged", got.Nickname)
			},
		},
		{
			// (4) 존재하지 않는 ID -> 에러 없이 그냥 영향 row 0건
			name:     "(4) 존재하지 않는 ID",
			seedUser: domain.User{Nickname: "someone"},
			param: func(seededID uint) payload.PatchUserParam {
				return payload.PatchUserParam{
					ID:       seededID + 999,
					Nickname: &newNickname,
				}
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err) // gorm Updates는 매칭 row 없어도 에러 아님
			},
			checkResult: func(t *testing.T, db *sqlite.Sqlite, seededID uint) {
				var got domain.User
				require.NoError(t, db.First(&got, seededID).Error)
				assert.Equal(t, "someone", got.Nickname) // 원본 그대로
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db := testdb.NewTestDBSQLite(t)
			userStore := NewUserStore(db)

			require.NoError(t, db.Create(&tt.seedUser).Error)

			// when
			err := userStore.PatchUser(t.Context(), tt.param(tt.seedUser.ID))

			// then
			tt.expectedError(t, err)
			tt.checkResult(t, db, tt.seedUser.ID)
		})
	}
}

func TestPatchUsers(t *testing.T) {
	nickA := "new-a"
	nickB := "new-b"

	t.Run("(1) 여러 명 동시에 수정", func(t *testing.T) {
		// given
		db := testdb.NewTestDBSQLite(t)
		userStore := NewUserStore(db)

		userA := domain.User{Nickname: "a"}
		userB := domain.User{Nickname: "b"}
		require.NoError(t, db.Create(&userA).Error)
		require.NoError(t, db.Create(&userB).Error)

		// when
		err := userStore.PatchUsers(t.Context(), payload.PatchUsersParam{
			{ID: userA.ID, Nickname: &nickA},
			{ID: userB.ID, Nickname: &nickB},
		})

		// then
		require.NoError(t, err)

		var gotA, gotB domain.User
		require.NoError(t, db.First(&gotA, userA.ID).Error)
		require.NoError(t, db.First(&gotB, userB.ID).Error)
		assert.Equal(t, "new-a", gotA.Nickname)
		assert.Equal(t, "new-b", gotB.Nickname)
	})

	t.Run("(2) 중간에 실패하면 전체 롤백", func(t *testing.T) {
		// given
		db := testdb.NewTestDBSQLite(t)
		userStore := NewUserStore(db)

		userA := domain.User{Nickname: "a"}
		userB := domain.User{Nickname: "existing-nickname"} // 아래 duplicate 유발용
		require.NoError(t, db.Create(&userA).Error)
		require.NoError(t, db.Create(&userB).Error)

		duplicateNickname := "existing-nickname" // userB와 동일 -> unique 위반 유발

		// when: userA를 duplicateNickname으로 바꾸려다 실패해야 함
		err := userStore.PatchUsers(t.Context(), payload.PatchUsersParam{
			{ID: userA.ID, Nickname: &duplicateNickname},
		})

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrDuplicatedKey)

		// 롤백 확인: userA는 여전히 원래 닉네임
		var gotA domain.User
		require.NoError(t, db.First(&gotA, userA.ID).Error)
		assert.Equal(t, "a", gotA.Nickname)
	})
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name string

		seedUsers []domain.User
		param     func(seededIDs []uint) payload.DeleteUserParam

		expectedError func(t *testing.T, err error)
		checkResult   func(t *testing.T, db *sqlite.Sqlite, seededIDs []uint)
	}{
		{
			// (1) ID로 soft delete
			name: "(1) ID로 soft delete",
			seedUsers: []domain.User{
				{Nickname: "a"},
			},
			param: func(seededIDs []uint) payload.DeleteUserParam {
				return payload.DeleteUserParam{ID: seededIDs[0]}
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, db *sqlite.Sqlite, seededIDs []uint) {
				// soft delete -> 기본 쿼리로는 안 보임
				var got domain.User
				err := db.First(&got, seededIDs[0]).Error
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

				// Unscoped로 보면 DeletedAt이 채워져 있어야 함
				var gotUnscoped domain.User
				require.NoError(t, db.Unscoped().First(&gotUnscoped, seededIDs[0]).Error)
				assert.True(t, gotUnscoped.DeletedAt.Valid)
			},
		},
		{
			// (2) 닉네임으로 soft delete
			name: "(2) 닉네임으로 soft delete",
			seedUsers: []domain.User{
				{Nickname: "target"},
			},
			param: func(seededIDs []uint) payload.DeleteUserParam {
				return payload.DeleteUserParam{Nickname: "target"}
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, db *sqlite.Sqlite, seededIDs []uint) {
				var got domain.User
				err := db.First(&got, seededIDs[0]).Error
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			},
		},
		{
			// (3) hard delete
			name: "(3) hard delete",
			seedUsers: []domain.User{
				{Nickname: "hard"},
			},
			param: func(seededIDs []uint) payload.DeleteUserParam {
				return payload.DeleteUserParam{
					ID:           seededIDs[0],
					IsHardDelete: true,
				}
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, db *sqlite.Sqlite, seededIDs []uint) {
				// Unscoped로 봐도 완전히 사라져야 함
				var got domain.User
				err := db.Unscoped().First(&got, seededIDs[0]).Error
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db := testdb.NewTestDBSQLite(t)
			userStore := NewUserStore(db)

			var seededIDs []uint
			for _, seed := range tt.seedUsers {
				require.NoError(t, db.Create(&seed).Error)
				seededIDs = append(seededIDs, seed.ID)
			}

			// when
			err := userStore.DeleteUser(t.Context(), tt.param(seededIDs))

			// then
			tt.expectedError(t, err)
			tt.checkResult(t, db, seededIDs)
		})
	}
}
