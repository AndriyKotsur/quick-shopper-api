package authentication

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/AndriyKotsur/quick-shopper-api/internal/database"
	"github.com/AndriyKotsur/quick-shopper-api/internal/utility/json"
	"github.com/AndriyKotsur/quick-shopper-api/internal/utility/password"
	"github.com/AndriyKotsur/quick-shopper-api/internal/utility/respond"
	"github.com/AndriyKotsur/quick-shopper-api/internal/utility/token"
	"github.com/AndriyKotsur/quick-shopper-api/internal/utility/validate"
)

type AuthController struct {
	db  *database.Queries
	ctx context.Context
}

func NewAuthController(db *database.Queries, ctx context.Context) *AuthController {
	return &AuthController{db, ctx}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var params *RegisterRequest

	err := json.DecodeJSON(w, r, &params)
	if err != nil {
		log.Error().Msgf("Error decoding JSON: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	errs := validate.Validate(params)
	if errs != nil {
		log.Error().Msgf("Error validating parameters: %v", errs)
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	hashedPassword, err := password.HashPassword(params.Password)
	if err != nil {
		log.Error().Msgf("Error encrypting password: %v", err)
		respond.Error(w, http.StatusInternalServerError, nil)
		return
	}

	userArgs := &database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Email:     params.Email,
		Role:      params.Role,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = c.db.CreateUser(c.ctx, *userArgs)
	if err != nil {
		log.Error().Msgf("Error registering user: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	respond.Status(w, http.StatusCreated)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var params *LoginRequest

	err := json.DecodeJSON(w, r, &params)
	if err != nil {
		log.Error().Msgf("Error decoding JSON: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	user, err := c.db.GetUserByEmail(c.ctx, params.Email)
	if err != nil {
		log.Error().Msgf("Invalid email or password: %v", err)
		respond.Error(w, http.StatusBadRequest, errors.New("Invalid email or password"))
		return
	}

	if err := password.ComparePassword(user.Password, params.Password); err != nil {
		log.Error().Msgf("Invalid email or password: %v", err)
		respond.Error(w, http.StatusBadRequest, errors.New("Invalid email or password"))
		return
	}

	accessToken, err := token.GenerateToken(user.ID, viper.GetString("access_token.private_key"), viper.GetDuration("access_token.expires_in"))
	if err != nil {
		log.Error().Msgf("Error generating token: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	refreshToken, err := token.GenerateToken(user.ID, viper.GetString("refresh_token.private_key"), viper.GetDuration("refresh_token.expires_in"))
	if err != nil {
		log.Error().Msgf("Error generating token: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   viper.GetInt("refresh_token.max_age") * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	respond.JSON(w, http.StatusOK, accessToken)
}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	claims, err := token.ValidateToken(cookie.Value, viper.GetString("access_token.public_key"))
	if err != nil {
		log.Error().Msgf("Error validating token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := token.GenerateToken(claims, viper.GetString("access_token.private_key"), viper.GetDuration("access_token.expires_in"))
	if err != nil {
		log.Error().Msgf("Error generating token: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	newRefreshToken, err := token.GenerateToken(claims, viper.GetString("refresh_token.private_key"), viper.GetDuration("refresh_token.expires_in"))
	if err != nil {
		log.Error().Msgf("Error generating token: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	respond.JSON(w, http.StatusOK, newAccessToken)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	})

	respond.JSON(w, http.StatusOK, nil)
}
