package endpoints

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/vishenosik/CherryWatch/internal/services/endpoints/models"
	"github.com/vishenosik/web-tools/api"
)

type Endpoints interface {
	SaveEndpoint(
		ctx context.Context,
		endpoint *models.Endpoint,
	) (endpointID string, err error)
}

type authenticationAPI struct {
	log     *slog.Logger
	service Endpoints
}

type server = *authenticationAPI

func NewAuthenticationServer(
	log *slog.Logger,
	service Endpoints,
) *authenticationAPI {

	return &authenticationAPI{
		log:     log,
		service: service,
	}

}

func (srv server) Routers() *chi.Mux {

	// Creating a New Router
	endpointsRouter := chi.NewMux()
	endpointsRouter.Post("", srv.saveEndpoint())

	router := chi.NewMux()
	router.Mount(api.ApiV1("/endpoints"), endpointsRouter)

	return router
}
