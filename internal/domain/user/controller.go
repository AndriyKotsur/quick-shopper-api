package user

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/AndriyKotsur/quick-shopper-api/internal/database"
	"github.com/AndriyKotsur/quick-shopper-api/internal/utility/json"
	"github.com/AndriyKotsur/quick-shopper-api/internal/utility/password"
	"github.com/AndriyKotsur/quick-shopper-api/internal/utility/respond"
	"github.com/AndriyKotsur/quick-shopper-api/internal/utility/validate"
)

type UserController struct {
	db *database.Queries
}

func NewUserController(db *database.Queries) *UserController {
	return &UserController{db}
}

func (c *UserController) GetMe(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value("userId").(uuid.UUID)
	log.Info().Interface("userClaims", userClaims).Msg("Parsed user claims")
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := c.db.GetUserByID(r.Context(), userClaims)
	if err != nil {
		log.Error().Msgf("User not found: %v", err)
		respond.Error(w, http.StatusNotFound, nil)
		return
	}

	respond.JSON(w, http.StatusOK, UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

func (c *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	urlParams := chi.URLParam(r, "id")

	userId, err := uuid.Parse(urlParams)
	if err != nil {
		log.Error().Msgf("Error parsing URL params: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	user, err := c.db.GetUserByID(r.Context(), userId)
	if err != nil {
		log.Error().Msgf("User not found: %v", err)
		respond.Error(w, http.StatusNotFound, nil)
		return
	}

	respond.JSON(w, http.StatusOK, UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var bodyParams *CreateUserRequest

	err := json.DecodeJSON(w, r, &bodyParams)
	if err != nil {
		log.Error().Msgf("Error decoding JSON: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	errs := validate.Validate(bodyParams)
	if errs != nil {
		log.Error().Msgf("Error validating parameters: %v", errs)
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	hashedPassword, err := password.HashPassword(bodyParams.Password)
	if err != nil {
		log.Error().Msgf("Error encrypting password: %v", err)
		respond.Error(w, http.StatusInternalServerError, nil)
		return
	}

	userArgs := &database.CreateUserParams{
		ID:        uuid.New(),
		Name:      bodyParams.Name,
		Email:     bodyParams.Email,
		Role:      bodyParams.Role,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newUser, err := c.db.CreateUser(r.Context(), *userArgs)
	if err != nil {
		log.Error().Msgf("Error creating user: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	respond.JSON(w, http.StatusCreated, UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
		Name:  newUser.Name,
	})
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	urlParams := chi.URLParam(r, "id")

	userId, err := uuid.Parse(urlParams)
	if err != nil {
		log.Error().Msgf("Error parsing URL params: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	var bodyParams *UpdateUserRequest

	err = json.DecodeJSON(w, r, &bodyParams)
	if err != nil {
		log.Error().Msgf("Error decoding JSON: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	errs := validate.Validate(bodyParams)
	if errs != nil {
		log.Error().Msgf("Error validating parameters: %v", errs)
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	updatedUser, err := c.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:        userId,
		Name:      bodyParams.Name,
		Email:     bodyParams.Email,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Error().Msgf("Error updating user: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	respond.JSON(w, http.StatusOK, UserResponse{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
		Name:  updatedUser.Name,
	})
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	urlParams := chi.URLParam(r, "id")

	userId, err := uuid.Parse(urlParams)
	if err != nil {
		log.Error().Msgf("Error parsing URL params: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	err = c.db.DeleteUser(r.Context(), userId)
	if err != nil {
		log.Error().Msgf("Error deleting user: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	respond.Status(w, http.StatusOK)
}
