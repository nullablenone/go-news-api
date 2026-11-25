package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
	JWT   JWTConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DBConfig struct {
	Host         string
	User         string
	Password     string
	Name         string
	Port         string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type JWTConfig struct {
	Secret string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	appPort := getEnv("APP_PORT", "")
	if appPort == "" {
		appPort = getEnv("PORT", "8888")
	}
	appPort = strings.TrimPrefix(appPort, ":")

	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "go-news-api"),
			Env:  getEnv("APP_ENV", "development"),
			Port: appPort,
		},
		DB: DBConfig{
			Host:         os.Getenv("DB_HOST"),
			User:         os.Getenv("DB_USER"),
			Password:     os.Getenv("DB_PASS"),
			Name:         os.Getenv("DB_NAME"),
			Port:         os.Getenv("DB_PORT"),
			SSLMode:      os.Getenv("DB_SSLMODE"),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 50),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
	}

	log.Println("config initialized successfully")

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(val)
	if err != nil || parsed <= 0 {
		log.Printf("warning: invalid %s=%q, using default %d", key, val, fallback)
		return fallback
	}
	return parsed
}
