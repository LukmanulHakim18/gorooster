package repositories

import (
	"bytes"
	"encoding/json"
	"event-scheduler/helpers"
	"event-scheduler/models"
	"fmt"
	"net/http"
)

type jobAPIRepository struct {
	jobEvent models.JobAPI
}

func NewJobAPI() Contract {
	return jobAPIRepository{}
}

// execute the http request task built from the job api event
func (jar jobAPIRepository) DoJob(evenString string) error {

	if err := helpers.ExtractEvent(evenString, &jar.jobEvent); err != nil {
		// log error
		return err
	}
	if ok := jar.jobEvent.IsMethodSupport(); !ok {
		// log error
		return fmt.Errorf("mthod not support")
	}
	if err := jar.post(); err != nil {
		return err
	}
	return nil

}

func (jar jobAPIRepository) post() error {

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

	// Success is indicated with 2xx status codes:
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return fmt.Errorf("error status code %d", resp.StatusCode)
	}
 	return nil
}

func (jar jobAPIRepository) setHeaders(req *http.Request) {
	if len(jar.jobEvent.Headers) > 0 {
		for _, header := range jar.jobEvent.Headers {
			req.Header.Set(header.Key, header.Value)
		}
	}
}
