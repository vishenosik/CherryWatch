package sqlite

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	srv_models "github.com/vishenosik/CherryWatch/internal/services/models"
)

func Test_getAllEndpoints(t *testing.T) {

	x, cancel := suite(t)
	defer cancel()

	store := NewEndpoints(x)

	edps := srv_models.Endpoints{
		{
			ID:           "1",
			ServiceName:  "service1",
			URL:          "urlurl",
			SuccessCodes: []int{200, 201},
		},
		{
			ID:           "2",
			ServiceName:  "service2",
			URL:          "urlurl2",
			SuccessCodes: []int{},
		},
		{
			ID:           "3",
			ServiceName:  "service3",
			URL:          "urlur3",
			SuccessCodes: []int{204, 304},
		},
	}

	created, err := store.CreateEndpoints(context.Background(), edps...)
	require.NoError(t, err)
	require.Len(t, created, 3)

	actual, err := store.GetEndpoints()
	require.NoError(t, err)
	require.Len(t, actual, 3)

	actual, err = store.GetEndpoints("1", "2")
	require.NoError(t, err)
	require.Len(t, actual, 2)

	actual, err = store.GetEndpoints("1")
	require.NoError(t, err)
	require.Len(t, actual, 1)

	assert.Equal(t, actual[0].ID, "1")
	assert.Equal(t, actual[0].URL, "urlurl")
	assert.Equal(t, actual[0].ServiceName, "service1")
	assert.Equal(t, actual[0].SuccessCodes, []int{200, 201})
}

func Test_createEndpoints(t *testing.T) {

	x, cancel := suite(t)
	defer cancel()

	store := NewEndpoints(x)

	edps := srv_models.Endpoints{
		{
			ID:           "1",
			ServiceName:  "service1",
			URL:          "urlurl",
			SuccessCodes: []int{200, 201},
		},
		{
			ID:           "1",
			ServiceName:  "service2",
			SuccessCodes: []int{},
		},
		{
			ID:           "1",
			ServiceName:  "service3",
			SuccessCodes: []int{204, 304},
		},
	}

	created, err := store.CreateEndpoints(context.Background(), edps...)
	log.Println(created, err)
	require.Error(t, err)
	require.Len(t, created, 1)

	actual, err := store.GetEndpoints()
	require.NoError(t, err)
	require.Len(t, actual, 1)
}
