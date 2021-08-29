package redis

import (
	"github.com/go-redis/redis/v8"
	"sync"
)

type RedisClientInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

var RedisClientInstance *redis.Client
var RedisClientInstanceOnce sync.Once

func GetRedisClient(info RedisClientInfo) *redis.Client {
	RedisClientInstanceOnce.Do(func() {
		RedisClientInstance = redis.NewClient(&redis.Options{
			Addr:     info.Host + ":" + info.Port,
			Password: info.Password,
			DB:       0,
		})
	})

	return RedisClientInstance
}
