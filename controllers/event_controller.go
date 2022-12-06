package controllers

import (
	"encoding/json"
	"git.bluebird.id/mybb/gorooster/database"
	"git.bluebird.id/mybb/gorooster/helpers"
	"git.bluebird.id/mybb/gorooster/logger"
	"git.bluebird.id/mybb/gorooster/models"
	"git.bluebird.id/mybb/gorooster/services"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func GetEvent(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()

	redisClient := database.GetRedisClient()
	eventManager := services.GetServiceEventManaget(redisClient)
	var (
		eventReleaseIn time.Duration
		event          models.Event
	)

	eventKey := chi.URLParam(r, "event_key")
	logger.AddData("event_key", eventKey)
	if ok := helpers.ValidatorClinetNameAndKey(eventKey); !ok {
		logger.Log.Errorw("error_event_key", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_key"))
		return
	}

	clientName := r.Header.Get("X-CLIENT-NAME")
	logger.AddData("client_name", clientName)
	if ok := helpers.ValidatorClinetNameAndKey(clientName); !ok {
		logger.Log.Errorw("error_client_name", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("client_name"))
		return
	}

	eventReleaseIn, err := eventManager.GetEvent(clientName, eventKey, &event)
	if err != nil {
		if err.Error() == "data not found" {
			logger.Log.Errorw(err.Error(), logger.Data()...)
			helpers.ResponseErrorWithData(w, helpers.ErrorDataNotFound("event"))
			return
		}
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorServer)
		return
	}
	logger.AddData("event", event)
	logger.AddData("event_release_in", eventReleaseIn)

	logger.Log.Infow("success", logger.Data()...)
	res := helpers.SuccessResponse{
		Event:          event,
		EventReleaseIn: eventReleaseIn.String(),
	}
	helpers.ResponseSuccessWithData(w, http.StatusOK, res)

}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()

	redisClient := database.GetRedisClient()
	eventManager := services.GetServiceEventManaget(redisClient)
	var (
		eventReleaseIn time.Duration
		event          models.Event
		err            error
	)

	eventKey := chi.URLParam(r, "event_key")
	logger.AddData("event_key", eventKey)
	if ok := helpers.ValidatorClinetNameAndKey(eventKey); !ok {
		logger.Log.Errorw("error_event_key", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_key"))
		return
	}

	clientName := r.Header.Get("X-CLIENT-NAME")
	logger.AddData("client_name", clientName)
	if ok := helpers.ValidatorClinetNameAndKey(clientName); !ok {
		logger.Log.Errorw("error_client_name", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("client_name"))
		return
	}

	eventReleaseInStr := chi.URLParam(r, "event_release_in")
	if eventReleaseIn, err = time.ParseDuration(eventReleaseInStr); err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_release_in"))
		return
	}
	logger.AddData("event_release_in", eventReleaseIn)

	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadBody)
		return
	}
	logger.AddData("event", event)

	if err := eventManager.SetEvent(clientName, eventKey, eventReleaseIn, event); err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadBody)
		return
	}
	logger.Log.Infow("success", logger.Data()...)

	res := helpers.SuccessResponse{
		Event:          event,
		EventReleaseIn: eventReleaseIn.String(),
	}
	helpers.ResponseSuccessWithData(w, http.StatusCreated, res)
}

func UpdateReleaseEvent(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()

	redisClient := database.GetRedisClient()
	eventManager := services.GetServiceEventManaget(redisClient)
	var (
		eventReleaseIn time.Duration
		event          models.Event
		err            error
	)

	eventKey := chi.URLParam(r, "event_key")
	logger.AddData("event_key", eventKey)
	if ok := helpers.ValidatorClinetNameAndKey(eventKey); !ok {
		logger.Log.Errorw("error_event_key", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_key"))
		return
	}
	logger.AddData("event_key", eventKey)

	clientName := r.Header.Get("X-CLIENT-NAME")
	eventReleaseInStr := chi.URLParam(r, "event_release_in")
	if ok := helpers.ValidatorClinetNameAndKey(clientName); !ok {
		logger.Log.Errorw("error_client_name", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("client_name"))
		return
	}
	logger.AddData("client_name", clientName)

	eventReleaseIn, err = time.ParseDuration(eventReleaseInStr)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_release_in"))
		return
	}
	logger.AddData("event_release_in", eventReleaseIn)

	err = eventManager.UpdateExpiredEvent(clientName, eventKey, eventReleaseIn)
	if err != nil {
		if err.Error() == "data not found" {
			logger.Log.Errorw(err.Error(), logger.Data()...)
			helpers.ResponseErrorWithData(w, helpers.ErrorDataNotFound("event"))
			return
		}
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorServer)
		return
	}

	eventReleaseIn, err = eventManager.GetEvent(clientName, eventKey, &event)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorServer)
		return
	}

	logger.AddData("event", event)

	logger.Log.Infow("success", logger.Data()...)
	res := helpers.SuccessResponse{
		Event:          event,
		EventReleaseIn: eventReleaseIn.String(),
	}
	helpers.ResponseSuccessWithData(w, http.StatusAccepted, res)

}
func UpdateDataEvent(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()

	redisClient := database.GetRedisClient()
	eventManager := services.GetServiceEventManaget(redisClient)
	var (
		eventReleaseIn time.Duration
		event          models.Event
	)

	eventKey := chi.URLParam(r, "event_key")
	logger.AddData("event_key", eventKey)
	clientName := r.Header.Get("X-CLIENT-NAME")
	logger.AddData("client_name", clientName)

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadBody)
		return
	}
	logger.AddData("event", event)

	if ok := helpers.ValidatorClinetNameAndKey(eventKey); !ok {
		logger.Log.Errorw("error_event_key", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_key"))
		return
	}
	if ok := helpers.ValidatorClinetNameAndKey(clientName); !ok {
		logger.Log.Errorw("error_client_name", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("client_name"))
		return
	}

	err = eventManager.UpdateDataEvent(clientName, eventKey, event)
	if err != nil {
		if err.Error() == "data not found" {
			logger.Log.Errorw(err.Error(), logger.Data()...)
			helpers.ResponseErrorWithData(w, helpers.ErrorDataNotFound("event"))
			return
		}
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorServer)
		return
	}

	eventReleaseIn, err = eventManager.GetEvent(clientName, eventKey, &event)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorServer)
		return
	}

	logger.AddData("event", event)
	logger.AddData("event_release_in", eventReleaseIn)

	logger.Log.Infow("success", logger.Data()...)
	res := helpers.SuccessResponse{
		Event:          event,
		EventReleaseIn: eventReleaseIn.String(),
	}
	helpers.ResponseSuccessWithData(w, http.StatusAccepted, res)

}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()

	redisClient := database.GetRedisClient()
	eventManager := services.GetServiceEventManaget(redisClient)
	var (
		eventReleaseIn time.Duration
	)

	eventKey := chi.URLParam(r, "event_key")
	logger.AddData("event_key", eventKey)
	if ok := helpers.ValidatorClinetNameAndKey(eventKey); !ok {
		logger.Log.Errorw("error_event_key", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_key"))
		return
	}

	clientName := r.Header.Get("X-CLIENT-NAME")
	logger.AddData("client_name", clientName)
	if ok := helpers.ValidatorClinetNameAndKey(clientName); !ok {
		logger.Log.Errorw("error_client_name", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("client_name"))
		return
	}

	err := eventManager.DeleteEvent(clientName, eventKey)
	if err != nil {
		if err.Error() == "data not found" {
			logger.Log.Errorw(err.Error(), logger.Data()...)
			helpers.ResponseErrorWithData(w, helpers.ErrorDataNotFound("event"))
			return
		}
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorServer)
		return
	}
	logger.AddData("event_release_in", eventReleaseIn)

	logger.Log.Infow("success", logger.Data()...)

	helpers.ResponseSuccessWithData(w, http.StatusNoContent, nil)

}
