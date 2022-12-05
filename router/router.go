package router

import (
	"github.com/LukmanulHakim18/gorooster/controllers"
	"net/http"

	"github.com/go-chi/chi"
)

func GetRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", HealtCheck)
	r.Post("/event/{event_key}/{event_release_in}", controllers.CreateEvent)       // Save event
	r.Get("/event/{event_key}", controllers.GetEvent)                              // Get event
	r.Put("/event/{event_key}/{event_release_in}", controllers.UpdateReleaseEvent) // Update expired event
	r.Put("/event/{event_key}", controllers.UpdateDataEvent)                       // Update data event
	r.Delete("/event/{event_key}", controllers.DeleteEvent)                        // Delete event
	return r
}
func HealtCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
