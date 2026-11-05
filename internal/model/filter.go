package model

type FilterMahasiswa struct {
	Cari       string
	Angkatan   int
	Halaman    int
	PerHalaman int
}

type FilterMatakuliah struct {
	Cari       string
	Semester   int
	Halaman    int
	PerHalaman int
}
