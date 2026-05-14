package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DBHost    string
	DBUser    string
	DBPass    string
	DBName    string
	DBPort    string
	DBSSLMode string
}

func NewEnv() (*Env, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	env := Env{
		DBHost:    os.Getenv("DB_HOST"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASS"),
		DBName:    os.Getenv("DB_NAME"),
		DBPort:    os.Getenv("DB_PORT"),
		DBSSLMode: os.Getenv("DB_SSLMODE"),
	}

	return &env, nil
}
