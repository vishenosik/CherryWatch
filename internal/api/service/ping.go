package service

import (
	"encoding/json"
	"net/http"
)

type response struct {
	ApiVersion string `json:"api_version"`
	Status     string `json:"status"`
}

func (srv server) ping_1_0() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		response := response{
			ApiVersion: "1.0",
			Status:     "ok buddy",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func (srv server) ping_1_1() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		response := response{
			ApiVersion: "1.1",
			Status:     "ok just ok",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
