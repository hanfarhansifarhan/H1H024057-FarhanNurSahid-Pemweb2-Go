package auth

import (
	"errors"
	"time"

	"gorm.io/gorm/clause"

	"latihan-fiber-api/internal/config"
	"latihan-fiber-api/internal/model"
)

func SudahDiBlacklist(jti string) (bool, error) {
	if jti == "" {
		return false, errors.New("identitas token kosong")
	}
	if config.DB == nil {
		return false, errors.New("koneksi basis data belum tersedia")
	}

	var total int64
	if err := config.DB.Model(&model.TokenBlacklist{}).Where("jti = ?", jti).Count(&total).Error; err != nil {
		return false, err
	}

	return total > 0, nil
}

func CatatBlacklist(jti string, expiresAt time.Time) error {
	if jti == "" {
		return errors.New("identitas token kosong")
	}
	if config.DB == nil {
		return errors.New("koneksi basis data belum tersedia")
	}

	return config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "jti"}},
		DoNothing: true,
	}).Create(&model.TokenBlacklist{JTI: jti, ExpiresAt: expiresAt}).Error
}
