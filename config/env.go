package config

import (
	"log"
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
	JWTSecret string

	RedisHost string
	RedisPort string
	RedisPass string
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
		JWTSecret: os.Getenv("JWT_SECRET"),

		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),
		RedisPass: os.Getenv("REDIS_PASSWORD"),
	}

	log.Println("Berhasil Inisialisasi Env file")

	return &env, nil
}
