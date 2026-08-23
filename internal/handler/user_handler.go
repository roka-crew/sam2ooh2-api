package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/roka-crew/sam2ooh2-api/internal/app"
	"github.com/roka-crew/sam2ooh2-api/internal/payload"
	"github.com/roka-crew/sam2ooh2-api/internal/service"
)

type UserHandler struct {
	App         *app.App
	UserService *service.UserService
}

func NewUserHandler(
	app *app.App,
	userService *service.UserService,
) *UserHandler {
	return &UserHandler{
		App:         app,
		UserService: userService,
	}
}

func UserHandlerRouteSetup(userHandler *UserHandler) {
	users := userHandler.App.Group("/users")
	{
		users.Post("/", userHandler.CreateUser)
	}
}

func (u *UserHandler) CreateUser(c fiber.Ctx) error {
	var (
		request  payload.CreateUserRequest
		response payload.CreateUserResponse
		err      error
	)

	if err = c.Bind().Body(&request); err != nil {
		return err
	}

	response, err = u.UserService.CreateUser(c.Context(), request)
	if err != nil {
		return err
	}

	return c.JSON(response)
}
