package dto

import "net/http"

type AppError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	ErrorCode  string `json:"error_code,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, msg string, code string) *AppError {
	return &AppError{
		StatusCode: status,
		Message:    msg,
		ErrorCode:  code,
	}
}

// Các lỗi thường gặp
var (
	ErrInvalidRequest = &AppError{StatusCode: http.StatusBadRequest, Message: "invalid request format", ErrorCode: "INVALID_REQUEST"}
	ErrUnauthorized   = &AppError{StatusCode: http.StatusUnauthorized, Message: "unauthorized access", ErrorCode: "UNAUTHORIZED"}
	ErrInternal       = &AppError{StatusCode: http.StatusInternalServerError, Message: "internal server error", ErrorCode: "INTERNAL_ERROR"}
)
