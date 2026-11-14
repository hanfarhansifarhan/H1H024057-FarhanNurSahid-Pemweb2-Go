package middleware

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"latihan-fiber-api/internal/auth"
)

func Autentikasi() fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"sukses": false,
				"pesan":  "Token tidak valid atau tidak dikirim",
			})
		}

		bagian := strings.Fields(header)
		if len(bagian) != 2 || !strings.EqualFold(bagian[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"sukses": false,
				"pesan":  "Format Authorization tidak valid",
			})
		}

		klaim, err := auth.PeriksaToken(bagian[1])
		if err != nil || klaim == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"sukses": false,
				"pesan":  "Token tidak valid atau kedaluwarsa",
			})
		}

		sudahDicabut, err := auth.SudahDiBlacklist(klaim.JTI)
		if err != nil {
			log.Printf("Gagal memeriksa blacklist token: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Gagal memeriksa autentikasi")
		}
		if sudahDicabut {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"sukses": false,
				"pesan":  "Token tidak valid atau dicabut",
			})
		}

		idPengguna, err := strconv.ParseUint(klaim.Subject, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"sukses": false,
				"pesan":  "Token tidak valid",
			})
		}

		c.Locals("id_pengguna", uint(idPengguna))
		c.Locals("peran", klaim.Peran)
		return c.Next()
	}
}

func HanyaPeran(peranDibolehkan ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		peran, ok := c.Locals("peran").(string)
		if !ok || peran == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"sukses": false,
				"pesan":  "Akses ditolak",
			})
		}

		for _, peranIzin := range peranDibolehkan {
			if strings.EqualFold(peran, peranIzin) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Akses ditolak",
		})
	}
}
