package endpoints

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/vishenosik/CherryWatch/internal/services/models"
)

type Endpoints interface {
	SaveEndpoints(
		ctx context.Context,
		endpoints models.Endpoints,
	) (added models.Endpoints, err error)
}

type endpointsAPI struct {
	log     *slog.Logger
	service Endpoints
}

type server = *endpointsAPI

func NewAuthenticationServer(
	log *slog.Logger,
	service Endpoints,
) *endpointsAPI {

	return &endpointsAPI{
		log:     log,
		service: service,
	}

}

func (srv server) Routers() *chi.Mux {

	router := chi.NewMux()
	router.Post(mount("save"), srv.saveEndpoint())

	return router
}

func mount(method string) string {
	return fmt.Sprintf("/endpoints.%s", method)
}
