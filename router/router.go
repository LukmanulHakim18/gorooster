package router

import (
	"net/http"

	"git.bluebird.id/mybb/gorooster/v2/controllers"

	"github.com/go-chi/chi"
)

func GetRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", HealtCheck)
	r.Post("/event/{event_key}/{event_release_in}", controllers.CreateEvent)                    // Save event
	r.Get("/event/{event_key}", controllers.GetEvent)                                           // Get event
	r.Put("/event/{event_key}/{event_release_in}", controllers.UpdateReleaseEvent)              // Update expired event
	r.Put("/event/{event_key}", controllers.UpdateDataEvent)                                    // Update data event
	r.Delete("/event/{event_key}", controllers.DeleteEvent)                                     // Delete event
	r.Post("/event/release_at/{event_key}/{event_release_at}", controllers.CreateEventAt)       // Save event with release at
	r.Put("/event/release_at/{event_key}/{event_release_at}", controllers.UpdateReleaseEventAt) // Update event with release at
	return r
}
func HealtCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
