package dto

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status    bool        `json:"status"`
	Message   string      `json:"message,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func RespondWithSuccess(w http.ResponseWriter, status int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  true,
		Message: msg,
		Data:    data,
	})
}

func RespondWithError(w http.ResponseWriter, err error) {
	var appErr *AppError
	
	if e, ok := err.(*AppError); ok {
		appErr = e
	} else {
		appErr = &AppError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			ErrorCode:  "INTERNAL_ERROR",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Status:    false,
		Message:   appErr.Message,
		ErrorCode: appErr.ErrorCode,
	})
}

