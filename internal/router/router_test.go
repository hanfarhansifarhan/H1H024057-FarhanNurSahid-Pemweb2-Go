package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestProtectedRoutesRequireJWT(t *testing.T) {
	app := fiber.New()
	Daftarkan(app)

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/auth/profil"},
		{http.MethodPost, "/api/auth/segarkan"},
		{http.MethodPost, "/api/auth/logout"},
		{http.MethodGet, "/api/mahasiswa/"},
		{http.MethodPost, "/api/mahasiswa/"},
		{http.MethodGet, "/api/mahasiswa/1"},
		{http.MethodPut, "/api/mahasiswa/1"},
		{http.MethodDelete, "/api/mahasiswa/1"},
		{http.MethodPost, "/api/mahasiswa/1/restore"},
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
				t.Fatalf("status = %d, diharapkan %d", response.StatusCode, fiber.StatusUnauthorized)
			}
		})
	}
}

func TestAuthPublicRoutesDoNotRequireAPIKey(t *testing.T) {
	app := fiber.New()
	Daftarkan(app)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/registrasi", strings.NewReader(`{}`))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("gagal menjalankan request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusUnprocessableEntity {
		t.Fatalf("status = %d, diharapkan %d", response.StatusCode, fiber.StatusUnprocessableEntity)
	}
}

func TestStatusRouteIsPublic(t *testing.T) {
	app := fiber.New()
	Daftarkan(app)

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/status", nil))
	if err != nil {
		t.Fatalf("gagal menjalankan request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, diharapkan %d", response.StatusCode, fiber.StatusOK)
	}
}

func TestLoginIsRateLimited(t *testing.T) {
	app := fiber.New()
	Daftarkan(app)

	for attempt := 1; attempt <= 6; attempt++ {
		request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{}`))
		request.Header.Set("Content-Type", "application/json")
		response, err := app.Test(request)
		if err != nil {
			t.Fatalf("gagal menjalankan percobaan ke-%d: %v", attempt, err)
		}
		response.Body.Close()

		expected := fiber.StatusUnprocessableEntity
		if attempt == 6 {
			expected = fiber.StatusTooManyRequests
		}
		if response.StatusCode != expected {
			t.Fatalf("percobaan ke-%d: status = %d, diharapkan %d", attempt, response.StatusCode, expected)
		}
	}
}

func TestMatakuliahWritesRequireBearerTokenWithAPIKey(t *testing.T) {
	t.Setenv("API_KEY", "test-api-key")

	app := fiber.New()
	Daftarkan(app)

	request := httptest.NewRequest(http.MethodPost, "/api/matakuliah/", strings.NewReader(`{}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", "test-api-key")
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("gagal menjalankan request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("status = %d, diharapkan %d", response.StatusCode, fiber.StatusUnauthorized)
	}
}

func TestLegacyRoutesRequireAPIKeyEvenWithBearerHeader(t *testing.T) {
	t.Setenv("API_KEY", "test-api-key")

	app := fiber.New()
	Daftarkan(app)

	tests := []struct {
		name          string
		authorization string
	}{
		{name: "missing API key"},
		{name: "fabricated bearer", authorization: "Bearer not-a-real-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/program-studi", nil)
			if tt.authorization != "" {
				request.Header.Set("Authorization", tt.authorization)
			}
			response, err := app.Test(request)
			if err != nil {
				t.Fatalf("gagal menjalankan request: %v", err)
			}
			defer response.Body.Close()

			if response.StatusCode != fiber.StatusUnauthorized {
				t.Fatalf("status = %d, diharapkan %d", response.StatusCode, fiber.StatusUnauthorized)
			}
		})
	}
}
