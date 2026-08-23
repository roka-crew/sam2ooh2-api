package app

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/roka-crew/sam2ooh2-api/internal/apperr"
)

type ValidationErrorDetail struct {
	Field  string `json:"field"`  // e.g. "Nickname"
	Reason string `json:"reason"` // e.g. "Nickname must be at most 10 characters in length"
}

func errorHandler(c fiber.Ctx, err error) error {
	var appError *apperr.AppError
	if errors.As(err, &appError) && appError != nil {
		return c.Status(appError.StatusCode).JSON(appError)
	}

	var internalError *apperr.InternalError
	if errors.As(err, &internalError) && internalError != nil {
		fmt.Println("internal error: ", err)
		fmt.Println("StackTrace()")
		fmt.Println(internalError.StackTrace(func(file, _ string, line int) string {
			return fmt.Sprintf("\t%s:%d", file, line)
		}))

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return fiber.DefaultErrorHandler(c, err)
}
