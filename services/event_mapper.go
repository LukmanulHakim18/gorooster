package services

import (
	"encoding/json"
	"event-scheduler/logger"
	"event-scheduler/models"
	"event-scheduler/repositories"
)

type Mapper struct{}

func NewEventMapper() Mapper {
	return Mapper{}
}

func (m Mapper) CreateEvent(eventString string) {
	logger := logger.GetLogger()
	logger.AddData("eventString", eventString)
	event := models.Event{}
	err := json.Unmarshal([]byte(eventString), &event)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		return
	}
	// initiate repo interface
	logger.AddData("event", event)
	var jobRepo repositories.Contract

	// set which repository to use
	switch event.Type {
	case models.API_EVENT:
		jobRepo = repositories.NewJobAPI()
	default:
		logger.Log.Errorw("Unknown event type", logger.Data()...)
		return
	}

	// run the command according to the repo type
	// that has been set above
	if err = jobRepo.DoJob(eventString); err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		return
	}
	logger.Log.Infow("Success", logger.Data()...)

}
