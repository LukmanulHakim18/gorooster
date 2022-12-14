package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LukmanulHakim18/gorooster/v2/models"
)

func ResponseSuccessWithData(w http.ResponseWriter, statusCode int, payload any) {
	dataByte, err := json.Marshal(payload)
	if err != nil {
		ResponseErrorServer(w)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dataByte)
}

func ResponseErrorServer(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	dataByte, _ := json.Marshal(ErrorServer)
	w.WriteHeader(500)
	w.Write(dataByte)
}

func ResponseErrorWithData(w http.ResponseWriter, errFmt *Error) {
	dataByte, err := json.Marshal(errFmt)
	if err != nil {
		ResponseErrorServer(w)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errFmt.StatusCode)
	w.Write(dataByte)
}

// ================================ success format ================================
type SuccessResponse struct {
	Event          models.Event `json:"event"`
	EventReleaseIn string       `json:"event_release_in,omitempty"`
	EventReleaseAt int64        `json:"event_release_at,omitempty"`
}

func (sr *SuccessResponse) SetEventRelease(releaseEventFormat string, eventReleaseIn time.Duration) {
	switch releaseEventFormat {
	case ReleaseEventAT:
		sr.EventReleaseAt = time.Now().Add(eventReleaseIn).Unix()
	case ReleaseEventIN:
		sr.EventReleaseIn = eventReleaseIn.String()
	default:
		sr.EventReleaseAt = time.Now().Add(eventReleaseIn).Unix()
		sr.EventReleaseIn = eventReleaseIn.String()
	}
}

// ================================ error format ================================
type Error struct {
	StatusCode       int                    `json:"-"`
	ErrorCode        string                 `json:"error_code"`
	ErrorMessage     string                 `json:"error_message"`
	ErrorField       string                 `json:"error_field,omitempty"`
	LocalizedMessage Message                `json:"localized_message"`
	Data             map[string]interface{} `json:"data,omitempty"`
	ErrorData        interface{}            `json:"error_data,omitempty"`
}
type Message struct {
	English   string `json:"en"`
	Indonesia string `json:"id"`
}

func (err Error) Error() string {
	return err.ErrorMessage
}

func ErrorReadField(field string) *Error {
	return &Error{
		StatusCode:   http.StatusBadRequest,
		ErrorCode:    "CROW-400",
		ErrorMessage: fmt.Sprintf("there is an error in the %s parameter", field),
		LocalizedMessage: Message{
			English:   fmt.Sprintf("there is an error in the %s parameter", field),
			Indonesia: fmt.Sprintf("ada kesalahan pada parameter %s", field),
		},
	}
}

func ErrorDataNotFound(data string) *Error {
	return &Error{
		StatusCode:   http.StatusNotFound,
		ErrorCode:    "CROW-404",
		ErrorMessage: fmt.Sprintf("%s not Found", data),
		LocalizedMessage: Message{
			English:   fmt.Sprintf("%s not Found", data),
			Indonesia: fmt.Sprintf("%s tidak dapat ditemukan", data),
		},
	}
}

var ErrorServer = &Error{
	StatusCode:   http.StatusInternalServerError,
	ErrorCode:    "CROW-500",
	ErrorMessage: "failed to process the request",
	LocalizedMessage: Message{
		English:   "failed to process the request",
		Indonesia: "gagal memperoses permintaan",
	},
}

var ErrorReadBody = &Error{
	StatusCode:   http.StatusInternalServerError,
	ErrorCode:    "CROW-001",
	ErrorMessage: "failed to read data",
	LocalizedMessage: Message{
		English:   "failed to read data",
		Indonesia: "gagal membaca data",
	},
}
var ErrorTimeReleaseAt = &Error{
	StatusCode:   http.StatusBadRequest,
	ErrorCode:    "CROW-401",
	ErrorMessage: "unable to release past events.",
	LocalizedMessage: Message{
		English:   "unable to release past events.",
		Indonesia: "tidak dapat merilis event masa lalu.",
	},
}

var ErrorDuplicateKey = &Error{
	StatusCode:   http.StatusBadRequest,
	ErrorCode:    "CROW-402",
	ErrorMessage: "duplicat key.",
	LocalizedMessage: Message{
		English:   "duplicat key.",
		Indonesia: "kunci duplikat.",
	},
}
