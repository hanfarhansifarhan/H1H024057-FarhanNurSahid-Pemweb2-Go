package handler

import (
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"latihan-fiber-api/internal/model"
	"latihan-fiber-api/internal/service"
)

func DaftarMatakuliah(c fiber.Ctx) error {
	halaman, _ := strconv.Atoi(c.Query("halaman", "1"))
	perHalaman, _ := strconv.Atoi(c.Query("per_halaman", "10"))
	semester, _ := strconv.Atoi(c.Query("semester", "0"))

	if halaman < 1 {
		halaman = 1
	}
	if perHalaman < 1 || perHalaman > 100 {
		perHalaman = 10
	}

	daftar, total, err := service.DaftarMatakuliah(model.FilterMatakuliah{
		Cari:       c.Query("cari", ""),
		Semester:   semester,
		Halaman:    halaman,
		PerHalaman: perHalaman,
	})
	if err != nil {
		log.Printf("Gagal mengambil daftar matakuliah: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
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

func DetailMatakuliah(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "ID harus berupa angka positif")
	}

	matakuliah, err := service.DetailMatakuliah(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Matakuliah tidak ditemukan")
		}
		log.Printf("Gagal mengambil matakuliah %d: %v", id, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"data":   matakuliah,
	})
}

func BuatMatakuliah(c fiber.Ctx) error {
	var permintaan PermintaanMatakuliah
	if err := c.Bind().Body(&permintaan); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format body tidak valid")
	}
	if permintaan.Kode == nil || permintaan.Nama == nil ||
		permintaan.SKS == nil || permintaan.Semester == nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Semua field matakuliah wajib diisi")
	}

	matakuliah := model.Matakuliah{
		Kode:     *permintaan.Kode,
		Nama:     *permintaan.Nama,
		SKS:      *permintaan.SKS,
		Semester: *permintaan.Semester,
	}
	if err := service.BuatMatakuliah(&matakuliah); err != nil {
		log.Printf("Gagal membuat matakuliah: %v", err)
		if errors.Is(err, service.ErrMatakuliahTidakValid) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Gagal menyimpan data, periksa keunikan kode")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data matakuliah berhasil dibuat",
		"data":   matakuliah,
	})
}

func PerbaruiMatakuliah(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "ID harus berupa angka positif")
	}

	matakuliah, err := service.DetailMatakuliah(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Matakuliah tidak ditemukan")
		}
		log.Printf("Gagal mengambil matakuliah %d untuk diperbarui: %v", id, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	var permintaan PermintaanMatakuliah
	if err := c.Bind().Body(&permintaan); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format body tidak valid")
	}
	if permintaan.Kode != nil {
		matakuliah.Kode = *permintaan.Kode
	}
	if permintaan.Nama != nil {
		matakuliah.Nama = *permintaan.Nama
	}
	if permintaan.SKS != nil {
		matakuliah.SKS = *permintaan.SKS
	}
	if permintaan.Semester != nil {
		matakuliah.Semester = *permintaan.Semester
	}

	if err := service.PerbaruiMatakuliah(&matakuliah); err != nil {
		log.Printf("Gagal memperbarui matakuliah %d: %v", id, err)
		if errors.Is(err, service.ErrMatakuliahTidakValid) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Gagal memperbarui data, periksa keunikan kode")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data matakuliah berhasil diperbarui",
		"data":   matakuliah,
	})
}

func HapusMatakuliah(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "ID harus berupa angka positif")
	}

	matakuliah, err := service.DetailMatakuliah(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Matakuliah tidak ditemukan")
		}
		log.Printf("Gagal mengambil matakuliah %d untuk dihapus: %v", id, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data")
	}

	if err := service.HapusMatakuliah(&matakuliah); err != nil {
		log.Printf("Gagal menghapus matakuliah %d: %v", id, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menghapus data")
	}

	return c.JSON(fiber.Map{
		"sukses": true,
		"pesan":  "Data matakuliah berhasil dihapus",
	})
}
