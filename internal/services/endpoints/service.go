package endpoints

import (
	"context"
	"log/slog"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/vishenosik/CherryWatch/internal/services/models"
)

type EndpointsSaver interface {
	CreateEndpoints(
		ctx context.Context,
		edps ...*models.Endpoint,
	) (models.Endpoints, error)
}

type service struct {
	endpointsSaver EndpointsSaver
	tokenTTL       time.Duration
	tasksCH        chan models.Task
}

type Config struct {
	TokenTTL time.Duration
}

// New returns a new instance of Auth
func NewService(
	logger *slog.Logger,
	config Config,
	endpointsSaver EndpointsSaver,
) *service {

	return &service{
		tokenTTL:       config.TokenTTL,
		endpointsSaver: endpointsSaver,
		tasksCH:        make(chan models.Task, 1024),
	}
}

func (srv *service) TasksChan() chan models.Task {
	return srv.tasksCH
}

func (srv *service) SaveEndpoints(
	ctx context.Context,
	endpoints models.Endpoints,
) (added models.Endpoints, err error) {

	var errs *multierror.Error

	filtered, err := models.FilterValidEndpoints(endpoints)
	if err != nil {
		if validErrs, ok := err.(*multierror.Error); ok {
			errs = multierror.Append(errs, validErrs.Errors...)
		} else {
			errs = multierror.Append(errs, err)
		}
	}

	if len(filtered) == 0 {
		errs = multierror.Append(errs, models.ErrNothingToAdd)
		return nil, errs.ErrorOrNil()
	}

	created, err := srv.endpointsSaver.CreateEndpoints(ctx, filtered...)
	if err != nil {
		if storeErrs, ok := err.(*multierror.Error); ok {
			errs = multierror.Append(errs, storeErrs.Errors...)
		} else {
			errs = multierror.Append(errs, err)
		}
	}

	return created, errs.ErrorOrNil()
}
