package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jamaalul/proyek1pbl/helper"
)

// RequireRoles membatasi akses endpoint hanya untuk role tertentu
func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user").(*helper.JWTClaims)
		if !ok || claims == nil {
			return helper.Error(c, fiber.StatusUnauthorized, "Sesi autentikasi tidak ditemukan", nil)
		}

		userRole := claims.Role
		for _, allowedRole := range roles {
			if userRole == allowedRole {
				return c.Next()
			}
		}

		return helper.Error(c, fiber.StatusForbidden, "Akses ditolak: role Anda tidak memiliki izin untuk fitur ini", nil)
	}
}
