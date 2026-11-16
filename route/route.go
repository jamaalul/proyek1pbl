package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jamaalul/proyek1pbl/helper"
)

type Dependencies struct {
	Pool        *pgxpool.Pool
	JWT         *helper.JWTManager
	Permissions *helper.PermissionSet
}

func Register(app *fiber.App, deps Dependencies) {
	app.Get("/", func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.StatusOK, "Boilerplate berjalan", fiber.Map{
			"version": "1.0.0",
			"health":  "/health",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "ok"
		if deps.Pool == nil {
			dbStatus = "fallback (offline)"
		} else {
			if err := deps.Pool.Ping(c.Context()); err != nil {
				dbStatus = "error: " + err.Error()
			}
		}

		return helper.Success(c, fiber.StatusOK, "server dan basis data berjalan normal", fiber.Map{
			"postgres": dbStatus,
		})
	})
}
