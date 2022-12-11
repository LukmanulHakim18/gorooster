package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"git.bluebird.id/mybb/gorooster/v2/helpers"
	"git.bluebird.id/mybb/gorooster/v2/logger"
	"git.bluebird.id/mybb/gorooster/v2/models"
)

type jobAPIRepository struct {
	jobEvent models.JobAPI
}

func NewJobAPI() Contract {
	return jobAPIRepository{}
}

// Runing job from data event
// Validate Event
func (jar jobAPIRepository) DoJob(eventString string) (err error) {
	event := models.Event{
		JobData: &jar.jobEvent,
	}
	if err := json.Unmarshal([]byte(eventString), &event); err != nil {
		return err
	}
	if err = jar.jobEvent.Validate(); err != nil {
		return err
	}

	// Get config from env file, is in retry mode
	retryMode := helpers.EnvGetBool("RETRY_MODE", false)
	if retryMode {
		retryCount := helpers.EnvGetInt("RETRY_COUNT", 3)
		return jar.retryManager(event.Id, retryCount)
	}
	return jar.sendRequest()
}

// Do retry n time if error send request
func (jar jobAPIRepository) retryManager(eventId string, retryCount int) (err error) {
	for i := 1; i <= retryCount; i++ {
		go logger.GetLogger().Log.Infof("%d attempt to send event with request id: %s", i, eventId)
		err = jar.sendRequest()
		if err == nil {
			return nil
		}
	}
	return err
}

// Build Http request from Event.Data
func (jar jobAPIRepository) sendRequest() error {
	logger := logger.GetLogger()
	payloadBytes, err := json.Marshal(jar.jobEvent.Data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(jar.jobEvent.Method.ToString(), jar.jobEvent.Endpoint, body)
	if err != nil {
		return err
	}
	// setup headers
	jar.setHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	logger.AddData("status_code", resp.StatusCode)
	// Success is indicated with 2xx status codes:
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return fmt.Errorf("error status code %d", resp.StatusCode)
	}

	var res map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		go logger.Log.Error("error read response body")
	}
	logger.AddData("response_body", res)

	go logger.Log.Infow("success_hit_endpoint", logger.Data()...)
	return nil
}

func (jar jobAPIRepository) setHeaders(req *http.Request) {
	if len(jar.jobEvent.Headers) > 0 {
		for _, header := range jar.jobEvent.Headers {
			req.Header.Set(header.Key, header.Value)
		}
	}
}
