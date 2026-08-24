package store

import (
	"testing"

	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/payload"
	"github.com/roka-crew/sam2ooh2-api/internal/testutil/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
	tests := []struct {
		name string

		param payload.CreateGroupParam

		expectedError func(t *testing.T, err error)
		checkResult   func(t *testing.T, created domain.Group)
	}{
		{
			// (1) 성공 - 정상 데이터 입력
			name: "(1) 성공 - 정상 데이터 입력",
			param: payload.CreateGroupParam{
				Title:       "클린 코드 독서회",
				Author:      "로버트 C. 마틴",
				PageCount:   464,
				Publisher:   "인사이트",
				Description: "클린 코드를 함께 읽고 토론하는 모임입니다.",
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, created domain.Group) {
				assert.NotZero(t, created.ID)
				assert.Equal(t, "클린 코드 독서회", created.Title)
				assert.Equal(t, "로버트 C. 마틴", created.Author)
				assert.Equal(t, 464, created.PageCount)
				assert.Equal(t, "인사이트", created.Publisher)
				assert.Equal(t, "클린 코드를 함께 읽고 토론하는 모임입니다.", created.Description)
				assert.False(t, created.CreatedAt.IsZero())
			},
		},
		{
			// (2) 실패 - PageCount 0 이하 (DB check:page_count > 0 제약조건 위반)
			name: "(2) 실패 - PageCount 0 이하 (check 제약조건 위반)",
			param: payload.CreateGroupParam{
				Title:       "페이지 수 오류 그룹",
				Author:      "테스터",
				PageCount:   0, // check 위반
				Publisher:   "테스트 출판사",
				Description: "잘못된 페이지 수입니다.",
			},
			expectedError: func(t *testing.T, err error) {
				require.Error(t, err)
			},
			checkResult: func(t *testing.T, created domain.Group) {
				assert.Zero(t, created.ID) // 생성 실패했으므로 빈 값이어야 함
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db := testdb.NewTestDBSQLite(t)
			groupStore := NewGroupStore(db)

			// when
			created, err := groupStore.CreateGroup(t.Context(), tt.param)

			// then
			tt.expectedError(t, err)
			tt.checkResult(t, created)
		})
	}
}
