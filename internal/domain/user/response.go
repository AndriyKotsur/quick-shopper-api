package user

import "github.com/google/uuid"

type UserResponse struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Name  string    `json:"name,omitempty"`
	Email string    `json:"email,omitempty"`
	Role  string    `json:"role,omitempty"`
}
