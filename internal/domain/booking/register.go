package booking

import (
	"github.com/go-chi/chi"

	"github.com/AndriyKotsur/quick-booking-api/internal/database"
	"github.com/AndriyKotsur/quick-booking-api/internal/middleware"
)

func RegisterBookingEndpoints(r *chi.Mux, db *database.Queries) {
	bookingController := NewBookingController(db)

	r.Route("/booking", func(r chi.Router) {
		r.Use(middleware.AuthenticateUser)

		r.Post("/", bookingController.CreateBooking)
		r.Get("/", bookingController.GetBookings)
		r.Get("/{id}", bookingController.GetBookingByID)
		r.Put("/{id}", bookingController.UpdateBooking)
		r.Delete("/{id}", bookingController.DeleteBooking)
	})
}
