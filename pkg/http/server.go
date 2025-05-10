package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/pkg/errors"

	"github.com/vishenosik/web/config"
	logger "github.com/vishenosik/web/log"
)

type Server struct {
	log    *slog.Logger
	server *http.Server
	port   uint16
}

type Config struct {
	Server config.Server
}

func NewHttpApp(config Config, logger *slog.Logger, handler http.Handler) *Server {
	return NewHttpAppContext(context.Background(), config, logger, handler)
}

func NewHttpAppContext(
	ctx context.Context,
	config Config,
	logger *slog.Logger,
	handler http.Handler,
) *Server {

	err := config.Server.Validate()
	if err != nil {
		panic(errors.Wrap(err, "failed to validate http app config"))
	}

	if logger == nil {
		panic("logger can't be nil")
	}

	if handler == nil {
		panic("handler can't be nil")
	}

	return &Server{
		log: logger,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Server.Port),
			Handler: handler,
		},
		port: config.Server.Port,
	}
}

func (a *Server) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *Server) Run() error {
	const op = "http.Server.Run"

	log := a.log.With(logger.Operation(op), slog.Any("port", a.port))

	log.Info("starting server")

	if err := a.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, op)
		}
	}

	log.Info("server is running")
	return nil
}

func (a *Server) Stop(ctx context.Context) {

	const op = "http.Server.Stop"

	a.log.Info("stopping server", logger.Operation(op), slog.Any("port", a.port))

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("server shutdown failed", logger.Error(err))
	}
}
