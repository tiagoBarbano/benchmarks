package redis

import (
	"context"
	"log"
	"strconv"

	"my-fiber-app/pkg/config"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect() {
	db, _ := strconv.Atoi(config.GetEnv("REDIS_DB", "0"))
	Client = redis.NewClient(&redis.Options{
		Addr:     config.GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: config.GetEnv("REDIS_PASSWORD", ""),
		DB:       db,
	})

	if err := Client.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Redis ping error:", err)
		panic("Failed to connect to Redis")
	}
	// redisotel.InstrumentTracing(Client)

	log.Println("✅ Connected to Redis")
}
