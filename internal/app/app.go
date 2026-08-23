package app

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/roka-crew/sam2ooh2-api/pkg/config"
	"go.uber.org/fx"
)

type App struct {
	*fiber.App
}

func NewApp(
	cfg *config.Config,
	lifeCycle fx.Lifecycle,
) *App {
	structValidator, err := newStructValidator()
	if err != nil {
		log.Panicf("failed to new struct validator: %+v\n", err)
	}

	app := &App{
		App: fiber.New(fiber.Config{
			AppName:         cfg.Name,
			StructValidator: structValidator,
			ErrorHandler:    errorHandler,
		}),
	}

	lifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(cfg.Listen); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	app.Use(recover.New())
	if cfg.Env == config.EnvDev {
		app.Use(requestid.New()) // Ensure requestid middleware is used before the logger
		app.Use(logger.New(logger.Config{
			// 시간 | PID | ReqID | 상태코드 | 처리시간 | 메서드 | 경로 - 에러메시지
			Format: "[${time}] ${pid} | ${requestid} | ${status} | ${latency} | ${method} ${path} ${error}\n",
			// TimeZone을 UTC로 지정
			TimeZone: "UTC",
			// ISO-8601 UTC 표준 포맷
			TimeFormat: "2006-01-02T15:04:05Z",
		}))
	}

	return app
}
