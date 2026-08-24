package store

import (
	"testing"

	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/payload"
	"github.com/roka-crew/sam2ooh2-api/internal/sqlite"
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

func TestListGroups(t *testing.T) {
	// 테스트용 기본 시드 데이터
	seedGroups := []domain.Group{
		{
			Title:       "클린 코드",
			Author:      "로버트 C. 마틴",
			PageCount:   464,
			Publisher:   "인사이트",
			Description: "가독성 좋은 코드를 작성하는 방법",
		},
		{
			Title:       "리팩터링",
			Author:      "마틴 파울러",
			PageCount:   500,
			Publisher:   "한빛미디어",
			Description: "코드 구조를 개선하는 아키텍처",
		},
		{
			Title:       "도메인 주도 설계",
			Author:      "에릭 반스",
			PageCount:   500,
			Publisher:   "인사이트",
			Description: "복잡한 소프트웨어를 다루는 아키텍처",
		},
	}

	tests := []struct {
		name string

		param payload.ListGroupsParam

		expectedError func(t *testing.T, err error)
		checkResult   func(t *testing.T, groups domain.Groups)
	}{
		{
			name:  "(1) 조건 없이 전체 조회",
			param: payload.ListGroupsParam{},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, groups domain.Groups) {
				assert.Len(t, groups, 3)
			},
		},
		{
			name: "(2) Publisher 단일 조건 필터링",
			param: payload.ListGroupsParam{
				Publishers: []string{"인사이트"},
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, groups domain.Groups) {
				assert.Len(t, groups, 2)
				for _, g := range groups {
					assert.Equal(t, "인사이트", g.Publisher)
				}
			},
		},
		{
			name: "(3) PageCount 및 Publisher 복합 조건 필터링",
			param: payload.ListGroupsParam{
				PageCounts: []int{500},
				Publishers: []string{"인사이트"},
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, groups domain.Groups) {
				require.Len(t, groups, 1)
				assert.Equal(t, "도메인 주도 설계", groups[0].Title)
				assert.Equal(t, "에릭 반스", groups[0].Author)
			},
		},
		{
			name: "(4) 페이징 (Limit & Offset) 적용",
			param: payload.ListGroupsParam{
				Limit:  1,
				Offset: 1,
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, groups domain.Groups) {
				require.Len(t, groups, 1)
				assert.Equal(t, "리팩터링", groups[0].Title)
			},
		},
		{
			name: "(5) 일치하는 조건이 없어 빈 배열 반환",
			param: payload.ListGroupsParam{
				Authors: []string{"존재하지 않는 저자"},
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, groups domain.Groups) {
				assert.Empty(t, groups)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db := testdb.NewTestDBSQLite(t)
			groupStore := NewGroupStore(db)

			// 시드 데이터 삽입
			for i := range seedGroups {
				require.NoError(t, db.Create(&seedGroups[i]).Error)
			}

			// when
			groups, err := groupStore.ListGroups(t.Context(), tt.param)

			// then
			tt.expectedError(t, err)
			tt.checkResult(t, groups)
		})
	}
}

func TestPatchGroup(t *testing.T) {
	tests := []struct {
		name string

		seedGroup domain.Group
		param     payload.PatchGroupParam

		expectedError func(t *testing.T, err error)
		checkResult   func(t *testing.T, db *sqlite.Sqlite, groupID uint)
	}{
		{
			name: "(1) 성공 - Title 및 Author 부분 수정",
			seedGroup: domain.Group{
				Title:     "클린 코드",
				Author:    "로버트 C. 마틴",
				PageCount: 464,
				Publisher: "인사이트",
			},
			param: payload.PatchGroupParam{
				Title:  "클린 코드 (개정판)",
				Author: "Uncle Bob",
			},
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, db *sqlite.Sqlite, groupID uint) {
				var updated domain.Group
				require.NoError(t, db.First(&updated, groupID).Error)

				assert.Equal(t, "클린 코드 (개정판)", updated.Title)
				assert.Equal(t, "Uncle Bob", updated.Author)
				assert.Equal(t, 464, updated.PageCount)    // 기존 값 유지 확인
				assert.Equal(t, "인사이트", updated.Publisher) // 기존 값 유지 확인
			},
		},
		{
			name: "(2) 성공 - 수정할 필드가 없는 경우 (Zero Value)",
			seedGroup: domain.Group{
				Title:     "변경 없음",
				PageCount: 100,
			},
			param: payload.PatchGroupParam{}, // 모든 필드가 Zero Value
			expectedError: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
			checkResult: func(t *testing.T, db *sqlite.Sqlite, groupID uint) {
				var current domain.Group
				require.NoError(t, db.First(&current, groupID).Error)

				assert.Equal(t, "변경 없음", current.Title)
				assert.Equal(t, 100, current.PageCount)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			db := testdb.NewTestDBSQLite(t)
			groupStore := NewGroupStore(db)

			require.NoError(t, db.Create(&tt.seedGroup).Error)
			tt.param.ID = tt.seedGroup.ID // 생성된 시드 데이터의 ID 동적 설정

			// when
			err := groupStore.PatchGroup(t.Context(), tt.param)

			// then
			tt.expectedError(t, err)
			tt.checkResult(t, db, tt.seedGroup.ID)
		})
	}
}
