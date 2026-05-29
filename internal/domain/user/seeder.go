package user

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RunAdminSeeder bertugas mengecek dan membuat default admin jika belum ada
func RunAdminSeeder(db *gorm.DB) {
	var adminCount int64

	err := db.Model(&User{}).Where("role = ?", "admin").Count(&adminCount).Error
	if err != nil {
		log.Printf("[SEEDER ERROR] Gagal mengecek admin: %v\n", err)
		return
	}

	if adminCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("#admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("[SEEDER ERROR] Gagal melakukan hash password admin: %v\n", err)
		}

		defaultAdmin := User{
			Name:     "Super Admin",
			Email:    "admin@news.com",
			Password: string(hashedPassword),
			Role:     "admin",
		}

		if err := db.Create(&defaultAdmin).Error; err != nil {
			log.Printf("[SEEDER ERROR] Gagal membuat default admin: %v\n", err)
		} else {
			log.Println("====== [SEEDER] Berhasil membuat akun admin default: admin@news.com / #admin123 ======")
		}
	} else {
		log.Println("====== [SEEDER] Akun admin sudah tersedia di database ======")
	}
}
