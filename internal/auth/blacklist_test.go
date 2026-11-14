package auth

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/glebarez/sqlite"

	"latihan-fiber-api/internal/config"
	"latihan-fiber-api/internal/model"
)

func TestBlacklistOperationsFailWhenDatabaseIsUnavailable(t *testing.T) {
	originalDB := config.DB
	config.DB = nil
	t.Cleanup(func() {
		config.DB = originalDB
	})

	revoked, err := SudahDiBlacklist("audit-jti")
	if err == nil || revoked {
		t.Fatalf("pemeriksaan blacklist tanpa DB harus gagal tertutup, revoked=%t err=%v", revoked, err)
	}

	if err := CatatBlacklist("audit-jti", time.Now()); err == nil {
		t.Fatal("pencatatan blacklist tanpa DB harus menghasilkan error")
	}
}

func TestBlacklistDuplicateJTIIsIgnored(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("gagal membuka database sqlite untuk uji blacklist: %v", err)
	}

	originalDB := config.DB
	config.DB = db
	t.Cleanup(func() {
		config.DB = originalDB
	})

	if err := db.AutoMigrate(&model.TokenBlacklist{}); err != nil {
		t.Fatalf("gagal migrasi tabel blacklist: %v", err)
	}

	if err := CatatBlacklist("duplicate-jti", time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("pencatatan blacklist pertama gagal: %v", err)
	}
	if err := CatatBlacklist("duplicate-jti", time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("pencatatan blacklist duplikat harus diabaikan: %v", err)
	}

	revoked, err := SudahDiBlacklist("duplicate-jti")
	if err != nil {
		t.Fatalf("pemeriksaan blacklist gagal: %v", err)
	}
	if !revoked {
		t.Fatal("JTI duplikat harus terdeteksi sebagai dicabut")
	}

	var count int64
	if err := db.Model(&model.TokenBlacklist{}).Where("jti = ?", "duplicate-jti").Count(&count).Error; err != nil {
		t.Fatalf("gagal menghitung catatan blacklist: %v", err)
	}
	if count != 1 {
		t.Fatalf("jumlah catatan blacklist = %d, diharapkan 1", count)
	}
}
