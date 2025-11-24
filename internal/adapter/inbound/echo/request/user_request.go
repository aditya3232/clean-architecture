package request

type SignInRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"min=8,required"`
}

type SignUpRequest struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"email,required"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"email,required"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"password,omitempty"`
	NewPassword     string `json:"password_new" validate:"required"`
	ConfirmPassword string `json:"password_confirmation" validate:"required"`
}

type UpdateDataUserRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email" validate:"omitempty,email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Lat     string `json:"lat"`
	Lng     string `json:"lng"`
	Photo   string `json:"photo"`
}
