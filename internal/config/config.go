package config

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Konfigurasi struct {
	NamaAplikasi string
	Port         string
	Lingkungan   string

	DbHost  string
	DbPort  string
	DbNama  string
	DbUser  string
	DbSandi string

	JwtSecret   string
	JwtTtlMenit int
}

func ambil(kunci, nilaiDefault string) string {
	nilai := os.Getenv(kunci)

	if nilai == "" {
		return nilaiDefault
	}

	return nilai
}

func buatSecretJwt() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		panic("gagal menghasilkan rahasia JWT")
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}

func Muat() Konfigurasi {
	godotenv.Load()

	ttlDefault := 120
	if nilaiTtl, err := strconv.Atoi(ambil("JWT_TTL_MENIT", "120")); err == nil && nilaiTtl > 0 && nilaiTtl <= 153722867 {
		ttlDefault = nilaiTtl
	}

	secret := strings.TrimSpace(ambil("JWT_SECRET", ""))
	if secret == "" {
		secret = buatSecretJwt()
		if err := os.Setenv("JWT_SECRET", secret); err != nil {
			panic("gagal menyimpan rahasia JWT pada lingkungan proses")
		}
	}
	if len([]byte(secret)) < 32 || strings.HasPrefix(strings.ToLower(secret), "replace-with-") {
		panic("JWT_SECRET harus memiliki panjang minimal 32 byte")
	}

	return Konfigurasi{
		NamaAplikasi: ambil("APP_NAME", "Fiber API"),
		Port:         ambil("APP_PORT", "3000"),
		Lingkungan:   ambil("APP_ENV", "local"),

		DbHost:  ambil("DB_HOST", "127.0.0.1"),
		DbPort:  ambil("DB_PORT", "3306"),
		DbNama:  ambil("DB_NAME", "pemweb2_fiber"),
		DbUser:  ambil("DB_USER", "root"),
		DbSandi: ambil("DB_PASS", ""),

		JwtSecret:   secret,
		JwtTtlMenit: ttlDefault,
	}
}
