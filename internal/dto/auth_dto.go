package dto

// LoginRequest defines the data structure for a user login request.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse defines the data structure for a user response.
// It omits sensitive information like the password.
type UserResponse struct {
	ID      string  `json:"id"`
	Email   string  `json:"email"`
	Balance float64 `json:"balance"`
	Role    string  `json:"role"`
	Status  string  `json:"status"`
}
