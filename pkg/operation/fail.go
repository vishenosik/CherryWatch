package operation

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FailWrapError is a generic function that returns a closure that takes an error as input,
// wraps it with the provided operation string using github.com/pkg/errors.Wrap, and returns
// the original result along with the wrapped error.
//
// This function is useful for handling errors in a consistent and flexible manner,
// allowing you to wrap errors with descriptive operation strings without modifying the original result.
//
// The function takes two parameters:
// - result: The original result of type Type. This value will be returned along with the wrapped error.
// - op: A string representing the operation that failed. This string will be used to wrap the error.
//
// The function returns a closure that takes an error as input and returns a tuple containing:
// - The original result of type Type.
// - An error that is the result of wrapping the input error with the provided operation string.
func FailWrapError[Type any](result Type, op string) func(err error) (Type, error) {
	return func(err error) (Type, error) {
		return result, errors.Wrap(err, op)
	}
}

// FailNilWrapError is a convenience function that returns a closure that wraps an error with a provided operation string.
// It uses the FailWrapError function to achieve this, but simplifies the process by setting the result to nil.
//
// This function is useful for handling errors in a consistent and flexible manner,
// allowing you to wrap errors with descriptive operation strings without modifying the original result.
//
// The function takes a single parameter:
// - op (string): A string representing the operation that failed. This string will be used to wrap the error.
//
// The function returns a closure that takes an error as input and returns a tuple containing:
// - nil: The original result of type any is set to nil.
// - An error that is the result of wrapping the input error with the provided operation string.
func FailNilWrapError(op string) func(err error) (any, error) {
	return FailWrapError[any](nil, op)
}

func FailWrapErrorStatus[Type any](result Type, message string) func(code codes.Code) (Type, error) {
	return func(code codes.Code) (Type, error) {
		return result, status.Error(code, message)
	}
}

func FailNilWrapErrorStatus(message string) func(code codes.Code) (any, error) {
	return FailWrapErrorStatus[any](nil, message)
}
