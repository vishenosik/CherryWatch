package http

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an http API error
type ErrorResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

// sendError sends a JSON error response
func SendErrors(w http.ResponseWriter, statusCode int, messages ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Message: http.StatusText(statusCode),
		Errors:  messages,
	})
}
