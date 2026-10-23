package handler

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	}
}

type errorResponse struct {
	Error string `json:"error"`
	Field string `json:"field,omitempty"`
}

func writeError(w http.ResponseWriter, status int, message string, field string) {
	writeJSON(w, status, errorResponse{
		Error: message,
		Field: field,
	})
}

func writeInternalServerError(w http.ResponseWriter) {
	writeJSON(w, http.StatusInternalServerError, "Internal server error")
}
