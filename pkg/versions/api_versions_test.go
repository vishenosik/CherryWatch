package versions

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIVersion_ParseRequest(t *testing.T) {
	tests := []struct {
		name        string
		request     func() *http.Request
		expectError bool
		expected    *Version
	}{
		{
			name: "Valid query param",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/api?v=2.1", nil)
				return req
			},
			expected: &Version{Major: 2, Minor: 1},
		},
		{
			name: "Valid header",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/api", nil)
				req.Header.Set(VersionHeader, "1.0")
				return req
			},
			expected: &Version{Major: 1, Minor: 0},
		},
		{
			name: "Invalid format",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/api?v=invalid", nil)
				return req
			},
			expectError: true,
		},
		{
			name: "Version too high",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/api?v=3.0", nil)
				return req
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av, err := NewApiVersion("2.5")
			assert.NoError(t, err)

			err = av.ParseRequest(tt.request())
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, av.Version)
			}
		})
	}
}

func TestAPIVersion_WithContext(t *testing.T) {
	av := &APIVersion{Version: &Version{Major: 1, Minor: 0}}
	ctx := av.WithContext(context.Background())

	retrieved, err := ApiVersionFromContext(ctx)
	if err != nil {
		t.Fatalf("Failed to retrieve version from context: %v", err)
	}

	if retrieved != av.Version {
		t.Errorf("Expected %v, got %v", av, *retrieved)
	}
}

func TestApiVersionFromContext_Error(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{"No version", context.Background()},
		{"Wrong type", context.WithValue(context.Background(), apiVersionKey{}, "not a version")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ApiVersionFromContext(tt.ctx)
			if err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected *Version
		err      bool
	}{
		{"2.1", &Version{2, 1}, false},
		{"1.0", &Version{1, 0}, false},
		{"invalid", &Version{}, true},
		{"1", &Version{}, true},
		{"1.2.3", &Version{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseVersion(tt.input)
			if tt.err {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			require.Equal(t, tt.expected, result)

		})
	}
}

func Test_VersionGTE(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Version
		v2       *Version
		expected bool
	}{
		{"Equal", &Version{1, 0}, &Version{1, 0}, true},
		{"Major less", &Version{1, 0}, &Version{2, 0}, false},
		{"Major greater", &Version{2, 0}, &Version{1, 0}, true},
		{"Minor less", &Version{1, 0}, &Version{1, 1}, false},
		{"Minor greater", &Version{1, 1}, &Version{1, 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.GTE(tt.v2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
