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
		edps models.Endpoints,
	) error
}

type service struct {
	endpointsSaver EndpointsSaver
	tokenTTL       time.Duration
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
	}
}

func (srv *service) SaveEndpoints(
	ctx context.Context,
	endpoints models.Endpoints,
) (added models.Endpoints, err error) {
	var validationErrs *multierror.Error

	filtered, err := models.FilterValidEndpoints(endpoints)
	if err != nil {
		if errs, ok := err.(*multierror.Error); ok {
			validationErrs = errs
		} else {
			return nil, err
		}
	}

	if len(filtered) == 0 {
		return nil, models.ErrNothingToAdd
	}

	if err := srv.endpointsSaver.CreateEndpoints(ctx, filtered); err != nil {
		return nil, err
	}

	return filtered, validationErrs.ErrorOrNil()
}
