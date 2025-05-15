package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vishenosik/CherryWatch/internal/api/models"
	srvmodels "github.com/vishenosik/CherryWatch/internal/services/models"
	pkghttp "github.com/vishenosik/web/http"
	"github.com/vishenosik/web/multierr"
)

type SaveResponse struct {
	AddedEndpoints models.Endpoints `json:"added_endpoints,omitempty"`
	pkghttp.ErrorResponse
}

func (srv server) save_1_0() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		endpoints, err := pkghttp.Decode[models.Endpoints](r)
		if err != nil {
			pkghttp.SendErrors(w, http.StatusBadRequest, "failed to decode request body")
			return
		}

		if len(endpoints) == 0 {
			pkghttp.SendErrors(w, http.StatusNoContent, "nothing to add")
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		errs := new(multierr.Error)

		added, err := srv.service.SaveEndpoints(ctx, models.ToServiceEndpoints(endpoints))
		if err != nil {
			errs.Append(err)
		}

		// Prepare response

		response := SaveResponse{
			AddedEndpoints: models.FromServiceEndpoints(added),
		}

		writeErrorResponse := func(statusCode int) {
			w.WriteHeader(statusCode)
			response.ErrorResponse = pkghttp.NewErrorResponse(statusCode, errs.CriticalString(), errs.List()...)
		}

		handleErr := func(err error) {
			if err != nil {
				if errors.Is(err, srvmodels.ErrContentNotAdded) {
					writeErrorResponse(http.StatusNotAcceptable)
					return
				}
				writeErrorResponse(http.StatusPartialContent)
			}
		}

		handleErr(errs.ErrorOrNil())

		if err := json.NewEncoder(w).Encode(response); err != nil {
			pkghttp.SendErrors(w, http.StatusInternalServerError, "failed to encode response")
			return
		}
	}
}
