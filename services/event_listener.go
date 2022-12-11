package services

import (
	"context"
	"fmt"
	"os"

	"git.bluebird.id/mybb/gorooster/v2/database"
	"git.bluebird.id/mybb/gorooster/v2/helpers"
	"git.bluebird.id/mybb/gorooster/v2/logger"
	"github.com/go-redis/redis/v8"
)

func StartEventListeners(client *database.RedisClient) {
	for dbNumber, client := range client.DB {
		go StartEventListener(dbNumber, client)
	}
}

func StartEventListener(dbNumber int, client *redis.Client) {
	// init loger use zap
	logger := logger.GetLogger()
	// This is telling redis to publish events since it's off by default.
	// https://redis.io/docs/manual/keyspace-notifications/
	_, err := client.Do(context.Background(), "CONFIG", "SET", "notify-keyspace-events", "KEA").Result()
	if err != nil {
		logger.Log.Errorf("unable to set keyspace events %s", err.Error())
		os.Exit(1)
	}

	KeyEventChannel := fmt.Sprintf("__keyevent@%d__:expired", dbNumber)
	logger.AddData("key_event_channel", KeyEventChannel)

	// this is telling redis to subscribe to events published in the keyevent channel, specifically for expired events
	pubsub := client.PSubscribe(context.Background(), KeyEventChannel)

	logger.Log.Infow("start service ", logger.Data()...)
	logger.ClearData()

	eventMapper := NewEventMapper()
	ctx := context.Background()
	for {
		logger.ClearData()
		// Infinite loop for listening event
		// This listens in the background for messages.
		message, err := pubsub.ReceiveMessage(ctx)
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

		go logger.Log.Infow("create_event", logger.Data()...)
		go eventMapper.CreateEvent(ctx, client, dataKey)
	}
}
