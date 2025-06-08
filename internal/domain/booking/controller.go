package booking

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

type BookingController struct {
	db *database.Queries
}

func NewBookingController(db *database.Queries) *BookingController {
	return &BookingController{db}
}

func (c *BookingController) GetBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := c.db.GetBookings(r.Context())
	if err != nil {
		log.Error().Msgf("Bookings not found: %v", err)
		respond.Error(w, http.StatusNotFound, nil)
		return
	}

	bookingList := []BookingResponse{}

	for _, booking := range bookings {
		bookingList = append(bookingList, BookingResponse{
			ID:        booking.ID,
			VenueID:   booking.VenueID,
			StartTime: booking.StartTime,
			EndTime:   booking.EndTime,
			Status:    booking.Status,
		})
	}

	respond.JSON(w, http.StatusOK, bookingList)
}

func (c *BookingController) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	urlParams := chi.URLParam(r, "id")

	bookingId, err := uuid.Parse(urlParams)
	if err != nil {
		log.Error().Msgf("Error parsing URL params: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	booking, err := c.db.GetBookingByID(r.Context(), bookingId)
	if err != nil {
		log.Error().Msgf("Booking not found: %v", err)
		respond.Error(w, http.StatusNotFound, nil)
		return
	}

	respond.JSON(w, http.StatusOK, BookingResponse{
		ID:        booking.ID,
		VenueID:   booking.VenueID,
		StartTime: booking.StartTime,
		EndTime:   booking.EndTime,
		Status:    booking.Status,
	})
}

func (c *BookingController) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var bodyParams *CreateBookingParams

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

	bookingArgs := &database.CreateBookingParams{
		ID:        uuid.New(),
		UserID:    userId,
		VenueID:   bodyParams.VenueID,
		StartTime: bodyParams.StartTime,
		EndTime:   bodyParams.EndTime,
		Status:    bodyParams.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newBooking, err := c.db.CreateBooking(r.Context(), *bookingArgs)
	if err != nil {
		log.Error().Msgf("Error creating booking: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	respond.JSON(w, http.StatusCreated, BookingResponse{
		ID:        newBooking.ID,
		VenueID:   newBooking.VenueID,
		StartTime: newBooking.StartTime,
		EndTime:   newBooking.EndTime,
		Status:    newBooking.Status,
	})
}

func (c *BookingController) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	urlParams := chi.URLParam(r, "id")

	bookingId, err := uuid.Parse(urlParams)
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

	updatedBooking, err := c.db.UpdateBooking(r.Context(), database.UpdateBookingParams{
		ID:        bookingId,
		Status:    bodyParams.Status,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Error().Msgf("Error updating booking: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	respond.JSON(w, http.StatusOK, BookingResponse{
		ID:        updatedBooking.ID,
		VenueID:   updatedBooking.VenueID,
		StartTime: updatedBooking.StartTime,
		EndTime:   updatedBooking.EndTime,
		Status:    updatedBooking.Status,
	})
}

func (c *BookingController) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	urlParams := chi.URLParam(r, "id")

	bookingId, err := uuid.Parse(urlParams)
	if err != nil {
		log.Error().Msgf("Error parsing URL params: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	err = c.db.DeleteBooking(r.Context(), bookingId)
	if err != nil {
		log.Error().Msgf("Error deleting booking: %v", err)
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	respond.Status(w, http.StatusOK)
}
