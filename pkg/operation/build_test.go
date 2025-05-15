package operation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildOperation(t *testing.T) {

	t.Helper()
	t.Parallel()

	expect := "layer.service.method"
	result := buildOperation("layer", "service", "method")
	assert.Equal(t, expect, result)
}

func Test_ServicesOperation(t *testing.T) {

	t.Helper()
	t.Parallel()

	expect := fmt.Sprintf("%s.service.method", servicesLayer)
	result := ServicesOperation("service", "method")
	assert.Equal(t, expect, result)
}
