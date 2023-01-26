package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LukmanulHakim18/gorooster/v2/database"
	"github.com/LukmanulHakim18/gorooster/v2/helpers"
	"github.com/LukmanulHakim18/gorooster/v2/logger"
	"github.com/LukmanulHakim18/gorooster/v2/models"
	"github.com/LukmanulHakim18/gorooster/v2/services"

	"github.com/go-chi/chi"
)

var CLIENT_NAME = "X-CLIENT-NAME"

func GetEvent(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()

	redisClient := database.GetRedisClient()
	eventManager := services.GetServiceEventManager(redisClient)
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

	clientName := r.Header.Get(CLIENT_NAME)
	logger.AddData("client_name", clientName)
	if ok := helpers.ValidatorClinetNameAndKey(clientName); !ok {
		logger.Log.Errorw("error_client_name", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField(CLIENT_NAME))
		return
	}

	releaseEventFormat := r.Header.Get("X-RELEASE-FORMAT")
	releaseEventFormat = helpers.RelaseEventFormator(releaseEventFormat)

	logger.AddData("release_event_format", releaseEventFormat)

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
		Event: event,
	}
	res.SetEventRelease(releaseEventFormat, eventReleaseIn)

	helpers.ResponseSuccessWithData(w, http.StatusOK, res)
}

func CreateEventReleaseIn(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()

	redisClient := database.GetRedisClient()
	eventManager := services.GetServiceEventManager(redisClient)
	var (
		eventReleaseIn     time.Duration
		bodyEventReleaseIn helpers.BodyEventReleaseIn
		err                error
	)

	eventKey := chi.URLParam(r, "event_key")
	logger.AddData("event_key", eventKey)
	if ok := helpers.ValidatorClinetNameAndKey(eventKey); !ok {
		logger.Log.Errorw("error_event_key", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_key"))
		return
	}

	clientName := r.Header.Get(CLIENT_NAME)
	logger.AddData("client_name", clientName)
	if ok := helpers.ValidatorClinetNameAndKey(clientName); !ok {
		logger.Log.Errorw("error_client_name", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField(CLIENT_NAME))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&bodyEventReleaseIn)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadBody)
		return
	}
	logger.AddData("event", bodyEventReleaseIn.Event)

	if eventReleaseIn, err = time.ParseDuration(bodyEventReleaseIn.ReleaseIn); err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("release_in"))
		return
	}
	logger.AddData("event_release_in", eventReleaseIn)

	if err := eventManager.SetEventreleaseIn(clientName, eventKey, eventReleaseIn, bodyEventReleaseIn.Event); err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		if err.Error() == "duplicate key" {
			helpers.ResponseErrorWithData(w, helpers.ErrorDuplicateKey)
			return
		}
		helpers.ResponseErrorWithData(w, helpers.ErrorReadBody)
		return
	}
	logger.Log.Infow("success", logger.Data()...)

	res := helpers.SuccessResponse{
		Event:          bodyEventReleaseIn.Event,
		EventReleaseIn: bodyEventReleaseIn.ReleaseIn,
	}
	helpers.ResponseSuccessWithData(w, http.StatusCreated, res)
}

func UpdateReleaseEventIn(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()

	redisClient := database.GetRedisClient()
	eventManager := services.GetServiceEventManager(redisClient)
	var (
		eventReleaseIn     time.Duration
		event              models.Event
		err                error
		bodyEventReleaseIn helpers.BodyEventReleaseIn
	)

	eventKey := chi.URLParam(r, "event_key")
	logger.AddData("event_key", eventKey)
	if ok := helpers.ValidatorClinetNameAndKey(eventKey); !ok {
		logger.Log.Errorw("error_event_key", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_key"))
		return
	}
	logger.AddData("event_key", eventKey)

	clientName := r.Header.Get(CLIENT_NAME)
	logger.AddData("client_name", clientName)

	err = json.NewDecoder(r.Body).Decode(&bodyEventReleaseIn)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadBody)
		return
	}

	logger.AddData("release_in", eventReleaseIn)
	eventReleaseIn, err = time.ParseDuration(bodyEventReleaseIn.ReleaseIn)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_release_in"))
		return
	}
	logger.AddData("event_release_in", eventReleaseIn)

	err = eventManager.UpdateEventReleaseIn(clientName, eventKey, eventReleaseIn)
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
	eventManager := services.GetServiceEventManager(redisClient)
	var (
		bodyEventReleaseIn helpers.BodyEventReleaseIn
		eventReleaseIn     time.Duration
		event              models.Event
		err                error
	)

	eventKey := chi.URLParam(r, "event_key")
	if ok := helpers.ValidatorClinetNameAndKey(eventKey); !ok {
		logger.Log.Errorw("error_event_key", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_key"))
		return
	}
	logger.AddData("event_key", eventKey)

	clientName := r.Header.Get(CLIENT_NAME)
	if ok := helpers.ValidatorClinetNameAndKey(clientName); !ok {
		logger.Log.Errorw("error_client_name", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField(CLIENT_NAME))
		return
	}
	logger.AddData("client_name", clientName)

	err = json.NewDecoder(r.Body).Decode(&bodyEventReleaseIn)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadBody)
		return
	}
	event = bodyEventReleaseIn.Event
	logger.AddData("event", event)

	releaseEventFormat := r.Header.Get("X-RELEASE-FORMAT")
	releaseEventFormat = helpers.RelaseEventFormator(releaseEventFormat)
	logger.AddData("release_event_format", releaseEventFormat)

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
		Event: event,
	}
	res.SetEventRelease(releaseEventFormat, eventReleaseIn)
	helpers.ResponseSuccessWithData(w, http.StatusAccepted, res)

}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()

	redisClient := database.GetRedisClient()
	eventManager := services.GetServiceEventManager(redisClient)

	eventKey := chi.URLParam(r, "event_key")
	logger.AddData("event_key", eventKey)
	if ok := helpers.ValidatorClinetNameAndKey(eventKey); !ok {
		logger.Log.Errorw("error_event_key", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_key"))
		return
	}

	clientName := r.Header.Get(CLIENT_NAME)
	logger.AddData("client_name", clientName)
	if ok := helpers.ValidatorClinetNameAndKey(clientName); !ok {
		logger.Log.Errorw("error_client_name", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField(CLIENT_NAME))
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

	logger.Log.Infow("success", logger.Data()...)

	helpers.ResponseSuccessWithData(w, http.StatusNoContent, nil)

}

func CreateEventReleaseAt(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()

	redisClient := database.GetRedisClient()
	eventManager := services.GetServiceEventManager(redisClient)
	var (
		bodyEventReleaseAt helpers.BodyEventReleaseAt
		event              models.Event
		eventReleaseAt     time.Time
		err                error
	)

	eventKey := chi.URLParam(r, "event_key")
	logger.AddData("event_key", eventKey)
	if ok := helpers.ValidatorClinetNameAndKey(eventKey); !ok {
		logger.Log.Errorw("error_event_key", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_key"))
		return
	}

	clientName := r.Header.Get(CLIENT_NAME)
	logger.AddData("client_name", clientName)
	if ok := helpers.ValidatorClinetNameAndKey(clientName); !ok {
		logger.Log.Errorw("error_client_name", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField(CLIENT_NAME))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&bodyEventReleaseAt)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadBody)
		return
	}
	event = bodyEventReleaseAt.Event
	logger.AddData("event", event)

	eventReleaseAt = time.Unix(bodyEventReleaseAt.ReleaseAt, 0)
	if !eventReleaseAt.After(time.Now().Add(1 * time.Second)) {
		logger.Log.Errorw("error event_release_at", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorTimeReleaseAt)
		return
	}
	logger.AddData("event_release_at", eventReleaseAt)

	if err := eventManager.SetEventReleaseAt(clientName, eventKey, eventReleaseAt, event); err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		if err.Error() == "duplicate key" {
			helpers.ResponseErrorWithData(w, helpers.ErrorDuplicateKey)
			return
		}
		helpers.ResponseErrorWithData(w, helpers.ErrorReadBody)
		return
	}
	logger.Log.Infow("success", logger.Data()...)

	res := helpers.SuccessResponse{
		Event:          event,
		EventReleaseAt: eventReleaseAt.Unix(),
	}
	helpers.ResponseSuccessWithData(w, http.StatusCreated, res)
}

func UpdateReleaseEventAt(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetLogger()

	redisClient := database.GetRedisClient()
	eventManager := services.GetServiceEventManager(redisClient)
	var (
		bodyEventReleaseAt helpers.BodyEventReleaseAt
		eventReleaseAt     time.Time
		event              models.Event
		err                error
	)

	eventKey := chi.URLParam(r, "event_key")
	logger.AddData("event_key", eventKey)
	if ok := helpers.ValidatorClinetNameAndKey(eventKey); !ok {
		logger.Log.Errorw("error_event_key", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField("event_key"))
		return
	}

	clientName := r.Header.Get(CLIENT_NAME)
	logger.AddData("client_name", clientName)
	if ok := helpers.ValidatorClinetNameAndKey(clientName); !ok {
		logger.Log.Errorw("error_client_name", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadField(CLIENT_NAME))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&bodyEventReleaseAt)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadBody)
		return
	}

	eventReleaseAt = time.Unix(bodyEventReleaseAt.ReleaseAt, 0)
	if !eventReleaseAt.After(time.Now().Add(1 * time.Second)) {
		logger.Log.Errorw("error event_release_at", logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorTimeReleaseAt)
		return
	}
	logger.AddData("event_release_at", eventReleaseAt)

	if err := eventManager.UpdateEventReleaseAt(clientName, eventKey, eventReleaseAt); err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorReadBody)
		return
	}

	_, err = eventManager.GetEvent(clientName, eventKey, &event)
	if err != nil {
		logger.Log.Errorw(err.Error(), logger.Data()...)
		helpers.ResponseErrorWithData(w, helpers.ErrorServer)
		return
	}

	logger.AddData("event", event)

	logger.Log.Infow("success", logger.Data()...)
	res := helpers.SuccessResponse{
		Event:          event,
		EventReleaseAt: eventReleaseAt.Unix(),
	}
	helpers.ResponseSuccessWithData(w, http.StatusAccepted, res)

}
