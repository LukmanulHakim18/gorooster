package database

import (
	"github.com/LukmanulHakim18/gorooster/helpers"
	"sync"

	"github.com/go-redis/redis/v8"
)

var once sync.Once

var redisClient *RedisClient

type RedisClient struct {
	DB       *redis.Client
	DBNumber int
}

func GetRedisClient() *RedisClient {

	once.Do(func() {
		if redisClient == nil {

			redisDB := redis.NewClient(&redis.Options{
				Addr:     helpers.EnvGetString("REDIS_SERVER_IP", "localhost:6379"),
				Password: helpers.EnvGetString("REDIS_PASSWORD", ""),
				DB:       helpers.EnvGetInt("REDIS_SELECT_DB", 3),
			})
			redisClient = &RedisClient{
				DB:       redisDB,
				DBNumber: helpers.EnvGetInt("REDIS_SELECT_DB", 3),
			}
		}
	})
	return redisClient
}
