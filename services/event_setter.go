package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gorooster/database"
	"gorooster/helpers"
	"gorooster/models"
	"time"
)

type EventManaget struct {
	*database.RedisClient
}

func GetServiceEventManaget(redisClient *database.RedisClient) EventManaget {
	return EventManaget{
		redisClient,
	}
}

func (res EventManaget) GetEvent(clientName, key string, target interface{}) (eventReleaseIn time.Duration, err error) {
	ctx := context.Background()
	err = res.DB.Get(ctx, helpers.GenerateKeyEvent(clientName, key)).Err()
	if err != nil {
		return 0, fmt.Errorf("data not found")
	}

	val := res.DB.Get(ctx, helpers.GenerateKeyData(clientName, key)).Val()

	if err = json.Unmarshal([]byte(val), target); err != nil {
		return
	}
	eventReleaseIn, err = res.DB.TTL(ctx, helpers.GenerateKeyEvent(clientName, key)).Result()
	return
}

func (res EventManaget) SetEvent(clientName, key string, expired time.Duration, event models.Event) error {
	ctx := context.Background()
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := res.DB.Set(ctx, helpers.GenerateKeyEvent(clientName, key), "event-key", expired).Err(); err != nil {
		return err
	}
	if err := res.DB.Set(ctx, helpers.GenerateKeyData(clientName, key), string(data), -1).Err(); err != nil {
		return err
	}
	return nil
}

func (res EventManaget) UpdateExpiredEvent(clientName, key string, expired time.Duration) error {
	ctx := context.Background()
	err := res.DB.Get(ctx, helpers.GenerateKeyEvent(clientName, key)).Err()
	if err != nil {
		return fmt.Errorf("data not found")
	}
	if err := res.DB.Set(ctx, helpers.GenerateKeyEvent(clientName, key), "event-key", expired).Err(); err != nil {
		return err
	}
	return nil
}

func (res EventManaget) UpdateDataEvent(clientName, key string, event models.Event) error {
	ctx := context.Background()
	err := res.DB.Get(ctx, helpers.GenerateKeyEvent(clientName, key)).Err()
	if err != nil {
		return fmt.Errorf("data not found")
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := res.DB.Get(ctx, helpers.GenerateKeyData(clientName, key)).Err(); err != nil {
		return err
	}
	return res.DB.Set(ctx, helpers.GenerateKeyData(clientName, key), data, -1).Err()

}

func (res EventManaget) DeleteEvent(clientName, key string) error {
	if ok := helpers.ValidatorClinetNameAndKey(key); !ok {
		return fmt.Errorf("key can not contain ':'")
	}
	ctx := context.Background()
	err := res.DB.Get(ctx, helpers.GenerateKeyEvent(clientName, key)).Err()
	if err != nil {
		return fmt.Errorf("data not found")
	}
	if err := res.DB.Del(ctx, helpers.GenerateKeyEvent(clientName, key)).Err(); err != nil {
		return err
	}
	if err := res.DB.Del(ctx, helpers.GenerateKeyData(clientName, key)).Err(); err != nil {
		return err
	}
	return nil
}
