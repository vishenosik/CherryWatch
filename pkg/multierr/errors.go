package multierr

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type Error struct {
	*multierror.Error
}

func (er *Error) List() []string {

	if er.Error == nil {
		return nil
	}

	var errors []string

	for _, err := range er.Error.Errors {
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
		er.Error = multierror.Append(er.Error, err)
		return
	}

	for _, _err := range errs.Errors {
		er.Error = multierror.Append(er.Error, _err)
	}
}

func (er *Error) AppendWrap(err error, message string) {

	if err == nil {
		return
	}

	errs, ok := err.(*multierror.Error)
	if !ok {
		er.Error = multierror.Append(er.Error, errors.Wrap(err, message))
		return
	}

	for _, _err := range errs.Errors {
		er.Error = multierror.Append(er.Error, errors.Wrap(_err, message))
	}
}

func (er *Error) AppendWrapf(err error, format string, args ...any) {
	if err == nil {
		return
	}

	errs, ok := err.(*multierror.Error)
	if !ok {
		er.Error = multierror.Append(er.Error, errors.Wrapf(err, format, args...))
		return
	}

	for _, _err := range errs.Errors {
		er.Error = multierror.Append(er.Error, errors.Wrapf(_err, format, args...))
	}
}
