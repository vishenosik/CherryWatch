package endpoints

import (
	"context"
	"fmt"

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
	r.Post(mount("save"), srv.saveEndpoint())
}

func mount(method string) string {
	return fmt.Sprintf("/endpoints.%s", method)
}
