package restApp

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	appctx "github.com/vishenosik/CherryWatch/internal/app/context"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	// _ "github.com/vishenosik/CherryWatch/internal/gen/swagger"
	"github.com/vishenosik/web-tools/config"
	middleW "github.com/vishenosik/web-tools/middleware"
)

type App struct {
	log    *slog.Logger
	server *http.Server
	port   uint16
}

type Config struct {
	Server config.Server
}

func NewRestApp(
	ctx context.Context,
	config Config,
	// authenticationService authentication.Authentication,
) *App {

	err := config.Server.Validate()
	if err != nil {
		panic(errors.Wrap(err, "failed to validate REST config"))
	}

	app, err := newRestApp(ctx, config)
	if err != nil {
		panic(err)
	}
	return app
}

func newRestApp(
	ctx context.Context,
	config Config,
	// authenticationService authentication.Authentication,
) (*App, error) {

	appContext := appctx.AppCtx(ctx)

	log := appContext.Logger

	// authentication := authentication.NewAuthenticationServer(log, authenticationService)

	router := chi.NewRouter()
	router.Use(
		middleW.RequestLogger(log),
	)

	router.Get("/swagger/*", httpSwagger.Handler())

	setRouters(
		router,
		// authentication,
	)

	return &App{
		log: log,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Server.Port),
			Handler: router,
		},
		port: config.Server.Port,
	}, nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {

	const op = "restApp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Any("port", a.port),
	)

	log.Info("REST server is running", slog.String("addr", a.server.Addr))

	if err := a.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, op)
		}
	}

	return nil
}

func (a *App) Stop(ctx context.Context) {

	const op = "restApp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping Rest server", slog.Any("port", a.port))

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}

type Service interface {
	InitRouters() *chi.Mux
}

func setRouters(router *chi.Mux, services ...Service) {
	for i := range services {
		router.Mount("/api", services[i].InitRouters())
	}
}
