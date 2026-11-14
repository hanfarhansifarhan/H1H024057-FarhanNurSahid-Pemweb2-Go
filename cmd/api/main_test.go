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
