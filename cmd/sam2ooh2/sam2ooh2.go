package main

import (
	"github.com/roka-crew/sam2ooh2-api/internal/app"
	"github.com/roka-crew/sam2ooh2-api/internal/handler"
	"github.com/roka-crew/sam2ooh2-api/internal/service"
	"github.com/roka-crew/sam2ooh2-api/internal/sqlite"
	"github.com/roka-crew/sam2ooh2-api/internal/store"
	"github.com/roka-crew/sam2ooh2-api/pkg/config"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Supply("./configs/sam2ooh2.yaml"),
		fx.Provide(
			config.NewConfig,
			sqlite.NewSqlite,

			store.NewUserStore,

			service.NewUserService,

			handler.NewUserHandler,

			app.NewApp,
		),
		fx.Invoke(
			handler.UserHandlerRouteSetup,
		),
	).Run()
}
