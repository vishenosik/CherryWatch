package service

import (
	"github.com/go-chi/chi/v5"
	http_pkg "github.com/vishenosik/CherryWatch/pkg/http"
	"github.com/vishenosik/CherryWatch/pkg/versions"
)

var mount = http_pkg.MethodFunc("service")

type serviceAPI struct {
}

type server = *serviceAPI

func NewHttpServer() *serviceAPI {
	return &serviceAPI{}
}

func (srv server) Routers(r chi.Router) {
	r.Route(mount("ping"), func(r chi.Router) {

		r.Use(
			http_pkg.ApiVersionMiddleware(versions.DoubleVersion{}, "1.1"),
		)

		r.Get(http_pkg.BlankRoute, http_pkg.VersionedHandler(
			http_pkg.VersionedHandlersMap{
				"1.0": srv.ping_1_0(),
				"1.1": srv.ping_1_1(),
			},
		))
	})
}
