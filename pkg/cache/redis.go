package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr, password string, db int) *RedisCche {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		password: password,
		DB:       db,
	})
	return &RedisClient{
		client: client,
		ctx:    ctx,
	}
}
