package models

import "github.com/pkg/errors"

var (
	// no content added
	ErrContentNotAdded = errors.New("content is not added")
)

// validation errors
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
	// is required
	ErrRequired = errors.New("is required")
)
