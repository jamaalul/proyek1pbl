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

func ParamID64(c *fiber.Ctx) (int64, bool) {
	raw := c.Params("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func ParseAndValidate(c *fiber.Ctx, s any) error {
	if err := c.BodyParser(s); err != nil {
		return ErrBadRequest("BAD_REQUEST", "format JSON tidak valid")
	}
	if errs := ValidateStruct(s); len(errs) > 0 {
		return Validation(errs)
	}
	return nil
}
