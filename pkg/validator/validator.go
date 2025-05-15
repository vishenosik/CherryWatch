package validator

import (
	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
)

var (
	valid              = validator.New()
	ErrUuidNotProvided = errors.New("uuid not provided") //
)

func Struct(Struct any) error {
	return valid.Struct(Struct)
}

func UUID4(uuid string) error {
	if uuid == "" {
		return ErrUuidNotProvided
	}
	return valid.Var(uuid, "uuid4")
}
