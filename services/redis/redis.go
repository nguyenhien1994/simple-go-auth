package redis

import (
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var RedisClientInstance *redis.Client
var RedisClientInstanceOnce sync.Once

func GetRedisClient() *redis.Client {
	RedisClientInstanceOnce.Do(func() {
		RedisClientInstance = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	})

	return RedisClientInstance
}
