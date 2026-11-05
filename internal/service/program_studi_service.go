package service

import (
	"latihan-fiber-api/internal/model"
	"latihan-fiber-api/internal/repository"
)

func DaftarProgramStudi() ([]model.ProgramStudiRingkasan, error) {
	return repository.AmbilProgramStudi()
}
