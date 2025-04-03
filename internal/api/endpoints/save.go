package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vishenosik/CherryWatch/internal/api/models"
	"github.com/vishenosik/CherryWatch/pkg/httpjson"
)

func (srv server) saveEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		endpoints, err := httpjson.Decode[models.Endpoints](r)
		if err != nil {
			http.Error(w, "failed to decode request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		added, err := srv.service.SaveEndpoints(ctx, models.ToServiceEndpoints(endpoints))
		if err != nil {
			switch {
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		response := models.FromServiceEndpoints(added)

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
