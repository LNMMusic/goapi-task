package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Ok(w http.ResponseWriter, code int, msg string, data any) {
	// set status code
	w.WriteHeader(code)

	// set headers
	w.Header().Set("Content-Type", "application/json")

	// set body
	resp := Response{
		Message: msg,
		Data:    data,
	}
	json.NewEncoder(w).Encode(resp)
}

func Err(w http.ResponseWriter, code int, msg string) {
	// set status code
	w.WriteHeader(code)

	// set headers
	w.Header().Set("Content-Type", "application/json")

	// set body
	resp := Response{
		Message: msg,
	}
	json.NewEncoder(w).Encode(resp)
}