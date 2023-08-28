package db

import (
	"github.com/go-redis/redis"
	"github.com/hiosi123/udongmat/config"
)

func CreateRedisConnection(env config.EnvVars) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.REDIS_ADDR,
		Password: env.REDIS_PASSWORD,
		DB:       env.REDIS_DB,
	})

	return rdb
}
