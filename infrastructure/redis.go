package infrastructure

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vnnyx/rekadigital-tech-test/internal/config"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Host,
	})
	return client
}

func NewRedisContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
