package endpoints

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/hashicorp/go-multierror"
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

		var multiErr *multierror.Error

		added, err := srv.service.SaveEndpoints(ctx, models.ToServiceEndpoints(endpoints))
		if err != nil {
			errs, ok := err.(*multierror.Error)
			if ok {
				multiErr = errs
				log.Println(ok, multiErr == nil, multiErr.Errors)
			} else {
				switch err {
				default:
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
				return
			}
		}
		var errors []string

		if multiErr != nil {
			for _, err := range multiErr.Errors {
				errors = append(errors, err.Error())
			}
		}

		response := struct {
			AddedEndpoints models.Endpoints `json:"added_endpoints"`
			Errors         []string         `json:"errors,omitempty"`
		}{
			AddedEndpoints: models.FromServiceEndpoints(added),
			Errors:         errors,
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
