package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestAPIRequiresAPIKey(t *testing.T) {
	t.Setenv("API_KEY", "test-api-key")

	app := fiber.New()
	Daftarkan(app)

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/mahasiswa/"},
		{http.MethodPost, "/api/mahasiswa/"},
		{http.MethodGet, "/api/mahasiswa/1"},
		{http.MethodPut, "/api/mahasiswa/1"},
		{http.MethodDelete, "/api/mahasiswa/1"},
		{http.MethodPost, "/api/mahasiswa/1/restore"},
		{http.MethodGet, "/api/matakuliah/"},
		{http.MethodPost, "/api/matakuliah/"},
		{http.MethodGet, "/api/matakuliah/1"},
		{http.MethodPut, "/api/matakuliah/1"},
		{http.MethodDelete, "/api/matakuliah/1"},
		{http.MethodGet, "/api/program-studi"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			response, err := app.Test(request)
			if err != nil {
				t.Fatalf("gagal menjalankan request: %v", err)
			}
			defer response.Body.Close()

			if response.StatusCode != fiber.StatusUnauthorized {
				t.Fatalf(
					"status = %d, diharapkan %d; endpoint belum tentu terdaftar atau tidak terlindungi",
					response.StatusCode,
					fiber.StatusUnauthorized,
				)
			}
		})
	}
}
