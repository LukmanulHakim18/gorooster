package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"git.bluebird.id/mybb/gorooster/database"
	"git.bluebird.id/mybb/gorooster/helpers"
	"git.bluebird.id/mybb/gorooster/models"
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
	cmd := res.DBIndex.Get(ctx, helpers.GenerateKeyIndex(clientName, key))
	err = cmd.Err()
	if err != nil {
		return 0, fmt.Errorf("data not found")
	}
	dbAddress, err := cmd.Int()
	if err != nil {
		return 0, fmt.Errorf("data not found")
	}
	val := res.DB[dbAddress].Get(ctx, helpers.GenerateKeyData(clientName, key)).Val()

	if err = json.Unmarshal([]byte(val), target); err != nil {
		return
	}
	eventReleaseIn, err = res.DB[dbAddress].TTL(ctx, helpers.GenerateKeyEvent(clientName, key)).Result()
	return
}

func (res EventManaget) SetEvent(clientName, key string, expired time.Duration, event models.Event) error {
	ctx := context.Background()
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := res.DBIndex.Set(ctx, helpers.GenerateKeyIndex(clientName, key), res.DBPointer, expired).Err(); err != nil {
		return err
	}
	if err := res.DB[res.DBPointer].Set(ctx, helpers.GenerateKeyEvent(clientName, key), "event-key", expired).Err(); err != nil {
		return err
	}
	if err := res.DB[res.DBPointer].Set(ctx, helpers.GenerateKeyData(clientName, key), string(data), -1).Err(); err != nil {
		return err
	}
	res.Next()
	return nil
}

func (res EventManaget) UpdateExpiredEvent(clientName, key string, expired time.Duration) error {
	ctx := context.Background()
	cmd := res.DBIndex.Get(ctx, helpers.GenerateKeyIndex(clientName, key))
	err := cmd.Err()
	if err != nil {
		return fmt.Errorf("data not found")
	}
	dbAddress, err := cmd.Int()
	if err != nil {
		return fmt.Errorf("data not found")
	}
	if err := res.DBIndex.Set(ctx, helpers.GenerateKeyIndex(clientName, key), "event-key", expired).Err(); err != nil {
		return err
	}
	if err := res.DB[dbAddress].Set(ctx, helpers.GenerateKeyEvent(clientName, key), "event-key", expired).Err(); err != nil {
		return err
	}
	return nil
}

func (res EventManaget) UpdateDataEvent(clientName, key string, event models.Event) error {
	ctx := context.Background()
	cmd := res.DBIndex.Get(ctx, helpers.GenerateKeyIndex(clientName, key))
	err := cmd.Err()
	if err != nil {
		return fmt.Errorf("data not found")
	}
	dbAddress, err := cmd.Int()
	if err != nil {
		return fmt.Errorf("data not found")
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := res.DB[dbAddress].Get(ctx, helpers.GenerateKeyData(clientName, key)).Err(); err != nil {
		return err
	}
	return res.DB[dbAddress].Set(ctx, helpers.GenerateKeyData(clientName, key), data, -1).Err()

}

func (res EventManaget) DeleteEvent(clientName, key string) error {
	ctx := context.Background()
	cmd := res.DBIndex.Get(ctx, helpers.GenerateKeyIndex(clientName, key))
	err := cmd.Err()
	if err != nil {
		return fmt.Errorf("data not found")
	}
	dbAddress, err := cmd.Int()
	if err != nil {
		return fmt.Errorf("data not found")
	}
	if err := res.DBIndex.Del(ctx, helpers.GenerateKeyIndex(clientName, key)).Err(); err != nil {
		return err
	}
	if err := res.DB[dbAddress].Del(ctx, helpers.GenerateKeyEvent(clientName, key)).Err(); err != nil {
		return err
	}
	if err := res.DB[dbAddress].Del(ctx, helpers.GenerateKeyData(clientName, key)).Err(); err != nil {
		return err
	}
	return nil
}
