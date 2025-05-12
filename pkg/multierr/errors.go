package multierr

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type Error struct {
	Err *multierror.Error
}

func (er *Error) ErrorOrNil() error {
	return er.Err.ErrorOrNil()
}

func (er *Error) Error() string {
	return er.Err.GoString()
}

func (er *Error) List() []string {

	if er.Err == nil {
		return nil
	}

	var errors []string

	for _, err := range er.Err.Errors {
		errors = append(errors, err.Error())
	}
	return errors
}

func (er *Error) Append(err error) {

	if err == nil {
		return
	}

	errs, ok := err.(*multierror.Error)
	if !ok {
		er.Err = multierror.Append(er.Err, err)
		return
	}

	for _, _err := range errs.Errors {
		er.Err = multierror.Append(er.Err, _err)
	}
}

func (er *Error) AppendWrap(err error, message string) {

	if err == nil {
		return
	}

	errs, ok := err.(*multierror.Error)
	if !ok {
		er.Err = multierror.Append(er.Err, errors.Wrap(err, message))
		return
	}

	for _, _err := range errs.Errors {
		er.Err = multierror.Append(er.Err, errors.Wrap(_err, message))
	}
}

func (er *Error) AppendWrapf(err error, format string, args ...any) {
	if err == nil {
		return
	}

	errs, ok := err.(*multierror.Error)
	if !ok {
		er.Err = multierror.Append(er.Err, errors.Wrapf(err, format, args...))
		return
	}

	for _, _err := range errs.Errors {
		er.Err = multierror.Append(er.Err, errors.Wrapf(_err, format, args...))
	}
}
