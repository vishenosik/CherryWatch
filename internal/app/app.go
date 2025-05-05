package app

import (
	"context"
	"fmt"
	"log/slog"

	grpcApp "github.com/vishenosik/CherryWatch/internal/app/grpc"
	restApp "github.com/vishenosik/CherryWatch/internal/app/rest"
	"github.com/vishenosik/CherryWatch/internal/services/endpoints"
	"github.com/vishenosik/CherryWatch/internal/store/sql/sqlite"

	appctx "github.com/vishenosik/CherryWatch/internal/app/context"
	"github.com/vishenosik/web/config"
)

type App struct {
	log     *slog.Logger
	servers []Server
	pool    *Pool
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
	sqliteStore := sqlite.MustInitSqlite(appContext.Config.StorePath)

	endpointsService := endpoints.NewService(log, endpoints.Config{}, sqliteStore)

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
		endpointsService,
	)

	app := newApp(log, grpcServer, restServer)

	app.pool = MustNewPool(endpointsService.TasksChan())

	return app, nil
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

	app.pool.Start(context.TODO())
}

func (app *App) Stop(ctx context.Context) {

	const msg = "app stopping"

	signal, ok := appctx.SignalCtx(ctx)
	if ok {
		app.log.Info(msg, slog.String("signal", signal.Signal.String()))
	} else {
		app.log.Info(msg)
	}

	app.pool.Stop()

	for _, server := range app.servers {
		server.Stop(ctx)
	}

	app.log.Info("app stopped")
}
