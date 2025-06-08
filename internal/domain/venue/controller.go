package venue

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/AndriyKotsur/quick-booking-api/internal/database"
	"github.com/AndriyKotsur/quick-booking-api/internal/utility/json"
	"github.com/AndriyKotsur/quick-booking-api/internal/utility/respond"
	"github.com/AndriyKotsur/quick-booking-api/internal/utility/validate"
)

type VenueController struct {
	db *database.Queries
}

func NewVenueController(db *database.Queries) *VenueController {
	return &VenueController{db}
}

func (c *VenueController) GetVenues(w http.ResponseWriter, r *http.Request) {
	venues, err := c.db.GetVenues(r.Context())
	if err != nil {
		log.Error().Msgf("Venues not found: %v", err)
		respond.Error(w, http.StatusNotFound, nil)
		return
	}

	venueList := []VenueResponse{}

	for _, venue := range venues {
		venueList = append(venueList, VenueResponse{
			ID:       venue.ID,
			Name:     venue.Name,
			Location: venue.Location,
			Type:     venue.Type,
		})
	}

	respond.JSON(w, http.StatusOK, venueList)
}

func (c *VenueController) GetVenueByID(w http.ResponseWriter, r *http.Request) {
	urlParams := chi.URLParam(r, "id")

	venueId, err := uuid.Parse(urlParams)
	if err != nil {
		log.Error().Msgf("Error parsing URL params: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	venue, err := c.db.GetVenueByID(r.Context(), venueId)
	if err != nil {
		log.Error().Msgf("Venue not found: %v", err)
		respond.Error(w, http.StatusNotFound, nil)
		return
	}

	respond.JSON(w, http.StatusOK, VenueResponse{
		ID:       venue.ID,
		Name:     venue.Name,
		Location: venue.Location,
		Type:     venue.Type,
	})
}

func (c *VenueController) CreateVenue(w http.ResponseWriter, r *http.Request) {
	var bodyParams *CreateVenueRequest

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

	userId, ok := r.Context().Value("userId").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	venueArgs := &database.CreateVenueParams{
		ID:        uuid.New(),
		UserID:    userId,
		Name:      bodyParams.Name,
		Location:  bodyParams.Location,
		Type:      bodyParams.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newVenue, err := c.db.CreateVenue(r.Context(), *venueArgs)
	if err != nil {
		log.Error().Msgf("Error creating venue: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	respond.JSON(w, http.StatusCreated, VenueResponse{
		ID:       newVenue.ID,
		Name:     newVenue.Name,
		Location: newVenue.Location,
		Type:     newVenue.Type,
	})
}

func (c *VenueController) UpdateVenue(w http.ResponseWriter, r *http.Request) {
	urlParams := chi.URLParam(r, "id")

	venueId, err := uuid.Parse(urlParams)
	if err != nil {
		log.Error().Msgf("Error parsing URL params: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	var bodyParams *UpdateVenueRequest

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

	updatedVenue, err := c.db.UpdateVenue(r.Context(), database.UpdateVenueParams{
		ID:        venueId,
		Name:      bodyParams.Name,
		Location:  bodyParams.Location,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Error().Msgf("Error updating venue: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	respond.JSON(w, http.StatusOK, VenueResponse{
		ID:       updatedVenue.ID,
		Name:     updatedVenue.Name,
		Location: updatedVenue.Location,
		Type:     updatedVenue.Type,
	})
}

func (c *VenueController) DeleteVenue(w http.ResponseWriter, r *http.Request) {
	urlParams := chi.URLParam(r, "id")

	venueId, err := uuid.Parse(urlParams)
	if err != nil {
		log.Error().Msgf("Error parsing URL params: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	err = c.db.DeleteVenue(r.Context(), venueId)
	if err != nil {
		log.Error().Msgf("Error deleting venue: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	respond.Status(w, http.StatusOK)
}
