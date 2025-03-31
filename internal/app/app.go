package app

import (
	"context"
	"fmt"
	"log/slog"

	grpcApp "github.com/vishenosik/CherryWatch/internal/app/grpc"
	restApp "github.com/vishenosik/CherryWatch/internal/app/rest"

	appctx "github.com/vishenosik/CherryWatch/internal/app/context"
	"github.com/vishenosik/web-tools/config"
)

type App struct {
	log     *slog.Logger
	servers []Server
}

type Server interface {
	MustRun()
	Stop(ctx context.Context)
}

func MustInitApp() *App {
	app, err := NewApp()
	if err != nil {
		panic(fmt.Sprintf("failed to create app %s", err))
	}
	return app
}

func NewApp() (*App, error) {

	ctx := appctx.SetupAppCtx()
	appContext := appctx.AppCtx(ctx)

	log := appContext.Logger
	conf := appContext.Config

	// Stores init
	_, err := loadSqlStore(ctx)
	if err != nil {
		return nil, err
	}

	grpcServer := grpcApp.NewGrpcApp(
		log,
		grpcApp.Config{
			Server: config.Server{
				Port: conf.GrpcConfig.Port,
			},
		},
		// authenticationService,
	)

	restServer := restApp.NewRestApp(
		ctx,
		restApp.Config{
			Server: config.Server{
				Port: conf.RestConfig.Port,
			},
		},
		// authenticationService,
	)

	return newApp(log, grpcServer, restServer), nil
}

func newApp(
	logger *slog.Logger,
	apps ...Server,
) *App {
	return &App{
		log:     logger,
		servers: apps,
	}
}

func (app *App) MustRun() {

	app.log.Info("start app")

	for _, server := range app.servers {
		go server.MustRun()
	}
}

func (app *App) Stop(ctx context.Context) {

	const msg = "app stopping"

	signal, ok := appctx.SignalCtx(ctx)
	if !ok {
		app.log.Info(msg, slog.String("signal", signal.Signal.String()))
	} else {
		app.log.Info(msg)
	}

	for _, server := range app.servers {
		server.Stop(ctx)
	}

	app.log.Info("app stopped")
}
