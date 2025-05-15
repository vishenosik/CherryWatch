package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testContextKey struct{}

type testContext struct {
	testID string
}

func (ctx *testContext) Key() testContextKey {
	return testContextKey{}
}

func Test_Context(t *testing.T) {
	testID := "12345"
	requestCtx := &testContext{testID: testID}

	ctx := With(context.Background(), requestCtx)

	a := ctx.Value(testContextKey{})
	require.NotNil(t, a)

	actual, ok := From[*testContext](ctx)
	require.True(t, ok)
	require.NotNil(t, actual)

	require.Equal(t, testID, actual.testID)
}

func Test_RequestContext(t *testing.T) {

	requestID := "requestID"

	ctx := WithRequestCtx(context.Background(), requestID)

	actualGC, ok := RequestFromCtx(ctx)
	assert.True(t, ok)
	assert.Equal(t, requestID, actualGC.requestID)
}
