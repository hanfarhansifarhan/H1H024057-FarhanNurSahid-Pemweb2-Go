package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"latihan-fiber-api/internal/config"
	"latihan-fiber-api/internal/middleware"
	"latihan-fiber-api/internal/router"
)

func main() {
	konf := config.Muat()

	corsHandler, err := buatCorsMiddleware()
	if err != nil {
		log.Fatal("Konfigurasi CORS tidak valid: ", err)
	}

	if err := config.HubungkanBasisData(konf); err != nil {
		log.Fatal("Gagal terhubung ke basis data: ", err)
	}
	if err := config.Migrasi(); err != nil {
		log.Fatal("Gagal menjalankan migrasi: ", err)
	}
	if err := config.DataAwal(); err != nil {
		log.Fatal("Gagal membuat data awal: ", err)
	}

	app := fiber.New(fiber.Config{
		AppName:            konf.NamaAplikasi,
		ErrorHandler:       penanganGalat,
		TrustProxy:         true,
		TrustProxyConfig:   fiber.TrustProxyConfig{Loopback: true},
		ProxyHeader:        fiber.HeaderXForwardedFor,
		EnableIPValidation: true,
	})

	app.Use(requestid.New())
	app.Use(corsHandler)
	app.Use(logger.New())
	app.Use(recoverer.New())
	app.Use(middleware.WaktuProses())

	router.Daftarkan(app)

	// Menangani semua permintaan yang tidak cocok dengan rute.
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Endpoint tidak tersedia",
		})
	})

	log.Fatal(app.Listen("127.0.0.1:" + konf.Port))
}

func buatCorsMiddleware() (fiber.Handler, error) {
	rawOrigins := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if rawOrigins == "" {
		rawOrigins = "http://localhost:5173"
	}

	origins := strings.Split(rawOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
		if origins[i] == "" {
			return nil, errors.New("CORS_ALLOWED_ORIGINS memuat origin kosong")
		}
	}

	return cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodOptions,
		},
		AllowHeaders: []string{
			fiber.HeaderOrigin,
			fiber.HeaderContentType,
			fiber.HeaderAccept,
			fiber.HeaderAuthorization,
			"X-Api-Key",
		},
	}), nil
}

func penanganGalat(c fiber.Ctx, err error) error {
	kode := fiber.StatusInternalServerError
	pesan := "Terjadi kesalahan pada server"

	var galatFiber *fiber.Error
	if errors.As(err, &galatFiber) {
		kode = galatFiber.Code
		pesan = galatFiber.Message
	}

	if kode >= fiber.StatusInternalServerError {
		log.Println("Error server:", err)
		pesan = "Terjadi kesalahan pada server"
	}

	return c.Status(kode).JSON(fiber.Map{
		"sukses": false,
		"pesan":  pesan,
	})
}
