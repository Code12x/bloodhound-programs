package main

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx context.Context
)

func connectRedis() {
	ctx = context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}
