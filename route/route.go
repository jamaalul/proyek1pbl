package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jamaalul/proyek1pbl/app/service"
	"github.com/jamaalul/proyek1pbl/helper"
	"github.com/jamaalul/proyek1pbl/middleware"
)

// Dependencies menyimpan daftar service yang akan dihubungkan ke route
type Dependencies struct {
	Pool        *pgxpool.Pool
	AuthService *service.AuthService
}

// Register mendaftarkan seluruh endpoint API ke router Fiber
func Register(app *fiber.App, deps Dependencies) {
	// Root Endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.StatusOK, "Selamat datang di MediCare Backend API", fiber.Map{
			"service": "MediCare API Service",
			"version": "1.0.0",
			"docs":    "/api/v1/health",
		})
	})

	api := app.Group("/api/v1")

	// Health Check Endpoint untuk kebutuhan pengujian & demo
	api.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "connected"
		if deps.Pool == nil {
			dbStatus = "memory-fallback (database offline)"
		} else {
			if err := deps.Pool.Ping(c.Context()); err != nil {
				dbStatus = "disconnected: " + err.Error()
			}
		}

		return helper.Success(c, fiber.StatusOK, "Sistem MediCare berjalan dengan baik", fiber.Map{
			"status":   "UP",
			"database": dbStatus,
			"modules": fiber.Map{
				"core_auth":   "ready",
				"clinic":      "boilerplate ready",
				"transaction": "boilerplate ready",
			},
		})
	})

	// Modul Autentikasi & Akun Pengguna
	if deps.AuthService != nil {
		auth := api.Group("/auth")
		auth.Post("/register", deps.AuthService.Register)
		auth.Post("/login", deps.AuthService.Login)
		auth.Get("/me", middleware.Protected(), deps.AuthService.Profile)
	}
}
