package controllers

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, err error) error {
	return writeJSON(w, status, errorResponse{Error: err.Error()})
}
