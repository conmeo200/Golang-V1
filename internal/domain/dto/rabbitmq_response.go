package dto

import (
	"encoding/json"
	"net/http"
)

type RabbitMQResponse struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type OrderMessage struct {
	OrderUUID string  `json:"uuid"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
}

func SendOrderMessage(w http.ResponseWriter, event string, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(RabbitMQResponse{
		Event: event,
		Data:  data,
	})
}
