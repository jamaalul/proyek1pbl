package config

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jamaalul/proyek1pbl/app/model"
	"github.com/jamaalul/proyek1pbl/helper"
	"github.com/jamaalul/proyek1pbl/middleware"
)

func NewApp() *fiber.App {
	appName := GetEnv("APP_NAME", "4 Sehat 5 Haji")

	app := fiber.New(fiber.Config{
		AppName:      appName,
		ErrorHandler: newErrorHandler(),
	})

	middleware.Register(app)
	return app
}

func newErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		statusCode := fiber.StatusInternalServerError
		message := "terjadi kesalahan pada server"
		var errorsList []model.ErrorItem

		var appErr *helper.AppError
		if errors.As(err, &appErr) {
			statusCode = appErr.Status
			message = appErr.Message
			errorsList = appErr.Errors
		} else {
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				statusCode = fiberErr.Code
				message = fiberErr.Message
				errorsList = []model.ErrorItem{
					{Code: "BAD_REQUEST", Message: fiberErr.Message},
				}
			} else {
				// Log error internal tak terduga
				log.Printf("[SERVER_ERROR] %s %s: %v\n", c.Method(), c.Path(), err)
				if GetEnv("APP_ENV", "development") != "production" {
					message = err.Error()
				}
				errorsList = []model.ErrorItem{
					{Code: "INTERNAL", Message: message},
				}
			}
		}

		return c.Status(statusCode).JSON(model.ErrorResponse{
			Success: false,
			Message: message,
			Errors:  errorsList,
		})
	}
}
