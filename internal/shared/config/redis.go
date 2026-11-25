package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg RedisConfig) (*redis.Client, error) {
	if cfg.Host == "" || cfg.Port == "" {
		log.Println("Konfigurasi Redis tidak diisi, aplikasi berjalan tanpa cache Redis")
		return nil, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("gagal menyambung ke Redis: %v", err)
	}

	log.Println("Koneksi Redis Berhasil Terhubung")

	return rdb, nil
}
