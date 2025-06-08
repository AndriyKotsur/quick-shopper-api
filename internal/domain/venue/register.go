package venue

import (
	"github.com/go-chi/chi"

	"github.com/AndriyKotsur/quick-booking-api/internal/database"
	"github.com/AndriyKotsur/quick-booking-api/internal/middleware"
)

func RegisterVenueEndpoints(r *chi.Mux, db *database.Queries) {
	venueController := NewVenueController(db)

	r.Route("/venue", func(r chi.Router) {
		r.Use(middleware.AuthenticateUser)

		r.Post("/", venueController.CreateVenue)
		r.Get("/", venueController.GetVenues)
		r.Get("/{id}", venueController.GetVenueByID)
		r.Put("/{id}", venueController.UpdateVenue)
		r.Delete("/{id}", venueController.DeleteVenue)
	})
}
