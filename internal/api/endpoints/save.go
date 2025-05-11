package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vishenosik/CherryWatch/internal/api/models"
	pkghttp "github.com/vishenosik/CherryWatch/pkg/http"
	"github.com/vishenosik/CherryWatch/pkg/multierr"
)

func (srv server) save_1_0() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		endpoints, err := pkghttp.Decode[models.Endpoints](r)
		if err != nil {
			pkghttp.SendErrors(w, http.StatusBadRequest, "failed to decode request body")
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		errs := new(multierr.Error)

		added, err := srv.service.SaveEndpoints(ctx, models.ToServiceEndpoints(endpoints))
		if err != nil {
			errs.Append(err)
		}

		response := struct {
			AddedEndpoints models.Endpoints `json:"added_endpoints,omitempty"`
			Errors         []string         `json:"errors,omitempty"`
		}{
			AddedEndpoints: models.FromServiceEndpoints(added),
			Errors:         errs.List(),
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
