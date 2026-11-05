package config

import (
	"log"

	"latihan-fiber-api/internal/model"
)

func Migrasi() error {
	if err := DB.AutoMigrate(
		&model.ProgramStudi{},
		&model.Mahasiswa{},
		&model.Matakuliah{},
	); err != nil {
		return err
	}

	log.Println("Migrasi selesai")
	return nil
}

func DataAwal() error {
	var jumlah int64

	if err := DB.Model(&model.ProgramStudi{}).Count(&jumlah).Error; err != nil {
		return err
	}

	if jumlah > 0 {
		return nil
	}

	daftar := []model.ProgramStudi{
		{
			Kode:    "TK",
			Nama:    "Teknik Komputer",
			Jenjang: "S1",
		},
		{
			Kode:    "IF",
			Nama:    "Informatika",
			Jenjang: "S1",
		},
		{
			Kode:    "TE",
			Nama:    "Teknik Elektro",
			Jenjang: "S1",
		},
	}

	if err := DB.Create(&daftar).Error; err != nil {
		return err
	}

	log.Println("Data awal program studi dibuat")
	return nil
}
