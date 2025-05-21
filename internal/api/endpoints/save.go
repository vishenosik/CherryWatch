package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vishenosik/CherryWatch/internal/api/models"
	srvmodels "github.com/vishenosik/CherryWatch/internal/services/models"
	_errors "github.com/vishenosik/gocherry/pkg/errors"
	_http "github.com/vishenosik/gocherry/pkg/http"
)

type SaveResponse struct {
	AddedEndpoints models.Endpoints `json:"added_endpoints,omitempty"`
	_http.ErrorResponse
}

func (srv server) save() (string, func(chi.Router)) {

	versionMiddleware, versionHandler := _http.DotVersionMiddlewareHandler("1.0")

	return route("save"), func(r chi.Router) {
		r.Use(
			versionMiddleware,
		)
		r.Post(_http.BlankRoute, versionHandler(_http.HandlersMap{
			"1.0": srv.save_1_0(),
		}))
	}
}

func (srv server) save_1_0() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		endpoints, err := _http.Decode[models.Endpoints](r)
		if err != nil {
			_http.SendErrors(w, http.StatusBadRequest, "failed to decode request body")
			return
		}

		if len(endpoints) == 0 {
			_http.SendErrors(w, http.StatusNoContent, "nothing to add")
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		errs := new(_errors.Error)

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
			response.ErrorResponse = _http.NewErrorResponse(statusCode, errs.CriticalString(), errs.List()...)
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
			_http.SendErrors(w, http.StatusInternalServerError, "failed to encode response")
			return
		}
	}
}
