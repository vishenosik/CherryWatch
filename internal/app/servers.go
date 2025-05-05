package app

import (
	// std
	"log/slog"

	// pkg
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/vishenosik/web/config"
	logger "github.com/vishenosik/web/log"
	middleW "github.com/vishenosik/web/middleware"

	// internal

	"github.com/vishenosik/CherryWatch/pkg/http"
	middlewarepkg "github.com/vishenosik/CherryWatch/pkg/http/middleware"
	"github.com/vishenosik/CherryWatch/pkg/versions"
	// _ "github.com/vishenosik/CherryWatch/internal/gen/swagger"
)

type Service interface {
	Routers() *chi.Mux
}

func newHttpServer(conf Config, log *slog.Logger, services ...Service) Server {

	router := chi.NewRouter()
	router.Use(
		middlewarepkg.ApiVersionMiddleware(versions.MustInitApiVersion("1.0")),
		middleW.RequestLogger(log),
	)

	router.Get("/swagger/*", httpSwagger.Handler())

	for i := range services {
		router.Mount("/api", services[i].Routers())
	}

	return http.NewHttpApp(
		http.Config{
			Server: config.Server{
				Port: conf.RestConfig.Port,
			},
		},
		log.With(logger.AppComponent("http")),
		router,
	)
}
