package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vishenosik/CherryWatch/internal/services/models"
	http_pkg "github.com/vishenosik/web/http"
)

var route = http_pkg.MethodFunc("endpoints")

type Endpoints interface {
	SaveEndpoints(
		ctx context.Context,
		endpoints models.Endpoints,
	) (added models.Endpoints, err error)
}

type endpointsAPI struct {
	service Endpoints
}

type server = *endpointsAPI

func NewHttpServer(
	service Endpoints,
) *endpointsAPI {

	return &endpointsAPI{
		service: service,
	}

}

func (srv server) Routers(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(http_pkg.SetHeaders())
		r.Route(srv.save())
		r.Route(route("get"), func(r chi.Router) {
			r.Get(http_pkg.BlankRoute, func(w http.ResponseWriter, r *http.Request) {

				response := struct {
					Message string `json:"message"`
					Status  string `json:"status"`
				}{
					Message: "getter",
					Status:  "endpoints.get ok",
				}

				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "failed to encode response", http.StatusInternalServerError)
					return
				}
			})
		})
	})

}
