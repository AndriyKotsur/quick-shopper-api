package booking

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingParams struct {
	VenueID   uuid.UUID `json:"venue_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
	Status    string    `json:"status" validate:"required"`
}

type UpdateVenueRequest struct {
	Status string `json:"type" validate:"required"`
}
