package service

import (
	"errors"
	"strings"
	"unicode/utf8"

	"latihan-fiber-api/internal/model"
	"latihan-fiber-api/internal/repository"
)

var ErrMatakuliahTidakValid = errors.New("data matakuliah tidak valid")

func DaftarMatakuliah(filter model.FilterMatakuliah) ([]model.Matakuliah, int64, error) {
	return repository.AmbilMatakuliah(filter)
}

func DetailMatakuliah(id uint) (model.Matakuliah, error) {
	return repository.AmbilMatakuliahPerID(id)
}

func BuatMatakuliah(matakuliah *model.Matakuliah) error {
	if err := validasiMatakuliah(matakuliah); err != nil {
		return err
	}
	return repository.SimpanMatakuliah(matakuliah)
}

func PerbaruiMatakuliah(matakuliah *model.Matakuliah) error {
	if err := validasiMatakuliah(matakuliah); err != nil {
		return err
	}
	return repository.PerbaruiMatakuliah(matakuliah)
}

func HapusMatakuliah(matakuliah *model.Matakuliah) error {
	return repository.HapusMatakuliah(matakuliah)
}

func validasiMatakuliah(matakuliah *model.Matakuliah) error {
	if strings.TrimSpace(matakuliah.Kode) == "" ||
		strings.TrimSpace(matakuliah.Nama) == "" ||
		utf8.RuneCountInString(matakuliah.Kode) > 20 ||
		utf8.RuneCountInString(matakuliah.Nama) > 100 ||
		matakuliah.SKS < 1 || matakuliah.SKS > 6 ||
		matakuliah.Semester < 1 || matakuliah.Semester > 14 {
		return ErrMatakuliahTidakValid
	}
	return nil
}
