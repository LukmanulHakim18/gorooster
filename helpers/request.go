package helpers

import (
	"github.com/LukmanulHakim18/gorooster/v2/models"
)

type BodyEventReleaseIn struct {
	Event     models.Event `json:"event"`
	ReleaseIn string       `json:"release_in"`
}
type BodyEventReleaseAt struct {
	Event     models.Event `json:"event"`
	ReleaseAt int64        `json:"release_at"`
}
