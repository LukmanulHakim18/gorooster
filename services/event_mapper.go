package services

import (
	"context"
	"encoding/json"
	"time"

	"git.bluebird.id/mybb/gorooster/v2/helpers"
	"git.bluebird.id/mybb/gorooster/v2/logger"
	"git.bluebird.id/mybb/gorooster/v2/models"
	"git.bluebird.id/mybb/gorooster/v2/repositories"
	"github.com/go-redis/redis/v8"
)

type Mapper struct{}

func NewEventMapper() Mapper {
	return Mapper{}
}

// CreateEvent is function for build event
// From data string and formated to models.Event
func (m Mapper) CreateEvent(ctx context.Context, client *redis.Client, dataKey string) {
	logger := logger.GetLogger()

	defer m.Unlock(ctx, client, dataKey)
	// locking event for create event only this instence service
	if m.Lock(ctx, client, dataKey) {
		eventString := client.Get(ctx, dataKey).Val() // Get real data event from redis
		if eventString == "" {
			go logger.Log.Errorw("empty_dataEventStr", logger.Data()...)
			return
		}

		if err := client.Del(ctx, dataKey).Err(); err != nil { // Delete data from redis
			go logger.Log.Errorw(err.Error(), logger.Data()...)
		}
		logger.AddData("event_string", eventString)
		event := models.Event{}
		err := json.Unmarshal([]byte(eventString), &event)
		if err != nil {
			go logger.Log.Errorw(err.Error(), logger.Data()...)
			return
		}
		logger.AddData("event_struct", event)

		// Initiate repo interface
		var jobRepo repositories.Contract

		// Set which repository to use
		switch event.Type {
		case models.API_EVENT:
			jobRepo = repositories.NewJobAPI()
		default:
			go logger.Log.Errorw("unknown event type", logger.Data()...)
			return
		}

		// Run the command according to the repo type
		// That has been set above
		if err = jobRepo.DoJob(eventString); err != nil {
			go logger.Log.Errorw(err.Error(), logger.Data()...)
			return
		}
		go logger.Log.Infow("successfully do job", logger.Data()...)
		logger.ClearData()
	} else {

		logger.Log.Errorw("event run in another instance", logger.Data()...)
	}
}

func (m Mapper) Lock(ctx context.Context, client *redis.Client, dataKey string) bool {
	return client.SetNX(ctx, helpers.GenerateKeyLocking(dataKey), "locking", 1*time.Minute).Val()
}
func (m Mapper) Unlock(ctx context.Context, client *redis.Client, dataKey string) error {
	return client.Del(ctx, helpers.GenerateKeyLocking(dataKey)).Err()
}
