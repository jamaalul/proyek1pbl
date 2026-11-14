package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jamaalul/proyek1pbl/config"
	"github.com/jamaalul/proyek1pbl/helper"
)

// Protected memastikan request membawa JWT Bearer token yang valid
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helper.Error(c, fiber.StatusUnauthorized, "Header Authorization diperlukan", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return helper.Error(c, fiber.StatusUnauthorized, "Format token harus: Bearer <token>", nil)
		}

		tokenString := parts[1]
		secret := config.GetEnv("JWT_SECRET", "supersecretmedicarejwtkey2026ti")

		claims, err := helper.ValidateToken(tokenString, secret)
		if err != nil {
			return helper.Error(c, fiber.StatusUnauthorized, "Token tidak valid atau telah kedaluwarsa", err.Error())
		}

		// Simpan klaim data user ke context lokal Fiber
		c.Locals("user", claims)
		return c.Next()
	}
}
