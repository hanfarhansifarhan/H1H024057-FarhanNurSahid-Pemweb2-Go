package validasi

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type PesanGalat struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validasi = validator.New()

func init() {
	if err := validasi.RegisterValidation("nim", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() != reflect.String {
			return false
		}

		nilai := fl.Field().String()
		if nilai == "" {
			return false
		}

		if !strings.HasPrefix(nilai, "H") {
			return false
		}

		sisa := nilai[1:]
		if sisa == "" {
			return false
		}

		for _, karakter := range sisa {
			if karakter < '0' || karakter > '9' {
				return false
			}
		}

		return true
	}); err != nil {
		panic("gagal mendaftarkan validasi NIM")
	}
}

func Periksa(data any) []PesanGalat {
	if data == nil {
		return []PesanGalat{{Field: "data", Message: "Data tidak dapat divalidasi"}}
	}

	if err := validasi.Struct(data); err != nil {
		var galatValidasi validator.ValidationErrors
		if !errors.As(err, &galatValidasi) {
			return []PesanGalat{{Field: "data", Message: "Data tidak dapat divalidasi"}}
		}

		hasil := make([]PesanGalat, 0, len(galatValidasi))
		for _, galat := range galatValidasi {
			hasil = append(hasil, PesanGalat{
				Field:   namaBidang(data, galat.StructField()),
				Message: susunPesan(galat),
			})
		}

		return hasil
	}

	return nil
}

func namaBidang(data any, namaKolom string) string {
	tipe := reflect.TypeOf(data)
	if tipe == nil {
		return namaKolom
	}

	if tipe.Kind() == reflect.Pointer {
		tipe = tipe.Elem()
	}

	for i := 0; i < tipe.NumField(); i++ {
		field := tipe.Field(i)
		if field.Name != namaKolom {
			continue
		}

		if tag := field.Tag.Get("json"); tag != "" {
			bagian := strings.Split(tag, ",")
			if bagian[0] != "" && bagian[0] != "-" {
				return bagian[0]
			}
		}
		return field.Name
	}

	return namaKolom
}

func susunPesan(galat validator.FieldError) string {
	switch galat.Tag() {
	case "required":
		return "Field ini wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "min":
		return fmt.Sprintf("Nilai minimal adalah %s", galat.Param())
	case "max":
		return fmt.Sprintf("Nilai maksimal adalah %s", galat.Param())
	case "gte":
		return fmt.Sprintf("Nilai tidak boleh kurang dari %s", galat.Param())
	case "lte":
		return fmt.Sprintf("Nilai tidak boleh lebih dari %s", galat.Param())
	case "nim":
		return "Format NIM tidak valid"
	default:
		return "Nilai tidak valid"
	}
}

func (p PesanGalat) String() string {
	return p.Field + ": " + p.Message
}
