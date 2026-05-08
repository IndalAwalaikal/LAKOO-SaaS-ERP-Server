package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"lakoo/backend/pkg/config"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	var err error
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = rdb.Ping(ctx).Result()
		cancel()

		if err == nil {
			break
		}
		log.Printf("Failed to connect to Redis (attempt %d/10): %v. Retrying in 5s...", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to Redis after 10 attempts: %v", err)
	}

	log.Println("Connected to Redis successfully")
	return rdb
}
