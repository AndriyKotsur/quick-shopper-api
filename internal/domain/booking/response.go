package booking

import (
	"time"

	"github.com/google/uuid"
)

type BookingResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	VenueID   uuid.UUID `json:"venue_id,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
	Status    string    `json:"status,omitempty"`
}
