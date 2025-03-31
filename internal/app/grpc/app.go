package grpcApp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/pkg/errors"
	"github.com/vishenosik/web-tools/config"
	"google.golang.org/grpc"
)

// App represents the gRPC application structure.
// It encapsulates the core components needed to run a gRPC server.
type App struct {
	// log is a structured logger for the application.
	log *slog.Logger
	// server is the main gRPC server instance.
	server *grpc.Server
	// port is the network port on which the gRPC server will listen.
	port uint16
}

type Config struct {
	Server config.Server
}

// NewGrpcApp creates and initializes a new gRPC application.
//
// It sets up a gRPC server with authentication services and configures logging.
//
// Parameters:
//   - log: A pointer to a slog.Logger for application logging.
//   - conf: A GRPCConfig struct containing the gRPC server configuration.
//   - authService: An Authentication interface for handling authentication operations.
//
// Returns:
//   - *App: A pointer to the newly created App struct, ready to be run.
func NewGrpcApp(
	log *slog.Logger,
	config Config,
	// authService authentication.Authentication,
) *App {

	Log := log.WithGroup(
		"gRPC",
	)

	server := grpc.NewServer()

	// authentication_v1.RegisterAuthenticationServer(
	// 	server,
	// 	authentication.NewAuthenticationServer(
	// 		Log,
	// 		authService,
	// 	),
	// )

	return &App{
		log:    Log,
		server: server,
		port:   config.Server.Port,
	}
}

// MustRun starts the gRPC server and panics if an error occurs during startup.
// This method is useful for scenarios where the application should not continue
// if the gRPC server fails to start.
//
// It internally calls the Run method and panics if an error is returned.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run starts the gRPC server and listens for incoming connections.
//
// This method sets up a TCP listener on the configured port, starts the gRPC server,
// and handles any errors that occur during the process. It logs the server's startup
// and running status.
//
// The method uses the App's configured logger to provide context-aware logging,
// including the operation name and port number.
//
// Returns:
//   - error: An error if the server fails to start or encounters any issues while running.
//     Returns nil if the server starts and runs successfully.
func (a *App) Run() error {
	const op = "grpcApp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Any("port", a.port),
	)

	log.Info("starting gRPC server")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return errors.Wrap(err, op)
	}

	log.Info("gRPC server is running", slog.String("addr", listener.Addr().String()))

	if err := a.server.Serve(listener); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

// Stop gracefully stops the gRPC server.
//
// This method logs the stop operation and initiates a graceful shutdown of the gRPC server.
// It ensures that ongoing requests are allowed to complete before the server stops.
//
// The method does not take any parameters and does not return any value.
func (a *App) Stop(_ context.Context) {

	const op = "grpcApp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Any("port", a.port))

	a.server.GracefulStop()

}
