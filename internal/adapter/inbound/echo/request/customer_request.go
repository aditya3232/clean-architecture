package request

type CustomerRequest struct {
	Name                 string  `json:"name" validate:"required"`
	Email                string  `json:"email" validate:"required,email,uniqueEmail"`
	Password             string  `json:"password" validate:"required,min=8"`
	PasswordConfirmation string  `json:"password_confirmation" validate:"required,min=8"`
	Phone                string  `json:"phone" validate:"required,number"`
	Address              string  `json:"address"`
	Lat                  float64 `json:"lat"`
	Lng                  float64 `json:"lng"`
	Photo                string  `json:"photo"`
	RoleID               int64   `json:"role_id" validate:"required"`
}

type UpdateCustomerRequest struct {
	Name    string  `json:"name"`
	Email   string  `json:"email" validate:"omitempty,email,uniqueEmail"`
	Phone   string  `json:"phone" validate:"number"`
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Photo   string  `json:"photo"`
}
