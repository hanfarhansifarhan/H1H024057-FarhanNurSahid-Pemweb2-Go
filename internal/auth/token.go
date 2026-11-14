package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"latihan-fiber-api/internal/config"
)

type Klaim struct {
	Peran string `json:"peran"`
	JTI   string `json:"jti,omitempty"`
	jwt.RegisteredClaims
}

func TerbitkanToken(id uint, peran string) (string, error) {
	if id == 0 || !peranValid(peran) {
		return "", errors.New("identitas atau peran pengguna tidak valid")
	}

	konf := config.Muat()
	waktuSekarang := time.Now()
	klaim := &Klaim{
		Peran: peran,
		JTI:   uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(id), 10),
			Issuer:    "pemweb2-fiber-api",
			IssuedAt:  jwt.NewNumericDate(waktuSekarang),
			ExpiresAt: jwt.NewNumericDate(waktuSekarang.Add(time.Duration(konf.JwtTtlMenit) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, klaim)
	return token.SignedString([]byte(konf.JwtSecret))
}

func PeriksaToken(tokenString string) (*Klaim, error) {
	if tokenString == "" {
		return nil, errors.New("token kosong")
	}

	konf := config.Muat()
	klaim := &Klaim{}
	token, err := jwt.ParseWithClaims(tokenString, klaim, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("algoritma token tidak valid")
		}
		return []byte(konf.JwtSecret), nil
	}, jwt.WithIssuer("pemweb2-fiber-api"), jwt.WithIssuedAt())
	if err != nil {
		return nil, err
	}
	if token == nil || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	id, err := strconv.ParseUint(klaim.Subject, 10, 32)
	if err != nil || id == 0 {
		return nil, errors.New("sub tidak valid")
	}

	if !peranValid(klaim.Peran) {
		return nil, errors.New("peran tidak valid")
	}

	if klaim.IssuedAt == nil || klaim.ExpiresAt == nil || klaim.JTI == "" {
		return nil, errors.New("klaim wajib tidak tersedia")
	}

	return klaim, nil
}

func peranValid(peran string) bool {
	switch peran {
	case "mahasiswa", "dosen", "admin":
		return true
	default:
		return false
	}
}
