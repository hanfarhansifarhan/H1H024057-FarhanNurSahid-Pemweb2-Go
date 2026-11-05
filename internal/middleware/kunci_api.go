package middleware

import (
	"os"

	"github.com/gofiber/fiber/v3"
)

func KunciApi() fiber.Handler {
	return func(c fiber.Ctx) error {
		kunciBenar := os.Getenv("API_KEY")
		kunciDikirim := c.Get("X-Api-Key")

		if kunciBenar == "" || kunciDikirim != kunciBenar {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"sukses": false,
				"pesan":  "API key tidak valid atau tidak dikirim",
			})
		}

		return c.Next()
	}
}
