package operation

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_ReturnFailWithError(t *testing.T) {

	op := "op"
	Err := errors.New("test error")

	String := "string"
	result1, err := FailWrapError(String, op)(Err)
	require.Equal(t, String, result1)
	require.ErrorIs(t, err, Err)

	Bool := false
	result2, _ := FailWrapError(Bool, op)(Err)
	require.Equal(t, Bool, result2)

	Int := 9
	result3, _ := FailWrapError(Int, op)(Err)
	require.Equal(t, Int, result3)

	result4, _ := FailWrapError[any](nil, op)(Err)
	require.Equal(t, nil, result4)

	result4, _ = FailNilWrapError(op)(Err)
	require.Equal(t, nil, result4)
}
