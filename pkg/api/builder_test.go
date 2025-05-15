package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestingTable = []struct {
	name       string
	version    uint8
	routeParts []string
	expect     string
}{
	{"1", 1, []string{"/test"}, "/api/v1/test"},
	{"2", 1, []string{"test"}, "/api/v1/test"},
	{"3", 1, []string{"test/"}, "/api/v1/test"},
	{"4", 1, []string{"//test"}, "/api/v1/test"},
	{"5", 1, []string{"//test/direct"}, "/api/v1/test/direct"},
	{"6", 1, []string{"//test", "non-direct"}, "/api/v1/test/non-direct"},
	{"7", 1, nil, "/api/v1"},
	{"8", 1, []string{}, "/api/v1"},
}

func Test_BuildApi(t *testing.T) {

	t.Helper()
	t.Parallel()

	for _, tt := range TestingTable {
		t.Run(tt.name, func(t *testing.T) {
			Api := buildApi(tt.version, tt.routeParts...)
			assert.Equal(t, tt.expect, Api)
		})
	}
}

func Test_ApiV1(t *testing.T) {

	t.Helper()
	t.Parallel()

	for _, tt := range TestingTable {
		if tt.version != 1 {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			Api := ApiV1(tt.routeParts...)
			assert.Equal(t, tt.expect, Api)
		})
	}
}
