package services

import (
	"encoding/json"
	"github.com/LukmanulHakim18/gorooster/logger"
	"github.com/LukmanulHakim18/gorooster/models"
	"github.com/LukmanulHakim18/gorooster/repositories"
)

type Mapper struct{}

func NewEventMapper() Mapper {
	return Mapper{}
}

// CreateEvent is function for build event
// From data string and formated to models.Event
func (m Mapper) CreateEvent(eventString string) {
	logger := logger.GetLogger()
	logger.AddData("event_string", eventString)
	event := models.Event{}
	err := json.Unmarshal([]byte(eventString), &event)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
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
		logger.Log.Errorw("unknown event type", logger.Data()...)
		return
	}

	// Run the command according to the repo type
	// That has been set above
	if err = jobRepo.DoJob(eventString); err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		return
	}
	logger.Log.Infow("successfully do job", logger.Data()...)
}
