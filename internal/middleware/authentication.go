package middleware

import (
	"net/http"
	"strings"

	"github.com/AndriyKotsur/quick-shopper-api/internal/utility/token"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := token.ValidateToken(tokenString, viper.GetString("access_token.public_key"))
		if err != nil {
			log.Error().Msgf("Error validating token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(r.Context())

		next.ServeHTTP(w, r)
	})
}
