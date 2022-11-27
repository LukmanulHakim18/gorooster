package helpers

import (
	"encoding/json"
	"event-scheduler/models"
)

func ExtractEvent(eventString string, target any) error {
	event := models.Event{
		JobData: target,
	}
	if err := json.Unmarshal([]byte(eventString), &event); err != nil {
		return err
	}
	return nil
}
