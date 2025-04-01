package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_EndpointValidation(t *testing.T) {
	testingTable := []struct {
		name        string
		endpoint    *Endpoint
		expectError bool
	}{
		{
			name: "valid endpoint",
			endpoint: &Endpoint{
				ID:                   uuid.NewString(),
				ServiceName:          "valid_service",
				URL:                  "https://example.com",
				SuccessCodes:         []int{200, 201, 301},
				NotificationServices: []string{"slack", "email"},
				Interval:             30 * time.Minute,
			},
		},
		{
			name: "invalid interval",
			endpoint: &Endpoint{
				ID:                   uuid.NewString(),
				ServiceName:          "valid_service",
				URL:                  "https://example.com",
				SuccessCodes:         []int{200, 201, 301},
				NotificationServices: []string{"slack", "email"},
				Interval:             30 * time.Second,
			},
			expectError: true,
		},
		{
			name: "invalid UUID",
			endpoint: &Endpoint{
				ID:          "not-a-uuid",
				ServiceName: "valid_service",
				URL:         "https://example.com",
				Interval:    30 * time.Second,
			},
			expectError: true,
		},
		{
			name: "invalid service name - non-ascii",
			endpoint: &Endpoint{
				ID:          uuid.NewString(),
				ServiceName: "服务", // Chinese characters
				URL:         "https://example.com",
				Interval:    30 * time.Second,
			},
			expectError: true,
		},
		{
			name: "invalid URL",
			endpoint: &Endpoint{
				ID:          uuid.NewString(),
				ServiceName: "valid_service",
				URL:         "not-a-url",
				Interval:    30 * time.Second,
			},
			expectError: true,
		},
		{
			name: "invalid success code - too low",
			endpoint: &Endpoint{
				ID:           uuid.NewString(),
				ServiceName:  "valid_service",
				URL:          "https://example.com",
				SuccessCodes: []int{99}, // Below minimum
				Interval:     30 * time.Second,
			},
			expectError: true,
		},
		{
			name: "invalid success code - too high",
			endpoint: &Endpoint{
				ID:           uuid.NewString(),
				ServiceName:  "valid_service",
				URL:          "https://example.com",
				SuccessCodes: []int{600}, // Above maximum
				Interval:     30 * time.Second,
			},
			expectError: true,
		},
		{
			name: "invalid success code range format",
			endpoint: &Endpoint{
				ID:          uuid.NewString(),
				ServiceName: "valid_service",
				URL:         "https://example.com",
				Interval:    30 * time.Second,
			},
			expectError: true,
		},
		{
			name: "invalid success code range values",
			endpoint: &Endpoint{
				ID:          uuid.NewString(),
				ServiceName: "valid_service",
				URL:         "https://example.com",
				Interval:    30 * time.Second,
			},
			expectError: true,
		},
		{
			name: "invalid interval - too short",
			endpoint: &Endpoint{
				ID:          uuid.NewString(),
				ServiceName: "valid_service",
				URL:         "https://example.com",
				Interval:    5 * time.Second, // Potentially too frequent
			},
			expectError: true,
		},
	}

	for _, tt := range testingTable {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.endpoint.Validate()
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func Test_EndpointEdgeCases(t *testing.T) {
	t.Run("empty notification services", func(t *testing.T) {
		e := &Endpoint{
			ID:          uuid.NewString(),
			ServiceName: "test_service",
			URL:         "https://example.com",
			Interval:    30 * time.Minute,
		}
		if err := e.Validate(); err != nil {
			t.Errorf("empty notification services should be valid: %v", err)
		}
	})

	t.Run("empty success codes", func(t *testing.T) {
		e := &Endpoint{
			ID:          uuid.NewString(),
			ServiceName: "test_service",
			URL:         "https://example.com",
			Interval:    30 * time.Minute,
		}
		if err := e.Validate(); err != nil {
			t.Errorf("empty success codes should be valid: %v", err)
		}
	})
}
