package service

import (
	"github.com/go-chi/chi/v5"
	_http "github.com/vishenosik/gocherry/pkg/http"
)

var route = _http.MethodFunc("service")

type serviceAPI struct {
}

type server = *serviceAPI

func NewHttpServer() *serviceAPI {
	return &serviceAPI{}
}

func (srv server) Routers(r chi.Router) {
	r.Group(func(r chi.Router) {

		r.Use(_http.SetHeaders())

		r.Route(route("ping"), func(r chi.Router) {

			versionMiddleware, versionHandler := _http.DotVersionMiddlewareHandler(
				"1.1",
				_http.Min("0.1.1"),
			)

			r.Use(
				versionMiddleware,
			)

			r.Get(_http.BlankRoute, versionHandler(_http.HandlersMap{
				"1.0": srv.ping_1_0(),
				"1.1": srv.ping_1_1(),
			}))
		})
	})
}
