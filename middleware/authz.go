package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jamaalul/proyek1pbl/helper"
)

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := helper.CurrentUser(c)
		if !ok {
			return helper.Unauthorized("sesi autentikasi diperlukan")
		}

		if !helper.HasAnyRole(user.Role, roles...) {
			return helper.Forbidden("role Anda tidak memiliki hak akses ke resource ini")
		}

		return c.Next()
	}
}

func RequirePermission(perms *helper.PermissionSet, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := helper.CurrentUser(c)
		if !ok {
			return helper.Unauthorized("sesi autentikasi diperlukan")
		}

		if perms == nil || !perms.Can(user.Role, permission) {
			return helper.Forbidden("Anda tidak memiliki izin yang diperlukan: " + permission)
		}

		return c.Next()
	}
}
