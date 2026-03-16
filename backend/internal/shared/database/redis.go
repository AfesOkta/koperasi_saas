package database

import (
	"context"
	"fmt"
	"log"

	"github.com/koperasi-gresik/backend/config"
	"github.com/redis/go-redis/v9"
)

// NewRedis creates a new Redis client.
func NewRedis(cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("✅ Redis connected successfully")
	return client
}
