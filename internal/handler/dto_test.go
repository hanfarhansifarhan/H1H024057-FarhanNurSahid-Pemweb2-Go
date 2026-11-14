package handler

import (
	"testing"

	"latihan-fiber-api/internal/validasi"
)

func TestValidasiPembaruanMahasiswa(t *testing.T) {
	namaKosong := ""
	ipkNol := 0.0
	ipkTidakValid := 4.1
	nimInvalid := "A1024099"

	tests := []struct {
		name string
		data PermintaanPembaruanMahasiswa
		ok   bool
	}{
		{name: "ipk zero is valid", data: PermintaanPembaruanMahasiswa{Ipk: &ipkNol}, ok: true},
		{name: "empty name is invalid", data: PermintaanPembaruanMahasiswa{Nama: &namaKosong}},
		{name: "ipk above limit is invalid", data: PermintaanPembaruanMahasiswa{Ipk: &ipkTidakValid}},
		{name: "invalid NIM is rejected", data: PermintaanPembaruanMahasiswa{Nim: &nimInvalid}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			galat := validasi.Periksa(tt.data)
			if (len(galat) == 0) != tt.ok {
				t.Fatalf("jumlah galat = %d, diharapkan valid=%t", len(galat), tt.ok)
			}
		})
	}
}
