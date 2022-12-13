package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/LukmanulHakim18/gorooster/v2/database"
	"github.com/LukmanulHakim18/gorooster/v2/helpers"
	"github.com/LukmanulHakim18/gorooster/v2/models"
)

type EventManager struct {
	*database.RedisClient
	wipeDataEvent time.Duration
}

func GetServiceEventManager(redisClient *database.RedisClient) EventManager {
	return EventManager{
		redisClient,
		helpers.EnvGetTimeDuration("WIPE_DATA_EVENT", 24*time.Hour),
	}
}

func (em EventManager) GetEvent(clientName, key string, target interface{}) (eventReleaseIn time.Duration, err error) {
	ctx := context.Background()
	dbAddress, err := em.DBIndex.Get(ctx, helpers.GenerateKeyIndex(clientName, key)).Int()
	if err != nil {
		return 0, fmt.Errorf("data not found")
	}
	val := em.DB[dbAddress].Get(ctx, helpers.GenerateKeyData(clientName, key)).Val()

	if err = json.Unmarshal([]byte(val), target); err != nil {
		return
	}
	eventReleaseIn, err = em.DB[dbAddress].TTL(ctx, helpers.GenerateKeyEvent(clientName, key)).Result()
	return
}

func (em EventManager) SetEventreleaseIn(clientName, key string, releaseIn time.Duration, event models.Event) error {
	ctx := context.Background()
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	cmd := em.DBIndex.SetNX(ctx, helpers.GenerateKeyIndex(clientName, key), em.DBPointer, releaseIn)
	if cmd.Err(); err != nil {
		return err
	}
	if !cmd.Val() {
		return fmt.Errorf("duplicate key")
	}

	if err := em.DB[em.DBPointer].Set(ctx, helpers.GenerateKeyEvent(clientName, key), "event-key", releaseIn).Err(); err != nil {
		return err
	}

	if err := em.DB[em.DBPointer].Set(ctx, helpers.GenerateKeyData(clientName, key), string(data), releaseIn+em.wipeDataEvent).Err(); err != nil {
		return err
	}

	em.Next()
	return nil
}

func (em EventManager) UpdateEventReleaseIn(clientName, key string, releaseIn time.Duration) error {
	ctx := context.Background()
	dbAddress, err := em.DBIndex.Get(ctx, helpers.GenerateKeyIndex(clientName, key)).Int()
	if err != nil {
		return fmt.Errorf("data not found")
	}

	if err := em.DBIndex.Set(ctx, helpers.GenerateKeyIndex(clientName, key), dbAddress, releaseIn).Err(); err != nil {
		return err
	}

	if err := em.DB[dbAddress].Set(ctx, helpers.GenerateKeyEvent(clientName, key), "event-key", releaseIn).Err(); err != nil {
		return err
	}
	if err := em.DB[dbAddress].Expire(ctx, helpers.GenerateKeyData(clientName, key), releaseIn+em.wipeDataEvent).Err(); err != nil {
		return err
	}
	return nil
}

func (em EventManager) UpdateDataEvent(clientName, key string, event models.Event) error {
	ctx := context.Background()
	dbAddress, err := em.DBIndex.Get(ctx, helpers.GenerateKeyIndex(clientName, key)).Int()
	if err != nil {
		return fmt.Errorf("data not found")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if err := em.DB[dbAddress].Get(ctx, helpers.GenerateKeyData(clientName, key)).Err(); err != nil {
		return err
	}

	return em.DB[dbAddress].Set(ctx, helpers.GenerateKeyData(clientName, key), data, -1).Err()

}

func (em EventManager) DeleteEvent(clientName, key string) error {
	ctx := context.Background()
	dbAddress, err := em.DBIndex.Get(ctx, helpers.GenerateKeyIndex(clientName, key)).Int()
	if err != nil {
		return fmt.Errorf("data not found")
	}

	if err := em.DBIndex.Del(ctx, helpers.GenerateKeyIndex(clientName, key)).Err(); err != nil {
		return err
	}

	if err := em.DB[dbAddress].Del(ctx, helpers.GenerateKeyEvent(clientName, key)).Err(); err != nil {
		return err
	}

	if err := em.DB[dbAddress].Del(ctx, helpers.GenerateKeyData(clientName, key)).Err(); err != nil {
		return err
	}

	return nil
}

func (em EventManager) SetEventReleaseAt(clientName, key string, releaseAt time.Time, event models.Event) error {
	ctx := context.Background()
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	cmd := em.DBIndex.SetNX(ctx, helpers.GenerateKeyIndex(clientName, key), em.DBPointer, em.wipeDataEvent)
	if cmd.Err(); err != nil {
		return err
	}
	if !cmd.Val() {
		return fmt.Errorf("duplicate key")
	}

	if err := em.DBIndex.ExpireAt(ctx, helpers.GenerateKeyIndex(clientName, key), releaseAt).Err(); err != nil {
		return err
	}

	if err := em.DB[em.DBPointer].Set(ctx, helpers.GenerateKeyEvent(clientName, key), "event-key", em.wipeDataEvent).Err(); err != nil {
		return err
	}

	if err := em.DB[em.DBPointer].ExpireAt(ctx, helpers.GenerateKeyEvent(clientName, key), releaseAt).Err(); err != nil {
		return err
	}

	if err := em.DB[em.DBPointer].Set(ctx, helpers.GenerateKeyData(clientName, key), string(data), em.wipeDataEvent).Err(); err != nil {
		return err
	}

	if err := em.DB[em.DBPointer].ExpireAt(ctx, helpers.GenerateKeyData(clientName, key), releaseAt.Add(em.wipeDataEvent)).Err(); err != nil {
		return err
	}

	em.Next()
	return nil
}

func (em EventManager) UpdateEventReleaseAt(clientName, key string, releaseAt time.Time) error {
	ctx := context.Background()
	dbAddress, err := em.DBIndex.Get(ctx, helpers.GenerateKeyIndex(clientName, key)).Int()
	if err != nil {
		return fmt.Errorf("data not found")
	}

	if err := em.DBIndex.ExpireAt(ctx, helpers.GenerateKeyIndex(clientName, key), releaseAt).Err(); err != nil {
		return err
	}

	if err := em.DB[dbAddress].ExpireAt(ctx, helpers.GenerateKeyEvent(clientName, key), releaseAt).Err(); err != nil {
		return err
	}
	if err := em.DB[dbAddress].ExpireAt(ctx, helpers.GenerateKeyData(clientName, key), releaseAt.Add(em.wipeDataEvent)).Err(); err != nil {
		return err
	}

	return nil
}
