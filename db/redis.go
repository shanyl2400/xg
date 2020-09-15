package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"sync"
)

var (
	globalRedis *redis.Client
	_redisOnce sync.Once
)

func GetRedis() *redis.Client{
	_redisOnce.Do(func() {
		globalRedis = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("xg_redis_conn"),
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		err := globalRedis.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	})
	return globalRedis
}
