package user

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=20"`
	Email string `json:"email" validate:"required,email"`
}
