package model

import "time"

type TokenBlacklist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	JTI       string    `gorm:"size:255;uniqueIndex;not null" json:"jti"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"dibuat_pada"`
	UpdatedAt time.Time `json:"diubah_pada"`
}
