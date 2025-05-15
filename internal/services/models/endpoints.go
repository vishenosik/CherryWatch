package models

import (
	"time"

	"github.com/vishenosik/web/multierr"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
)

type Endpoint struct {
	// Endpoint identifier (uuid4 only)
	ID string
	// Name of checked service (ascii symbols only)
	ServiceName string
	// URL string to trigger during checks
	URL string
	// HTTP codes which are considered successful (should )
	SuccessCodes []int
	// Services used to notify about check failure
	NotificationServices []string
	// Time interval between checks
	Interval time.Duration
}

type Endpoints = []*Endpoint

func (edp *Endpoint) Validate() error {

	errs := new(multierr.Error)
	valid := validator.New()

	if edp.ServiceName == "" {
		errs.Append(errors.Wrap(ErrRequired, "service name"))
		return errs.ErrorOrNil()
	}

	if err := valid.Var(edp.ServiceName, "ascii"); err != nil {
		errs.Append(errors.Wrap(ErrAscii, "service name"))
		return errs.ErrorOrNil()
	}

	if err := valid.Var(edp.URL, "url"); err != nil {
		errs.Append(ErrURL)
	}

	if edp.ID == "" {
		errs.Append(ErrID)
	}

	if edp.Interval < time.Minute {
		errs.Append(ErrInterval)
	}

	for _, code := range edp.SuccessCodes {
		if code <= 0 || code >= 600 {
			errs.AppendWrapf(ErrCode, "code %d", code)
			break
		}
	}

	return errs.ErrorOrNil()
}

// FilterValidEndpoints takes a slice of Endpoints and returns a new slice
// containing only the endpoints that pass validation
func FilterValidEndpoints(endpoints Endpoints) (Endpoints, error) {

	if len(endpoints) == 0 {
		return nil, errors.Wrap(ErrContentNotAdded, "endpoints")
	}

	errs := new(multierr.Error)
	validEndpoints := make(Endpoints, 0, len(endpoints))

	for _, endpoint := range endpoints {
		err := endpoint.Validate()
		if err != nil {
			errs.AppendWrapf(err, "name=%s url=%s", endpoint.ServiceName, endpoint.URL)
			continue
		}
		validEndpoints = append(validEndpoints, endpoint)
	}

	return validEndpoints, errs.ErrorOrNil()
}
