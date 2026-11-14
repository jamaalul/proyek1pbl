package helper

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const DefaultDatabaseTimeout = 5 * time.Second

func RequestContext(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), DefaultDatabaseTimeout)
}

func ParamID(c *fiber.Ctx) (int, bool) {
	raw := c.Params("id")
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
