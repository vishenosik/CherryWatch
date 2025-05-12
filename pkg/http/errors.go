package http

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an http API error
type ErrorResponse struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

func NewErrorResponse(statusCode int, errors ...string) ErrorResponse {
	return ErrorResponse{
		Message: http.StatusText(statusCode),
		Errors:  errors,
	}
}

// sendError sends a JSON error response
func SendErrors(w http.ResponseWriter, statusCode int, errors ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(NewErrorResponse(statusCode, errors...))
}
