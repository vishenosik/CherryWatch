package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIVersion_ParseFirst(t *testing.T) {
	tests := []struct {
		name        string
		request     func() *http.Request
		expectError bool
		expected    APIVersion
	}{
		{
			name: "Valid query param",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/api?v=2.1", nil)
				return req
			},
			expected: APIVersion{Major: 2, Minor: 1},
		},
		{
			name: "Valid header",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/api", nil)
				req.Header.Set(VersionHeader, "1.0")
				return req
			},
			expected: APIVersion{Major: 1, Minor: 0},
		},
		{
			name: "Default version",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/api", nil)
			},
			expected: APIVersion{Major: 2, Minor: 1},
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
			name: "Version too low",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/api?v=0.9", nil)
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
			var av APIVersion
			err := av.ParseFirst(tt.request())

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if av != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, av)
			}
		})
	}
}

func TestAPIVersion_WithContext(t *testing.T) {
	av := APIVersion{Major: 1, Minor: 0}
	ctx := av.WithContext(context.Background())

	retrieved, err := ApiVersionFromContext(ctx)
	if err != nil {
		t.Fatalf("Failed to retrieve version from context: %v", err)
	}

	if *retrieved != av {
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

func TestApiVersionMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		request        func() *http.Request
		expectError    bool
		expectedStatus int
	}{
		{
			name: "Valid request",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/api?v=2.1", nil)
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid version",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/api?v=invalid", nil)
				return req
			},
			expectError:    true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handlerCalled := false

			middleware := ApiVersionMiddleware(&APIVersion{})
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				_, err := ApiVersionFromContext(r.Context())
				if err != nil {
					t.Errorf("Failed to get version from context: %v", err)
				}
			})

			middleware(testHandler).ServeHTTP(rr, tt.request())

			if tt.expectError {
				if handlerCalled {
					t.Error("Handler was called but shouldn't have been")
				}

				if rr.Code != tt.expectedStatus {
					t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
				}

				var errResp ErrorResponse
				if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
					t.Errorf("Failed to decode error response: %v", err)
				}
			} else {
				if !handlerCalled {
					t.Error("Handler was not called")
				}
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected APIVersion
		err      bool
	}{
		{"2.1", APIVersion{2, 1}, false},
		{"1.0", APIVersion{1, 0}, false},
		{"invalid", APIVersion{}, true},
		{"1", APIVersion{}, true},
		{"1.2.3", APIVersion{}, true},
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

			if *result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, *result)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       *APIVersion
		v2       *APIVersion
		expected int
	}{
		{"Equal", &APIVersion{1, 0}, &APIVersion{1, 0}, 0},
		{"Major less", &APIVersion{1, 0}, &APIVersion{2, 0}, -1},
		{"Major greater", &APIVersion{2, 0}, &APIVersion{1, 0}, 1},
		{"Minor less", &APIVersion{1, 0}, &APIVersion{1, 1}, -1},
		{"Minor greater", &APIVersion{1, 1}, &APIVersion{1, 0}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}
