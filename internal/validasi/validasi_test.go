package validasi

import "testing"

type permintaanUjiNIM struct {
	Nim string `json:"nim" validate:"required,min=5,max=20,nim"`
}

func TestValidasiNIMKustom(t *testing.T) {
	tests := []struct {
		nim string
		ok  bool
	}{
		{nim: "H1024099", ok: true},
		{nim: "1024099"},
		{nim: "A1024099"},
		{nim: "HABC123"},
		{nim: " H1024099"},
		{nim: "H1024099 "},
	}

	for _, tt := range tests {
		t.Run(tt.nim, func(t *testing.T) {
			galat := Periksa(permintaanUjiNIM{Nim: tt.nim})
			if (len(galat) == 0) != tt.ok {
				t.Fatalf("NIM %q menghasilkan %d galat; valid=%t", tt.nim, len(galat), tt.ok)
			}
			if !tt.ok {
				for _, item := range galat {
					if item.Field == "nim" && item.Message == "Format NIM tidak valid" {
						return
					}
				}
				t.Fatalf("NIM %q tidak menghasilkan galat format NIM", tt.nim)
			}
		})
	}
}

func TestPeriksaTidakPanicUntukDataBukanStruct(t *testing.T) {
	if galat := Periksa(nil); len(galat) == 0 {
		t.Fatal("data nil harus menghasilkan galat validasi")
	}

	if galat := Periksa("bukan struct"); len(galat) == 0 {
		t.Fatal("data bukan struct harus menghasilkan galat validasi")
	}
}
