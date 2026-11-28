package router

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"

	"latihan-fiber-api/internal/handler"
	"latihan-fiber-api/internal/middleware"
)

func Daftarkan(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/status", handler.Status)
	api.Post("/auth/registrasi", handler.Registrasi)
	api.Post("/auth/login", limiter.New(limiter.Config{
		Max:        5,
		Expiration: time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"sukses": false,
				"pesan":  "Terlalu banyak percobaan login",
			})
		},
	}), handler.Login)

	authBerlindung := api.Group("/auth", middleware.Autentikasi())
	authBerlindung.Get("/profil", handler.Profil)
	authBerlindung.Post("/segarkan", handler.SegarkanToken)
	authBerlindung.Post("/logout", handler.Logout)

	mahasiswa := api.Group("/mahasiswa", middleware.Autentikasi())
	mahasiswa.Get("/", handler.DaftarMahasiswa)
	mahasiswa.Post("/", middleware.HanyaPeran("admin", "dosen"), handler.BuatMahasiswa)
	mahasiswa.Get("/:id", handler.DetailMahasiswa)
	mahasiswa.Put("/:id", middleware.HanyaPeran("admin", "dosen"), handler.PerbaruiMahasiswa)
	mahasiswa.Delete("/:id", middleware.HanyaPeran("admin"), handler.HapusMahasiswa)
	mahasiswa.Post("/:id/restore", middleware.HanyaPeran("admin", "dosen"), handler.PulihkanMahasiswa)

	matakuliah := api.Group("/matakuliah", middleware.KunciApi())
	matakuliah.Get("/", handler.DaftarMatakuliah)
	matakuliah.Get("/:id", handler.DetailMatakuliah)
	matakuliah.Post("/", middleware.Autentikasi(), middleware.HanyaPeran("admin", "dosen"), handler.BuatMatakuliah)
	matakuliah.Put("/:id", middleware.Autentikasi(), middleware.HanyaPeran("admin", "dosen"), handler.PerbaruiMatakuliah)
	matakuliah.Delete("/:id", middleware.Autentikasi(), middleware.HanyaPeran("admin"), handler.HapusMatakuliah)

	api.Get("/program-studi", middleware.KunciApi(), handler.DaftarProgramStudi)
}
