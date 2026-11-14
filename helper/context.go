package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jamaalul/proyek1pbl/app/model"
)

const LocalsAuthUser = "authUser"

// CurrentUser mengambil data pengguna yang sedang login dari Locals context Fiber
func CurrentUser(c *fiber.Ctx) (model.AuthUser, bool) {
	val := c.Locals(LocalsAuthUser)
	if val == nil {
		return model.AuthUser{}, false
	}
	user, ok := val.(model.AuthUser)
	return user, ok
}
