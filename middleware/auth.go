package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/jamaalul/proyek1pbl/helper"
)

func RequireAuth(jwtManager *helper.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helper.Unauthorized("header Authorization diperlukan")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return helper.Unauthorized("format authorization harus 'Bearer <token>'")
		}

		tokenString := strings.TrimSpace(parts[1])
		authUser, err := jwtManager.Parse(tokenString)
		if err != nil {
			return helper.Unauthorized(err.Error())
		}

		c.Locals(helper.LocalsAuthUser, authUser)
		return c.Next()
	}
}

func LoginRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return helper.TooManyRequests("terlalu banyak percobaan login, coba lagi setelah 1 menit")
		},
		SkipSuccessfulRequests: true,
	})
}
