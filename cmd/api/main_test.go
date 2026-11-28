package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestPenanganGalatTidakMembocorkanDetailInternal(t *testing.T) {
	const privateDetail = "private database connection detail"
	app := fiber.New(fiber.Config{ErrorHandler: penanganGalat})
	app.Get("/internal-error", func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusInternalServerError, privateDetail)
	})

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/internal-error", nil))
	if err != nil {
		t.Fatalf("gagal menjalankan request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("status = %d, diharapkan %d", response.StatusCode, fiber.StatusInternalServerError)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("gagal membaca respons: %v", err)
	}
	if strings.Contains(string(body), privateDetail) {
		t.Fatal("respons membocorkan detail internal")
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("respons bukan JSON valid: %v", err)
	}
	if result["sukses"] != false || result["pesan"] != "Terjadi kesalahan pada server" {
		t.Fatalf("struktur respons error tidak sesuai: %v", result)
	}
}

func TestCorsMembatasiOriginSesuaiKonfigurasi(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")

	handler, err := buatCorsMiddleware()
	if err != nil {
		t.Fatalf("gagal membuat middleware CORS: %v", err)
	}

	app := fiber.New()
	app.Use(handler)
	app.Get("/status", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	tests := []struct {
		name          string
		origin        string
		allowedOrigin string
	}{
		{
			name:          "allowed origin",
			origin:        "http://localhost:5173",
			allowedOrigin: "http://localhost:5173",
		},
		{
			name:   "unlisted origin",
			origin: "http://127.0.0.1:5173",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/status", nil)
			request.Header.Set("Origin", tt.origin)
			response, err := app.Test(request)
			if err != nil {
				t.Fatalf("gagal menjalankan request: %v", err)
			}
			defer response.Body.Close()

			if got := response.Header.Get(fiber.HeaderAccessControlAllowOrigin); got != tt.allowedOrigin {
				t.Fatalf("Access-Control-Allow-Origin = %q, diharapkan %q", got, tt.allowedOrigin)
			}
		})
	}
}
