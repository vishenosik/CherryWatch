package app

import (
	"context"
	"path"

	"github.com/vishenosik/gocherry"
	_http "github.com/vishenosik/gocherry/pkg/http"
	"github.com/vishenosik/gocherry/pkg/sql"

	"log/slog"

	endpointsApi "github.com/vishenosik/CherryWatch/internal/api/endpoints"
	"github.com/vishenosik/CherryWatch/internal/api/service"
	"github.com/vishenosik/CherryWatch/internal/services/endpoints"
	"github.com/vishenosik/CherryWatch/internal/store/sql/sqlite"
	"github.com/vishenosik/gocherry/pkg/logs"

	"net/http"

	embed "github.com/vishenosik/CherryWatch"

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

	// Stores init

	store, err := sql.NewSqliteStore(
		sql.WithMigration(
			embed.Migrations,
			path.Join(embed.MigrationsPath, "sqlite"),
		),
	)
	if err != nil {
		panic(err)
	}

	db, err := store.Open(context.TODO())
	if err != nil {
		panic(err)
	}

	endpointsStore := sqlite.NewEndpoints(db)

	// Usecases init

	endpointsService := endpoints.NewService(
		log,
		endpoints.Config{},
		endpointsStore,
	)

	// Services init

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

func NewHttpServer(logger *slog.Logger, services ...Service) http.Handler {
	log := logger.With(logs.AppComponent("http"))

	router := chi.NewRouter()
	router.Use(
		_http.RequestLogger(log),
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
