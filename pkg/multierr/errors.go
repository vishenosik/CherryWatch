package multierr

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type Error struct {
	critical error
	errs     *multierror.Error
}

func (er *Error) ErrorOrNil() error {
	if er.errs == nil || len(er.errs.Errors) == 0 {
		return nil
	}

	return er
}

func (er *Error) Error() string {
	return er.errs.Error()
}

func (er *Error) Critical() error {
	return er.critical
}

func (er *Error) CriticalString() string {
	if er.critical != nil {
		return er.critical.Error()
	}
	return ""
}

func (er *Error) List() []string {
	if er.errs == nil {
		return nil
	}

	var errors []string

	for _, err := range er.errs.Errors {
		errors = append(errors, err.Error())
	}
	return errors
}

func (er *Error) append(err error, critical bool, wrapper func(error) error) {
	if err == nil {
		return
	}

	if critical {
		er.critical = err
	}

	switch err := err.(type) {
	case *multierror.Error:
		for _, _err := range err.Errors {
			er.errs = multierror.Append(er.errs, wrapper(_err))
		}
	case *Error:
		er.critical = err.critical
		er.errs = multierror.Append(er.errs, wrapper(err.errs))
	default:
		er.errs = multierror.Append(er.errs, wrapper(err))
	}
}

func (er *Error) Append(err error) {
	er.append(err, false, func(err error) error { return err })
}

func (er *Error) AppendWrap(err error, message string) {
	er.append(err, false, func(err error) error { return errors.Wrap(err, message) })
}

func (er *Error) AppendWrapf(err error, format string, args ...any) {
	er.append(err, false, func(err error) error { return errors.Wrapf(err, format, args...) })
}

func (er *Error) AppendCritical(err error) {
	er.append(err, true, func(err error) error { return err })
}

func (er *Error) AppendCriticalWrap(err error, message string) {
	er.append(err, true, func(err error) error { return errors.Wrap(err, message) })
}

func (er *Error) AppendCriticalWrapf(err error, format string, args ...any) {
	er.append(err, true, func(err error) error { return errors.Wrapf(err, format, args...) })
}
