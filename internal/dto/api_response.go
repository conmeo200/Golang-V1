package dto

import (
	"log"
)

type APIResponse struct {
	Status  bool      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendSuccess(data interface{}) APIResponse {

	log.Println(data)
	return APIResponse{
		Status: true,
		Data: data	,
	}
}

func SendError(msg string) APIResponse {
	return APIResponse{
		Status: false,
		Message: msg,
	}
}

