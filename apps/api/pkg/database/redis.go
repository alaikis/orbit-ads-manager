package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"orbit/apps/api/config"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func ConnectRedis() {
	cfg := config.AppCfg
	url := cfg.Redis.URL
	if url == "" {
		url = "redis://localhost:6379/0"
	}

	addr := "localhost:6379"
	if strings.HasPrefix(url, "redis://") {
		addr = strings.TrimPrefix(url, "redis://")
		addr = strings.Split(addr, "/")[0]
	}

	Redis = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	fmt.Println("redis connected")
}
