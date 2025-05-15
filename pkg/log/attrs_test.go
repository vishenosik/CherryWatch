package log

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Error(t *testing.T) {
	ErrTest := errors.New("test error")
	result := Error(ErrTest)
	assert.Equal(t, AttrError, result.Key)
	assert.Equal(t, ErrTest.Error(), result.Value.String())
}

func Test_Operation(t *testing.T) {
	op := "test"
	result := Operation(op)
	assert.Equal(t, AttrOperation, result.Key)
	assert.Equal(t, op, result.Value.String())
}
