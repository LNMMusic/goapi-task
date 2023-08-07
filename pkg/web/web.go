package web

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, body any) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}