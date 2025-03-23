package authentication

import (
	"context"

	"github.com/go-chi/chi"

	"github.com/AndriyKotsur/quick-shopper-api/internal/database"
)

func RegisterAuthEndpoints(r *chi.Mux, db *database.Queries, ctx context.Context) {
	authController := NewAuthController(db, ctx)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
		r.Post("/refresh", authController.Refresh)
		r.Post("/logout", authController.Logout)
	})
}
