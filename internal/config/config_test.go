package config

import "testing"

func TestMuatMenggunakanTTLDefaultUntukNilaiTidakValid(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-with-more-than-thirty-two-bytes")
	t.Setenv("JWT_TTL_MENIT", "bukan-angka")

	if konf := Muat(); konf.JwtTtlMenit != 120 {
		t.Fatalf("TTL = %d, diharapkan 120 menit", konf.JwtTtlMenit)
	}
}

func TestMuatMenolakSecretJWTYangLemah(t *testing.T) {
	t.Setenv("JWT_SECRET", "weak-secret")

	defer func() {
		if recover() == nil {
			t.Fatal("secret JWT yang lebih pendek dari 32 byte harus ditolak")
		}
	}()

	Muat()
}
