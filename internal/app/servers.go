package app

import (
	// std
	"log/slog"

	// pkg
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/vishenosik/web/config"
	logger "github.com/vishenosik/web/log"

	// internal

	"github.com/vishenosik/CherryWatch/pkg/http"
	// _ "github.com/vishenosik/CherryWatch/internal/gen/swagger"
)

type Service interface {
	Routers(r chi.Router)
}

func newHttpServer(conf Config, log *slog.Logger, services ...Service) Server {

	log_ := log.With(logger.AppComponent("http"))

	router := chi.NewRouter()
	router.Use(
		http.RequestLogger(log_),
		// http.ApiVersionMiddleware(versions.DoubleVersion{}, "2.0"),
	)

	router.Get("/swagger/*", httpSwagger.Handler())

	router.Route("/api", func(r chi.Router) {
		for i := range services {
			services[i].Routers(r)
		}
	})

	return http.NewHttpApp(
		http.Config{
			Server: config.Server{
				Port: conf.RestConfig.Port,
			},
		},
		log_,
		router,
	)
}
