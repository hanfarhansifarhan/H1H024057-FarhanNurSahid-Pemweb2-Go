package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestPeriksaTokenValidasiSignatureDanKlaim(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-with-more-than-thirty-two-bytes")
	t.Setenv("JWT_TTL_MENIT", "120")

	tests := []struct {
		name   string
		method jwt.SigningMethod
		claims Klaim
		secret string
		valid  bool
	}{
		{
			name:   "token valid",
			method: jwt.SigningMethodHS256,
			claims: klaimUji(),
			secret: "test-secret-with-more-than-thirty-two-bytes",
			valid:  true,
		},
		{
			name:   "signature mismatch",
			method: jwt.SigningMethodHS256,
			claims: klaimUji(),
			secret: "different-secret-with-more-than-thirty-two-bytes",
		},
		{
			name:   "wrong issuer",
			method: jwt.SigningMethodHS256,
			claims: func() Klaim {
				klaim := klaimUji()
				klaim.Issuer = "another-issuer"
				return klaim
			}(),
			secret: "test-secret-with-more-than-thirty-two-bytes",
		},
		{
			name:   "expired",
			method: jwt.SigningMethodHS256,
			claims: func() Klaim {
				klaim := klaimUji()
				klaim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Minute))
				return klaim
			}(),
			secret: "test-secret-with-more-than-thirty-two-bytes",
		},
		{
			name:   "missing expiration",
			method: jwt.SigningMethodHS256,
			claims: func() Klaim {
				klaim := klaimUji()
				klaim.ExpiresAt = nil
				return klaim
			}(),
			secret: "test-secret-with-more-than-thirty-two-bytes",
		},
		{
			name:   "missing issued at",
			method: jwt.SigningMethodHS256,
			claims: func() Klaim {
				klaim := klaimUji()
				klaim.IssuedAt = nil
				return klaim
			}(),
			secret: "test-secret-with-more-than-thirty-two-bytes",
		},
		{
			name:   "missing token id",
			method: jwt.SigningMethodHS256,
			claims: func() Klaim {
				klaim := klaimUji()
				klaim.JTI = ""
				return klaim
			}(),
			secret: "test-secret-with-more-than-thirty-two-bytes",
		},
		{
			name:   "wrong algorithm",
			method: jwt.SigningMethodHS384,
			claims: klaimUji(),
			secret: "test-secret-with-more-than-thirty-two-bytes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := jwt.NewWithClaims(tt.method, tt.claims).SignedString([]byte(tt.secret))
			if err != nil {
				t.Fatalf("gagal menandatangani token uji: %v", err)
			}

			_, err = PeriksaToken(raw)
			if (err == nil) != tt.valid {
				t.Fatalf("validasi token error=%v, diharapkan valid=%t", err, tt.valid)
			}
		})
	}
}

func TestTerbitkanTokenMenyertakanKlaimWajib(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-with-more-than-thirty-two-bytes")
	t.Setenv("JWT_TTL_MENIT", "120")

	raw, err := TerbitkanToken(42, "mahasiswa")
	if err != nil {
		t.Fatalf("gagal menerbitkan token: %v", err)
	}

	klaim, err := PeriksaToken(raw)
	if err != nil {
		t.Fatalf("token yang diterbitkan gagal diperiksa: %v", err)
	}
	if klaim.Subject != "42" || klaim.Peran != "mahasiswa" || klaim.Issuer != "pemweb2-fiber-api" {
		t.Fatalf("klaim inti tidak sesuai: %+v", klaim.RegisteredClaims)
	}
	if klaim.IssuedAt == nil || klaim.ExpiresAt == nil || klaim.JTI == "" {
		t.Fatal("token tidak memuat klaim iat, exp, dan jti")
	}
	if ttl := klaim.ExpiresAt.Time.Sub(klaim.IssuedAt.Time); ttl != 120*time.Minute {
		t.Fatalf("durasi token = %s, diharapkan 120 menit", ttl)
	}
}

func klaimUji() Klaim {
	now := time.Now()
	return Klaim{
		Peran: "mahasiswa",
		JTI:   uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			Issuer:    "pemweb2-fiber-api",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
}
