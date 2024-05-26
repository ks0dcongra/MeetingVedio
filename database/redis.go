package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisConn *redis.Client

func RedisDBinit() {
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	_, err := RedisConn.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
}
