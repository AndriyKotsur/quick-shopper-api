package venue

import "github.com/google/uuid"

type VenueResponse struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Location string    `json:"location,omitempty"`
	Type     string    `json:"type,omitempty"`
}
