package service

import (
	"latihan-fiber-api/internal/model"
	"latihan-fiber-api/internal/repository"
)

func DaftarMahasiswa(filter model.FilterMahasiswa) ([]model.Mahasiswa, int64, error) {
	return repository.Ambil(filter)
}

func DetailMahasiswa(id uint) (model.Mahasiswa, error) {
	return repository.AmbilPerId(id)
}

func BuatMahasiswa(mahasiswa *model.Mahasiswa) error {
	return repository.Simpan(mahasiswa)
}

func PerbaruiMahasiswa(mahasiswa *model.Mahasiswa) error {
	return repository.Perbarui(mahasiswa)
}

func HapusMahasiswa(mahasiswa *model.Mahasiswa) error {
	return repository.Hapus(mahasiswa)
}

func PulihkanMahasiswa(id uint) error {
	return repository.Pulihkan(id)
}
