package models

import (
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
)

var (
	// ID can't be empty
	ErrID = errors.New("ID can't be empty")
	// provided string is not URL
	ErrURL = errors.New("provided string is not URL")
	// time interval can't be less than time.Minute
	ErrInterval = errors.New("time interval can't be less than time.Minute")
	// string must consist of only ascii characters
	ErrAscii = errors.New("string must consist of only ascii characters")
	// must be in (0,600) interval
	ErrCode = errors.New("must be in (0,600) interval")
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

	var errs *multierror.Error

	if edp.ID == "" {
		errs = multierror.Append(errs, ErrID)
	}

	if edp.Interval < time.Minute {
		errs = multierror.Append(errs, ErrInterval)
	}

	for _, code := range edp.SuccessCodes {
		if code <= 0 || code >= 600 {
			errs = multierror.Append(errs, errors.Wrapf(ErrCode, "code %d", code))
		}
	}

	valid := validator.New()

	if err := valid.Var(edp.URL, "url"); err != nil {
		errs = multierror.Append(errs, ErrURL)
	}

	if err := valid.Var(edp.ServiceName, "ascii"); err != nil {
		errs = multierror.Append(errs, ErrAscii)
	}

	return errs.ErrorOrNil()
}

// FilterValidEndpoints takes a slice of Endpoints and returns a new slice
// containing only the endpoints that pass validation
func FilterValidEndpoints(endpoints Endpoints) (Endpoints, error) {

	if len(endpoints) == 0 {
		return nil, errors.Wrap(ErrNothingToAdd, "endpoints")
	}

	var errs *multierror.Error
	validEndpoints := make(Endpoints, 0, len(endpoints))

	for _, endpoint := range endpoints {
		err := endpoint.Validate()
		if err != nil {
			errs = multierror.Append(errs, errors.Wrap(err, endpoint.ID))
			continue
		}
		validEndpoints = append(validEndpoints, endpoint)

	}

	return validEndpoints, errs.ErrorOrNil()
}
