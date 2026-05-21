package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgre(env *Env) (*gorm.DB, error) {

	// dsn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", env.DBHost, env.DBUser, env.DBPass, env.DBName, env.DBPort, env.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("Error connecting to database: %v", err)
	}

	log.Println("database connected successfully!")

	return db, nil
}
