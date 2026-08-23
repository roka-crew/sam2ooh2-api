package sqlite

import (
	"fmt"
	"time"

	"github.com/roka-crew/sam2ooh2-api/pkg/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Sqlite struct {
	*gorm.DB
}

func NewSqlite(cfg *config.Config) (*Sqlite, error) {
	format := "%s.db"
	dsn := fmt.Sprintf(format,
		cfg.Sqlite.DBname,
	)

	var gormConfigLogger logger.Interface
	if cfg.Env == config.EnvDev {
		gormConfigLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         gormConfigLogger,
		NowFunc:        func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.Sqlite.Options.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Sqlite.Options.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Sqlite.Options.ConnMaxLifetime)

	return &Sqlite{DB: db}, nil
}
