package endpoints

import (
	"context"
	"log/slog"
	"time"

	"github.com/vishenosik/CherryWatch/internal/services/models"
	"github.com/vishenosik/CherryWatch/pkg/multierr"
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

	errs := new(multierr.Error)

	filtered, err := models.FilterValidEndpoints(endpoints)
	if err != nil {
		errs.Append(err)
	}

	if len(filtered) == 0 {
		errs.Append(models.ErrContentNotAdded)
		return nil, errs.ErrorOrNil()
	}

	created, err := srv.endpointsSaver.CreateEndpoints(ctx, filtered...)
	if err != nil {
		errs.Append(err)
	}

	if len(created) == 0 {
		errs.Append(models.ErrContentNotAdded)
		return nil, errs.ErrorOrNil()
	}

	return created, errs.ErrorOrNil()
}
