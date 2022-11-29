package main

import (
	"context"
	"fmt"
	"gorooster/helpers"
	"gorooster/logger"
	"gorooster/services"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// init loger use zap
	logger := logger.GetLogger()

	redisHost := helpers.EnvGetString("REDIS_SERVER_IP", "localhost:6379")
	dbNumber := helpers.EnvGetInt("REDIS_SELECT_DB", 3)

	logger.AddData("REDIS_SERVER_IP", redisHost)
	logger.AddData("REDIS_SELECT_DB", dbNumber)
	// connect to redis
	redisDB := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   3,
	})

	// This is telling redis to publish events since it's off by default.
	// https://redis.io/docs/manual/keyspace-notifications/
	_, err := redisDB.Do(context.Background(), "CONFIG", "SET", "notify-keyspace-events", "KEA").Result()
	if err != nil {
		logger.Log.Errorf("unable to set keyspace events %s", err.Error())
		os.Exit(1)
	}

	KeyEventChannel := fmt.Sprintf("__keyevent@%d__:expired", dbNumber)
	// this is telling redis to subscribe to events published in the keyevent channel, specifically for expired events
	pubsub := redisDB.PSubscribe(context.Background(), KeyEventChannel)

	logger.Log.Infow("start service ", logger.Data()...)

	eventMapper := services.NewEventMapper()
	// Infinite loop for listening event
	for {
		// This listens in the background for messages.
		message, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			logger.Log.Errorw(err.Error(), logger.Data()...)
			break
		}
		key := message.Payload
		dataKey, err := helpers.GetDataKey(key)
		if err != nil {
			logger.Log.Errorw(err.Error(), logger.Data()...)
			continue
		}
		logger.AddData("event_key", key)
		logger.AddData("data_key", dataKey)
		// Get real data event from redis
		ctx := context.Background()
		data := redisDB.Get(ctx, dataKey).Val()

		// Delete data from resis
		if err = redisDB.Del(ctx, dataKey).Err(); err != nil {
			logger.Log.Errorw(err.Error(), logger.Data()...)
		}

		go eventMapper.CreateEvent(data)
	}
}
