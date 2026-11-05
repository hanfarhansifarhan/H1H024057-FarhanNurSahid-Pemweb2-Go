package repository

import (
	"latihan-fiber-api/internal/config"
	"latihan-fiber-api/internal/model"
)

func AmbilProgramStudi() ([]model.ProgramStudiRingkasan, error) {
	var daftar []model.ProgramStudiRingkasan

	err := config.DB.Model(&model.ProgramStudi{}).
		Select(
			"program_studis.id, program_studis.kode, program_studis.nama, " +
				"program_studis.jenjang, program_studis.created_at, " +
				"program_studis.updated_at, COUNT(mahasiswas.id) AS jumlah_mahasiswa",
		).
		Joins(
			"LEFT JOIN mahasiswas ON mahasiswas.program_studi_id = program_studis.id " +
				"AND mahasiswas.deleted_at IS NULL",
		).
		Group(
			"program_studis.id, program_studis.kode, program_studis.nama, " +
				"program_studis.jenjang, program_studis.created_at, program_studis.updated_at",
		).
		Order("program_studis.nama ASC").
		Scan(&daftar).Error

	return daftar, err
}
