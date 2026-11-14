package model

import "time"

type Pengguna struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `gorm:"size:100;not null" json:"nama"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	KataSandi string    `gorm:"size:255;not null" json:"-"`
	Peran     string    `gorm:"size:20;default:'mahasiswa';not null" json:"peran"`
	CreatedAt time.Time `json:"dibuat_pada"`
	UpdatedAt time.Time `json:"diubah_pada"`
}
