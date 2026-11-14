package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jamaalul/proyek1pbl/app/model"
)

// Success mengirimkan response sukses 200 OK dengan format standar
func Success(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created mengirimkan response sukses 201 Created
func Created(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error mengirimkan response error dengan format standar
func Error(c *fiber.Ctx, status int, message string, errDetail any) error {
	return c.Status(status).JSON(model.WebResponse{
		Success: false,
		Message: message,
		Error:   errDetail,
	})
}
