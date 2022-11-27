package main

import (
	"context"
	"event-scheduler/helpers"
	"event-scheduler/logger"
	"event-scheduler/services"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis/v8"
)

func main() {

	// init loger use zap
	logger := logger.GetLogger()

	// Get setup redis
	redisSetup := helpers.GetRedisSetup()

	logger.AddData("REDIS_SERVER_IP", redisSetup.Host)
	logger.AddData("REDIS_SELECT_DB", redisSetup.SelectDB)
	// connect to redis
	redisDB := redis.NewClient(&redis.Options{
		Addr: redisSetup.Host,
		DB:   redisSetup.SelectDB,
	})

	// this is telling redis to publish events since it's off by default.
	// https://redis.io/docs/manual/keyspace-notifications/
	_, err := redisDB.Do(context.Background(), "CONFIG", "SET", "notify-keyspace-events", "KEA").Result()
	if err != nil {
		logger.Log.Errorf("unable to set keyspace events %s", err.Error())
		os.Exit(1)
	}

	// this is telling redis to subscribe to events published in the keyevent channel, specifically for expired events
	pubsub := redisDB.PSubscribe(context.Background(), redisSetup.GenerateKeyEventChannel())

	logger.Log.Infow("start service ", logger.Data()...)

	eventMapper := services.NewEventMapper()
	// infinite loop
	for {
		// this listens in the background for messages.
		message, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			logger.Log.Errorw(err.Error(), logger.Data()...)
			break
		}
		spew.Dump(message.Payload)
		key := message.Payload
		dataKey, err := helpers.GetDataKey(key)
		if err != nil {
			logger.Log.Errorw(err.Error(), logger.Data()...)
			continue
		}
		data := redisDB.Get(context.Background(), dataKey).Val()
		go eventMapper.CreateEvent(data)
	}
}
