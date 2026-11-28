package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestStatusReturnsOK(t *testing.T) {
	app := fiber.New()
	app.Get("/api/status", Status)

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/status", nil))
	if err != nil {
		t.Fatalf("gagal menjalankan request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, diharapkan %d", response.StatusCode, http.StatusOK)
	}
}

func TestDetailMahasiswaRejectsInvalidID(t *testing.T) {
	app := fiber.New()
	app.Get("/api/mahasiswa/:id", DetailMahasiswa)

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/mahasiswa/not-a-number", nil))
	if err != nil {
		t.Fatalf("gagal menjalankan request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, diharapkan %d", response.StatusCode, http.StatusBadRequest)
	}
}
