package user

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunAdminSeeder(db *gorm.DB) {
	var adminCount int64

	err := db.Model(&User{}).Where("role = ?", "admin").Count(&adminCount).Error
	if err != nil {
		log.Printf("[SEEDER ERROR] failed to check admin: %v\n", err)
		return
	}

	if adminCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("#admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("[SEEDER ERROR] failed to hash admin password: %v\n", err)
		}

		defaultAdmin := User{
			Name:     "Super Admin",
			Email:    "admin@news.com",
			Password: string(hashedPassword),
			Role:     "admin",
		}

		if err := db.Create(&defaultAdmin).Error; err != nil {
			log.Printf("[SEEDER ERROR] failed to create default admin: %v\n", err)
		} else {
			log.Println("====== [SEEDER] Default admin created: admin@news.com / #admin123 ======")
		}
	} else {
		log.Println("====== [SEEDER] Admin account already exists in database ======")
	}
}
