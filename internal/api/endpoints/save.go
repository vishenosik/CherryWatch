package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	apimodels "github.com/vishenosik/CherryWatch/internal/api/models"
	"github.com/vishenosik/CherryWatch/pkg/httpjson"
	dev "github.com/vishenosik/web-tools/log"
)

func (srv server) saveEndpoint() http.HandlerFunc {

	const op = "api.endpoints.saveEndpoint"

	log := srv.log.With(
		dev.Operation(op),
	)

	return func(w http.ResponseWriter, r *http.Request) {

		endpoints, err := httpjson.Decode[apimodels.Endpoints](r)
		if err != nil {
			http.Error(w, "failed to decode request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		req := apimodels.ToServiceEndpoints(endpoints)

		err = srv.service.SaveEndpoint(ctx, req)

		if err != nil {
			log.Error("failed to save endpoints", dev.Error(err))

			switch {
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		response := apimodels.FromServiceEndpoints(req)

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
