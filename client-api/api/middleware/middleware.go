package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RegisterMiddlewares(app *fiber.App, rateLimit int64) {
	app.Use(helmet.New())

	app.Use(limiter.New(limiter.Config{
		Max:        int(rateLimit),
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "too many requests"})
		},
	}))
}
