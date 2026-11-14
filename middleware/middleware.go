package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jamaalul/proyek1pbl/helper"
)

func Register(app *fiber.App) {
	app.Use(requestid.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		MaxAge:       86400,
	}))
	app.Use(RequestLogger())
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		reqID, _ := c.Locals("requestid").(string)

		log.Printf("[%s] %s %s | %d | %v | %s\n",
			c.Method(), c.Path(), c.IP(), c.Response().StatusCode(), duration, reqID,
		)
		return err
	}
}

func RequireJSON(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodPost || c.Method() == fiber.MethodPut || c.Method() == fiber.MethodPatch {
		ct := string(c.Request().Header.ContentType())
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return helper.BadRequest("Content-Type harus application/json")
		}
	}
	return c.Next()
}
