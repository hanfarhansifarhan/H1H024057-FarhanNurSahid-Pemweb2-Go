package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

func Status(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "API Pemweb II aktif",
		"waktu":  time.Now().Format(time.RFC3339),
	})
}
