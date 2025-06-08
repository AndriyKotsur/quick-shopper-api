package venue

type CreateVenueRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=20"`
	Location string `json:"location" validate:"required, min=2,max=50"`
	Type     string `json:"type" validate:"required"`
}

type UpdateVenueRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=20"`
	Location string `json:"location" validate:"required, min=2,max=50"`
}
