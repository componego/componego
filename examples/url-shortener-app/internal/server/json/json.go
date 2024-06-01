package json

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func Get[T any](request *http.Request) (value T, err error) {
	err = json.NewDecoder(request.Body).Decode(&value)
	return value, err
}

func Send(response http.ResponseWriter, data any, err error) {
	jsonResponse := &jsonResponse{
		Status: err == nil,
		Data:   data,
	}
	if !jsonResponse.Status {
		jsonResponse.Error = err.Error()
	}
	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(jsonResponse); err != nil {
		http.Error(response, "error sending response", http.StatusBadRequest)
	}
}
