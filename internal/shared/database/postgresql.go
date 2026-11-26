package database

import (
	"fmt"
	"log"
	"time"

	"github.com/nullablenone/go-news-api/internal/shared/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgre(cfg config.DBConfig) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("Error connecting to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.SetConnMaxIdleTime(15 * time.Minute)
	}

	log.Println("database connected successfully!")

	return db, nil
}
