package dto

import (
	"github.com/conmeo200/Golang-V1/internal/model"
)

type AuthResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Balance   float64 `json:"balance"`
	Role      string  `json:"role"`
	Status    string  `json:"status"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
	DeletedAt int64   `json:"deleted_at"`
}

func ToAuthResponse(user *model.User) UserResponse {
	return UserResponse{
		ID: user.ID.String(),
		Email: user.Email,
		Balance: user.Balance,
		Role: user.Role,
		Status: user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}


