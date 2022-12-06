package services

import (
	"context"
	"fmt"
	"os"

	"github.com/LukmanulHakim18/gorooster/database"
	"github.com/LukmanulHakim18/gorooster/helpers"
	"github.com/LukmanulHakim18/gorooster/logger"
)

func StartEventListener(client database.RedisClient) {
	// init loger use zap
	logger := logger.GetLogger()
	// This is telling redis to publish events since it's off by default.
	// https://redis.io/docs/manual/keyspace-notifications/
	_, err := client.DB.Do(context.Background(), "CONFIG", "SET", "notify-keyspace-events", "KEA").Result()
	if err != nil {
		logger.Log.Errorf("unable to set keyspace events %s", err.Error())
		os.Exit(1)
	}

	KeyEventChannel := fmt.Sprintf("__keyevent@%d__:expired", client.DBNumber)
	logger.AddData("key_event_channel", KeyEventChannel)

	// this is telling redis to subscribe to events published in the keyevent channel, specifically for expired events
	pubsub := client.DB.PSubscribe(context.Background(), KeyEventChannel)

	logger.Log.Infow("start service ", logger.Data()...)
	logger.ClearData()

	eventMapper := NewEventMapper()
	// Infinite loop for listening event
	for {
		logger.ClearData()
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
		dataEventStr := client.DB.Get(ctx, dataKey).Val()
		if dataEventStr == "" {
			logger.Log.Errorw("empty_dataEventStr", logger.Data()...)
			continue
		}
		// Delete data from resis
		if err = client.DB.Del(ctx, dataKey).Err(); err != nil {
			logger.Log.Errorw(err.Error(), logger.Data()...)
		}
		logger.Log.Infow("create_event", logger.Data()...)
		go eventMapper.CreateEvent(dataEventStr)
	}
}
