package authentication

import (
	"github.com/go-chi/chi"

	"github.com/AndriyKotsur/quick-booking-api/internal/database"
)

func RegisterAuthEndpoints(r *chi.Mux, db *database.Queries) {
	authController := NewAuthController(db)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
		r.Post("/refresh", authController.Refresh)
		r.Post("/logout", authController.Logout)
	})
}
