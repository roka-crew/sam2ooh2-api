package main

import (
	"log"

	"github.com/roka-crew/sam2ooh2-api/internal/domain"
	"github.com/roka-crew/sam2ooh2-api/internal/sqlite"
	"github.com/roka-crew/sam2ooh2-api/pkg/config"
)

func main() {
	cfg, err := config.NewConfig("./configs/sam2ooh2.yaml")
	if err != nil {
		log.Panicf("failed to new config: %v\n", err)
	}

	sqliteDB, err := sqlite.NewSqlite(cfg)
	if err != nil {
		log.Panicf("failed to new sqlite: %v\n", err)
	}

	if err := sqliteDB.AutoMigrate(new(domain.User)); err != nil {
		log.Panicf("failed to auto migrate: %v\n", err)
	}
}
