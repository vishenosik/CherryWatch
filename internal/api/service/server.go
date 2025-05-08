package service

import (
	"fmt"

	"github.com/go-chi/chi/v5"
)

type serviceAPI struct {
}

type server = *serviceAPI

func NewHttpServer() *serviceAPI {

	return &serviceAPI{}

}

func (srv server) Routers(r chi.Router) {
	r.Get(mount("ping"), srv.ping())
}

func mount(method string) string {
	return fmt.Sprintf("/service.%s", method)
}
