package service

import (
	"errors"
	"testing"

	"latihan-fiber-api/internal/model"
)

func TestValidasiMatakuliah(t *testing.T) {
	tests := []struct {
		nama string
		data model.Matakuliah
		ada  bool
	}{
		{
			nama: "valid boundary values",
			data: model.Matakuliah{Kode: "IF101", Nama: "Algoritma", SKS: 1, Semester: 14},
		},
		{
			nama: "missing code",
			data: model.Matakuliah{Nama: "Algoritma", SKS: 3, Semester: 1},
			ada:  true,
		},
		{
			nama: "blank name",
			data: model.Matakuliah{Kode: "IF101", Nama: " ", SKS: 3, Semester: 1},
			ada:  true,
		},
		{
			nama: "invalid credit count",
			data: model.Matakuliah{Kode: "IF101", Nama: "Algoritma", SKS: 7, Semester: 1},
			ada:  true,
		},
		{
			nama: "invalid semester",
			data: model.Matakuliah{Kode: "IF101", Nama: "Algoritma", SKS: 3, Semester: 0},
			ada:  true,
		},
		{
			nama: "code exceeds column length",
			data: model.Matakuliah{Kode: "ABCDEFGHIJKLMNOPQRSTU", Nama: "Algoritma", SKS: 3, Semester: 1},
			ada:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			err := validasiMatakuliah(&tt.data)
			if errors.Is(err, ErrMatakuliahTidakValid) != tt.ada {
				t.Fatalf("validasi memberi error %v; diharapkan error=%t", err, tt.ada)
			}
		})
	}
}
