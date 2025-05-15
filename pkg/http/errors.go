package http

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an http API error
type ErrorResponse struct {
	Message  string   `json:"message,omitempty"`
	Critical string   `json:"critical,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

func NewErrorResponse(statusCode int, critical string, errors ...string) ErrorResponse {
	return ErrorResponse{
		Message:  http.StatusText(statusCode),
		Critical: critical,
		Errors:   errors,
	}
}

// sendError sends a JSON error response
func SendErrors(w http.ResponseWriter, statusCode int, critical string, errors ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(NewErrorResponse(statusCode, critical, errors...))
}
