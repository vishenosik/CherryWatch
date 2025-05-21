package app

import (
	"context"

	"github.com/vishenosik/gocherry"
	_http "github.com/vishenosik/gocherry/pkg/http"

	"log/slog"

	endpointsApi "github.com/vishenosik/CherryWatch/internal/api/endpoints"
	"github.com/vishenosik/CherryWatch/internal/api/service"
	"github.com/vishenosik/CherryWatch/internal/services/endpoints"
	"github.com/vishenosik/CherryWatch/internal/store/sql/sqlite"
	"github.com/vishenosik/web/logs"

	"net/http"

	// pkg
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type App struct {
	Server
}

func NewApp() (*App, error) {
	app, err := gocherry.NewApp()
	if err != nil {
		panic(err)
	}

	log := app.Log

	conf := mustLoadEnvConfig()
	log.Debug("config loaded from env", slog.Any("config", conf))

	// Stores init
	sqliteStore := sqlite.MustInitSqlite(conf.StorePath, log)

	// Usecases init
	endpointsService := endpoints.NewService(log, endpoints.Config{}, sqliteStore)

	handler := _http.NewHttpServer(
		log.With(logs.AppComponent("http")),
		NewHttpServer(log,
			endpointsApi.NewHttpServer(endpointsService),
			service.NewHttpServer(),
		),
	)

	pool, err := gocherry.NewPool(
		log.With(logs.AppComponent("worker pool")),
		endpointsService.TasksChan(),
	)
	if err != nil {
		return nil, err
	}

	app.AddServices(handler, pool)

	return &App{
		Server: app,
	}, nil

}

type Service interface {
	Routers(r chi.Router)
}

func NewHttpServer(log *slog.Logger, services ...Service) http.Handler {
	log_ := log.With(logs.AppComponent("http"))

	router := chi.NewRouter()
	router.Use(
		_http.RequestLogger(log_),
		// http.ApiVersionMiddleware(versions.DotVersion{}, "2.0"),
	)

	router.Get("/swagger/*", httpSwagger.Handler())

	router.Route("/api", func(r chi.Router) {
		for i := range services {
			services[i].Routers(r)
		}
	})
	return router
}
