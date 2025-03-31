package models

import "github.com/pkg/errors"

var (
	// not found
	ErrNotFound = errors.New("not found")
	// exists already
	ErrAlreadyExists = errors.New("exists already")
)
