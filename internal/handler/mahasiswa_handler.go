package handler

import (
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"latihan-fiber-api/internal/model"
	"latihan-fiber-api/internal/service"
	"latihan-fiber-api/internal/validasi"
)

func DaftarMahasiswa(c fiber.Ctx) error {
	halaman, _ := strconv.Atoi(c.Query("halaman", "1"))
	perHalaman, _ := strconv.Atoi(c.Query("per_halaman", "10"))
	angkatan, _ := strconv.Atoi(c.Query("angkatan", "0"))

	if halaman < 1 {
		halaman = 1
	}

	if perHalaman < 1 || perHalaman > 100 {
		perHalaman = 10
	}

	filter := model.FilterMahasiswa{
		Cari:       c.Query("cari", ""),
		Angkatan:   angkatan,
		Halaman:    halaman,
		PerHalaman: perHalaman,
	}

	daftar, total, err := service.DaftarMahasiswa(filter)

	if err != nil {
		log.Printf("Gagal mengambil daftar mahasiswa: %v", err)
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"Gagal mengambil data",
		)
	}

	totalHalaman := int(total) / perHalaman

	if int(total)%perHalaman != 0 {
		totalHalaman++
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"data":   daftar,
		"meta": fiber.Map{
			"halaman":       halaman,
			"per_halaman":   perHalaman,
			"total":         total,
			"total_halaman": totalHalaman,
		},
	})
}

func DetailMahasiswa(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil || id < 1 {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"ID harus berupa angka positif",
		)
	}

	mahasiswa, err := service.DetailMahasiswa(uint(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(
				fiber.StatusNotFound,
				"Mahasiswa tidak ditemukan",
			)
		}

		log.Printf("Gagal mengambil mahasiswa %d: %v", id, err)
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"Gagal mengambil data",
		)
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"data":   mahasiswa,
	})
}

func BuatMahasiswa(c fiber.Ctx) error {
	var permintaan PermintaanMahasiswa

	if err := c.Bind().Body(&permintaan); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"Format body tidak valid",
		)
	}

	if galat := validasi.Periksa(permintaan); len(galat) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Data yang dikirim tidak valid",
			"galat":  galat,
		})
	}

	aktif := true

	if permintaan.Aktif != nil {
		aktif = *permintaan.Aktif
	}

	mahasiswa := model.Mahasiswa{
		ProgramStudiID: permintaan.ProgramStudiID,
		Nim:            permintaan.Nim,
		Nama:           permintaan.Nama,
		Email:          permintaan.Email,
		Angkatan:       permintaan.Angkatan,
		Ipk:            permintaan.Ipk,
		Aktif:          aktif,
	}

	err := service.BuatMahasiswa(&mahasiswa)

	if err != nil {
		log.Printf("Gagal menyimpan mahasiswa: %v", err)
		return fiber.NewError(
			fiber.StatusUnprocessableEntity,
			"Gagal menyimpan data, periksa keunikan NIM dan email",
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data mahasiswa berhasil dibuat",
		"data":   mahasiswa,
	})
}

func PerbaruiMahasiswa(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil || id < 1 {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"ID harus berupa angka positif",
		)
	}

	mahasiswa, err := service.DetailMahasiswa(uint(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Mahasiswa tidak ditemukan")
		}
		log.Printf("Gagal mengambil mahasiswa %d untuk diperbarui: %v", id, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	var permintaan PermintaanPembaruanMahasiswa

	if err := c.Bind().Body(&permintaan); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"Format body tidak valid",
		)
	}

	if galat := validasi.Periksa(permintaan); len(galat) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"sukses": false,
			"pesan":  "Data yang dikirim tidak valid",
			"galat":  galat,
		})
	}

	if permintaan.Nim != nil {
		mahasiswa.Nim = *permintaan.Nim
	}
	if permintaan.Nama != nil {
		mahasiswa.Nama = *permintaan.Nama
	}
	if permintaan.Email != nil {
		mahasiswa.Email = *permintaan.Email
	}
	if permintaan.Angkatan != nil {
		mahasiswa.Angkatan = *permintaan.Angkatan
	}
	if permintaan.Ipk != nil {
		mahasiswa.Ipk = *permintaan.Ipk
	}
	if permintaan.ProgramStudiID != nil {
		mahasiswa.ProgramStudiID = *permintaan.ProgramStudiID
	}

	if permintaan.Aktif != nil {
		mahasiswa.Aktif = *permintaan.Aktif
	}

	err = service.PerbaruiMahasiswa(&mahasiswa)

	if err != nil {
		log.Printf("Gagal memperbarui mahasiswa %d: %v", id, err)
		return fiber.NewError(
			fiber.StatusUnprocessableEntity,
			"Gagal memperbarui data",
		)
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data mahasiswa berhasil diperbarui",
		"data":   mahasiswa,
	})
}

func HapusMahasiswa(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil || id < 1 {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"ID harus berupa angka positif",
		)
	}

	mahasiswa, err := service.DetailMahasiswa(uint(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Mahasiswa tidak ditemukan")
		}
		log.Printf("Gagal mengambil mahasiswa %d untuk dihapus: %v", id, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	if err := service.HapusMahasiswa(&mahasiswa); err != nil {
		log.Printf("Gagal menghapus mahasiswa %d: %v", id, err)
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"Gagal menghapus data",
		)
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data mahasiswa berhasil dihapus",
	})
}

func PulihkanMahasiswa(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "ID harus berupa angka positif")
	}

	if err := service.PulihkanMahasiswa(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Mahasiswa terhapus tidak ditemukan")
		}
		log.Printf("Gagal memulihkan mahasiswa %d: %v", id, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memulihkan data")
	}

	mahasiswa, err := service.DetailMahasiswa(uint(id))
	if err != nil {
		log.Printf("Gagal mengambil mahasiswa %d setelah dipulihkan: %v", id, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Data dipulihkan tetapi gagal diambil")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data mahasiswa berhasil dipulihkan",
		"data":   mahasiswa,
	})
}
