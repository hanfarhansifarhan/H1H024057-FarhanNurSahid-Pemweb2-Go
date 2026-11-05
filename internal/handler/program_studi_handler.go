package handler

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"latihan-fiber-api/internal/service"
)

func DaftarProgramStudi(c fiber.Ctx) error {
	daftar, err := service.DaftarProgramStudi()
	if err != nil {
		log.Printf("Gagal mengambil daftar program studi: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"data":   daftar,
	})
}
