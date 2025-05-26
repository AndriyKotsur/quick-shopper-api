package user

import (
	"github.com/go-chi/chi"

	"github.com/AndriyKotsur/quick-booking-api/internal/database"
	"github.com/AndriyKotsur/quick-booking-api/internal/middleware"
)

func RegisterUserEndpoints(r *chi.Mux, db *database.Queries) {
	userController := NewUserController(db)

	r.Route("/user", func(r chi.Router) {
		r.Use(middleware.AuthenticateUser)

		r.Post("/", userController.CreateUser)
		r.Get("/me", userController.GetMe)
		r.Get("/{id}", userController.GetUserById)
		r.Put("/{id}", userController.UpdateUser)
		r.Delete("/{id}", userController.DeleteUser)
	})
}
