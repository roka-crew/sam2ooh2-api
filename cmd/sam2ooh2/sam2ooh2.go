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

	var (
		storeModule = fx.Module("store",
			fx.Provide(
				store.NewUserStore,
			),
		)

		serviceModule = fx.Module("service",
			fx.Provide(
				service.NewUserService,
			),
		)

		handlerModule = fx.Module("handler",
			fx.Provide(
				handler.NewUserHandler,
			),
			fx.Invoke(
				handler.UserHandlerRouteSetup,
			),
		)
	)

	fx.New(
		fx.Supply("./configs/sam2ooh2.yaml"),
		fx.Provide(
			config.NewConfig,
			sqlite.NewSqlite,
			app.NewApp,
		),
		storeModule,
		serviceModule,
		handlerModule,
	).Run()
}
