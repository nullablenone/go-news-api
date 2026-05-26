package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(env *Env) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort),
		Password: env.RedisPass,
		DB:       0, // Menggunakan database default Redis
	})

	// Membuat context dengan timeout 5 detik untuk tes koneksi (Ping)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Melakukan Ping ke Redis untuk memastikan koneksi benar-benar tersambung
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("gagal menyambung ke Redis: %v", err)
	}

	log.Println("Koneksi Redis Berhasil Terhubung")

	return rdb, nil
}
