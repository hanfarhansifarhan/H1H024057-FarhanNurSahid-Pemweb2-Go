package router

import (
	"github.com/gofiber/fiber/v3"

	"latihan-fiber-api/internal/handler"
	"latihan-fiber-api/internal/middleware"
)

func Daftarkan(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/auth/registrasi", handler.Registrasi)
	api.Post("/auth/login", handler.Login)

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
	matakuliah.Post("/", handler.BuatMatakuliah)
	matakuliah.Get("/:id", handler.DetailMatakuliah)
	matakuliah.Put("/:id", handler.PerbaruiMatakuliah)
	matakuliah.Delete("/:id", handler.HapusMatakuliah)

	api.Get("/program-studi", middleware.KunciApi(), handler.DaftarProgramStudi)
}
