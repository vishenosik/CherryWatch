package errors

import (
	// pkg
	"github.com/pkg/errors"
)

type ErrorsMap[Type any] struct {
	data map[error]Type
	def  Type
}

func NewErrorsMap[Type any](Initial map[error]Type, Default Type) *ErrorsMap[Type] {
	return &ErrorsMap[Type]{
		data: Initial,
		def:  Default,
	}
}

func (em *ErrorsMap[Type]) Get(err error) Type {
	for er := range em.data {
		if errors.Is(err, er) {
			return em.data[er]
		}
	}
	return em.def
}
