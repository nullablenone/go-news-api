package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nullablenone/go-news-api/internal/shared/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg config.RedisConfig) (*redis.Client, error) {
	if cfg.Host == "" || cfg.Port == "" {
		log.Println("Redis config not set, running without Redis cache")
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
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")

	return rdb, nil
}
