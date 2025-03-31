package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/vishenosik/CherryWatch/internal/services/endpoints/models"
	dev "github.com/vishenosik/web-tools/log"
)

func (srv server) saveEndpoint() http.HandlerFunc {

	const op = "authentication.http.IsAdmin"

	return func(w http.ResponseWriter, r *http.Request) {

		log := srv.log.With(
			dev.Operation(op),
		)

		endpointID, err := srv.service.SaveEndpoint(r.Context(), &models.Endpoint{})

		if err != nil {
			log.Error("failed to check admin status", dev.Error(err))

			switch {
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		response := struct {
			EndpointID string `json:"endpoint_id"`
		}{
			EndpointID: endpointID,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to encode response", dev.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
