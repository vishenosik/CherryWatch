package endpoints

import (
	"context"
	"log/slog"
	"time"

	"github.com/vishenosik/CherryWatch/internal/services/models"
	logger "github.com/vishenosik/web/logs"
	"github.com/vishenosik/web/multierr"
)

type EndpointsSaver interface {
	CreateEndpoints(
		ctx context.Context,
		edps ...*models.Endpoint,
	) (models.Endpoints, error)
}

type service struct {
	log            *slog.Logger
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
		log:            logger,
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
	log := srv.log.With(
		logger.Operation("service.SaveEndpoints"),
	)

	log.Info("filter endpoints", slog.Int("count", len(endpoints)))

	filtered, err := models.FilterValidEndpoints(endpoints)
	if err != nil {
		log.Warn("filtration errors", logger.Error(err))
		errs.Append(err)
	}

	if len(filtered) == 0 {
		log.Error("failed to filter endpoints", logger.Error(models.ErrContentNotAdded))
		errs.AppendCritical(models.ErrContentNotAdded)
		return nil, errs.ErrorOrNil()
	}

	log.Info("creating endpoints", slog.Int("count", len(filtered)))

	created, err := srv.endpointsSaver.CreateEndpoints(ctx, filtered...)
	if err != nil {
		log.Warn("creation errors", logger.Error(err))
		errs.Append(err)
	}

	if len(created) == 0 {
		log.Error("failed to create endpoints", logger.Error(models.ErrContentNotAdded))
		errs.AppendCritical(models.ErrContentNotAdded)
		return nil, errs.ErrorOrNil()
	}

	return created, errs.ErrorOrNil()
}
