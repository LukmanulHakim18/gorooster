package router

import (
	"net/http"

	"github.com/LukmanulHakim18/gorooster/v2/controllers"

	"github.com/go-chi/chi"
)

func GetRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", HealthCheck) // health check endpoint

	r.Put("/event/{event_key}", controllers.UpdateDataEvent) // Update data event
	r.Delete("/event/{event_key}", controllers.DeleteEvent)  // Delete event

	r.Route("/event", func(r chi.Router) {
		r.Route("/relin", func(r chi.Router) {
			r.Post("/{event_key}", controllers.CreateEventReleaseIn) // Save event release in
			r.Put("/{event_key}", controllers.UpdateReleaseEventIn)  // Update event release in
		})
		r.Route("/relat", func(r chi.Router) {
			r.Post("/{event_key}", controllers.CreateEventReleaseAt) // Save event release at
			r.Put("/{event_key}", controllers.UpdateReleaseEventAt)  // Update event release at
		})
		r.Get("/{event_key}", controllers.GetEvent)        // Get event
		r.Put("/{event_key}", controllers.UpdateDataEvent) // Update data event
		r.Delete("/{event_key}", controllers.DeleteEvent)  // Delete event
	})
	return r
}
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
