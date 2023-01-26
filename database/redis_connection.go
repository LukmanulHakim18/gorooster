package database

import (
	"context"
	"sync"

	"github.com/LukmanulHakim18/gorooster/v2/helpers"

	"github.com/go-redis/redis/v8"
)

var once sync.Once

var redisClient *RedisClient

type RedisClient struct {
	DB        map[int]*redis.Client
	DBIndex   *redis.Client
	DBPointer int
	UseDB     int
}

func GetRedisClient() *RedisClient {

	once.Do(func() {
		if redisClient == nil {
			redisClient = &RedisClient{
				DB:        map[int]*redis.Client{},
				DBPointer: 1,
				UseDB:     helpers.EnvGetInt("USE_DATABASE", 3),
			}

			for i := 1; i <= redisClient.UseDB; i++ {
				redisClient.DB[i] = redis.NewClient(&redis.Options{
					Addr:     helpers.EnvGetString("REDIS_SERVER_IP", "localhost:6379"),
					Password: helpers.EnvGetString("REDIS_PASSWORD", ""),
					DB:       i,
				})
			}
			redisClient.DBIndex = redis.NewClient(&redis.Options{
				Addr:     helpers.EnvGetString("REDIS_SERVER_IP", "localhost:6379"),
				Password: helpers.EnvGetString("REDIS_PASSWORD", ""),
				DB:       0,
			})
		}
	})
	return redisClient
}
func (rc *RedisClient) Next() {
	rc.DBPointer++
	if rc.DBPointer > rc.UseDB {
		rc.DBPointer = 1
	}
}

// isConnect
func (rc *RedisClient) IsConnected() bool {
	ping := rc.DBIndex.Ping(context.Background()).Val()
	return ping == "PONG"
}
