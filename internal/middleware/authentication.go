package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/AndriyKotsur/quick-booking-api/internal/utility/token"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const userKey = "userId"

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userClaims, err := token.ValidateToken(tokenString, viper.GetString("access_token.public_key"))
		if err != nil {
			log.Error().Msgf("Error validating token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, uuid.MustParse(fmt.Sprint(userClaims)))
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
