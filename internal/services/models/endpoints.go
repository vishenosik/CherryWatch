package models

import (
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
)

var (
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

func (ep *Endpoint) Validate() error {

	var errs *multierror.Error

	if ep.Interval < time.Minute {
		errs = multierror.Append(errs, ErrInterval)
	}

	for _, code := range ep.SuccessCodes {
		if code <= 0 || code >= 600 {
			errs = multierror.Append(errs, errors.Wrapf(ErrCode, "code %d", code))
		}
	}

	valid := validator.New()

	if err := valid.Var(ep.URL, "url"); err != nil {
		errs = multierror.Append(errs, ErrURL)
	}

	if err := valid.Var(ep.ServiceName, "ascii"); err != nil {
		errs = multierror.Append(errs, ErrAscii)
	}

	return errs.ErrorOrNil()
}
