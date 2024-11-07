package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct{
	client *redis.Client
	ctx context.Context
}
func NewRedisCache(addr, password string, db int) *RedisCache{
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db,
	})
	return &RedisCache{
		client: client,
		ctx: ctx,
	}
}
func (r *RedisCache) Set(key string, value[]byte, expiration time.Duration)error{
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisCache) Get(key string)([]byte, error){
	return r.client.Get(r.ctx, key).Bytes()
}
func (r *RedisCache) Close()error{
	return r.client.Close()

}
