package validator

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_UUID4(t *testing.T) {

	UUID := uuid.New().String()
	err := UUID4(UUID)
	assert.NoError(t, err)

	err = UUID4("")
	assert.Error(t, err)

	err = UUID4("dfvfvdfv")
	assert.Error(t, err)

}
