package models

type EventType string

const (
	API_EVENT = "api_event"
)

type Event struct {
	Name    string      `json:"name"`
	Id      string      `json:"id"`
	Type    EventType   `json:"type"`
	JobData interface{} `json:"job_data"`
}

func (e Event) IsTypeSupport() bool {
	return e.Type == API_EVENT
}
