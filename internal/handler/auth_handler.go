package handler

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"latihan-fiber-api/internal/auth"
	"latihan-fiber-api/internal/config"
	"latihan-fiber-api/internal/model"
	"latihan-fiber-api/internal/validasi"
)

func Registrasi(c fiber.Ctx) error {
	var permintaan PermintaanRegistrasi
	if err := c.Bind().Body(&permintaan); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format body tidak valid")
	}

	if galat := validasi.Periksa(permintaan); len(galat) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Data yang dikirim tidak valid",
			"galat":  galat,
		})
	}

	email := strings.ToLower(strings.TrimSpace(permintaan.Email))
	var pengguna model.Pengguna
	if err := config.DB.Where("email = ?", email).First(&pengguna).Error; err == nil {
		return fiber.NewError(fiber.StatusConflict, "Email sudah terdaftar")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Gagal memeriksa email saat registrasi: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memproses registrasi")
	}

	hashKataSandi, err := bcrypt.GenerateFromPassword([]byte(permintaan.KataSandi), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Gagal membuat hash password: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memproses registrasi")
	}

	penggunaBaru := model.Pengguna{
		Nama:      strings.TrimSpace(permintaan.Nama),
		Email:     email,
		KataSandi: string(hashKataSandi),
		Peran:     "mahasiswa",
	}

	if err := config.DB.Create(&penggunaBaru).Error; err != nil {
		log.Printf("Gagal menyimpan pengguna: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan data pengguna")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Registrasi berhasil",
		"data": fiber.Map{
			"id":    penggunaBaru.ID,
			"nama":  penggunaBaru.Nama,
			"email": penggunaBaru.Email,
			"peran": penggunaBaru.Peran,
		},
	})
}

func Login(c fiber.Ctx) error {
	var permintaan PermintaanLogin
	if err := c.Bind().Body(&permintaan); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format body tidak valid")
	}

	if galat := validasi.Periksa(permintaan); len(galat) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Data yang dikirim tidak valid",
			"galat":  galat,
		})
	}

	email := strings.ToLower(strings.TrimSpace(permintaan.Email))
	var pengguna model.Pengguna
	if err := config.DB.Where("email = ?", email).First(&pengguna).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"sukses": false,
				"pesan":  "Email atau kata sandi tidak sesuai",
			})
		}
		log.Printf("Gagal mengambil pengguna saat login: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memproses login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pengguna.KataSandi), []byte(permintaan.KataSandi)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Email atau kata sandi tidak sesuai",
		})
	}

	token, err := auth.TerbitkanToken(pengguna.ID, pengguna.Peran)
	if err != nil {
		log.Printf("Gagal menerbitkan token: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memproses login")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Login berhasil",
		"data": fiber.Map{
			"pengguna": fiber.Map{
				"id":    pengguna.ID,
				"nama":  pengguna.Nama,
				"email": pengguna.Email,
				"peran": pengguna.Peran,
			},
			"token": token,
			"tipe":  "Bearer",
		},
	})
}

func Profil(c fiber.Ctx) error {
	idPengguna, ok := c.Locals("id_pengguna").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Token tidak valid",
		})
	}

	var pengguna model.Pengguna
	if err := config.DB.First(&pengguna, idPengguna).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Pengguna tidak ditemukan")
		}
		log.Printf("Gagal mengambil profil pengguna %d: %v", idPengguna, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil profil")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"data": fiber.Map{
			"id":    pengguna.ID,
			"nama":  pengguna.Nama,
			"email": pengguna.Email,
			"peran": pengguna.Peran,
		},
	})
}

func SegarkanToken(c fiber.Ctx) error {
	header := c.Get("Authorization")
	bagian := strings.Fields(header)
	if len(bagian) != 2 || !strings.EqualFold(bagian[0], "Bearer") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Format Authorization tidak valid",
		})
	}

	klaim, err := auth.PeriksaToken(bagian[1])
	if err != nil || klaim == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Token tidak valid atau kedaluwarsa",
		})
	}

	sudahDicabut, err := auth.SudahDiBlacklist(klaim.JTI)
	if err != nil {
		log.Printf("Gagal memeriksa blacklist token saat refresh: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memeriksa autentikasi")
	}
	if sudahDicabut {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Token sudah dicabut",
		})
	}

	idPengguna, err := strconv.ParseUint(klaim.Subject, 10, 32)
	if err != nil || idPengguna == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Token tidak valid",
		})
	}

	tokenBaru, err := auth.TerbitkanToken(uint(idPengguna), klaim.Peran)
	if err != nil {
		log.Printf("Gagal menerbitkan token baru: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memperbarui token")
	}

	if err := auth.CatatBlacklist(klaim.JTI, klaim.ExpiresAt.Time); err != nil {
		log.Printf("Gagal menyimpan token lama ke blacklist: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memperbarui token")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Token berhasil diperbarui",
		"data": fiber.Map{
			"token": tokenBaru,
			"tipe":  "Bearer",
		},
	})
}

func Logout(c fiber.Ctx) error {
	header := c.Get("Authorization")
	bagian := strings.Fields(header)
	if len(bagian) != 2 || !strings.EqualFold(bagian[0], "Bearer") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Format Authorization tidak valid",
		})
	}

	klaim, err := auth.PeriksaToken(bagian[1])
	if err != nil || klaim == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Token tidak valid atau kedaluwarsa",
		})
	}

	if err := auth.CatatBlacklist(klaim.JTI, klaim.ExpiresAt.Time); err != nil {
		log.Printf("Gagal menambahkan token ke blacklist: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal logout")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Logout berhasil",
	})
}
