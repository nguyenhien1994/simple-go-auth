package redis

import (
	"github.com/go-redis/redis/v8"
)

type RedisClientInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

func NewRedisClient(info RedisClientInfo) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     info.Host + ":" + info.Port,
		Password: info.Password,
		DB:       0,
	})

	return redisClient
}
