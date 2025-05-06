package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	endpointsApi "github.com/vishenosik/CherryWatch/internal/api/endpoints"
	grpcApp "github.com/vishenosik/CherryWatch/internal/app/grpc"
	"github.com/vishenosik/CherryWatch/internal/services/endpoints"
	"github.com/vishenosik/CherryWatch/internal/store/sql/sqlite"

	appctx "github.com/vishenosik/CherryWatch/internal/app/context"
	"github.com/vishenosik/web/colors"
	"github.com/vishenosik/web/config"
	logger "github.com/vishenosik/web/log"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
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

	conf := mustLoadEnvConfig()
	log := setupLogger(conf.Env)

	log.Debug("config loaded from env", slog.Any("config", conf))

	// Stores init
	sqliteStore := sqlite.MustInitSqlite(conf.StorePath, log)

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

	httpServer := newHttpServer(
		conf, log,
		endpointsApi.NewHttpServer(endpointsService),
	)

	app := newApp(log, grpcServer, httpServer)

	app.pool = MustInitPool(endpointsService.TasksChan())

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

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler
	switch env {

	case EnvProd:
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)

	case EnvTest:
		handler = slog.NewJSONHandler(
			io.Discard,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)

	case EnvDev:
		handler = logger.NewHandler(
			logger.WithYamlMarshaller(),
			logger.WithNumbersHighlight(colors.Blue),
			logger.WithKeyWordsHighlight(map[string]colors.ColorCode{
				logger.AttrError:     colors.Red,
				logger.AttrOperation: colors.Green,
			}),
		)

	default:
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		)

	}
	return slog.New(handler)
}
