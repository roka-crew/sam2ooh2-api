package testdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/sqlite"
	"github.com/stretchr/testify/require"
	gormsqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewTestDBSQLite(t *testing.T) *sqlite.Sqlite {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.NewString())

	db, err := gorm.Open(gormsqlite.Open(dsn), &gorm.Config{
		TranslateError: true, // 프로덕션(sqlite.NewSqlite)과 동일하게 설정
		NowFunc:        func() time.Time { return time.Now().UTC() },
	})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(
		// 앞으로 추가되는 도메인 모델은 여기 한 곳에만 등록
		&domain.User{},
		&domain.Group{},
		&domain.Chapter{},
		&domain.Topic{},
	))

	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1) // in-memory DB가 커넥션마다 분리되는 것을 방지

	t.Cleanup(func() {
		sqlDB.Close()
	})

	return &sqlite.Sqlite{DB: db}
}
