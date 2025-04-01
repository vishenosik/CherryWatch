package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vishenosik/CherryWatch/internal/services/models"
)

var (
	requiredFieldsOnly = `
	{
  		"id": "550e8400-e29b-41d4-a716-446655440000",
  		"service_name": "user_service",
  		"url": "https://api.example.com/users",
  		"time_interval": 30000000000
	}`

	allFields = `
	{
  		"id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  		"service_name": "payment_gateway",
  		"url": "https://payments.example.com/status",
  		"success_codes": ["200", "201", "202", "204","200-299", "300-399"],
  		"notification_services": ["slack", "email", "sms"],
  		"time_interval": 60000000000
	}`

	codeRangesOnly = `
	{
  		"id": "123e4567-e89b-12d3-a456-426614174000",
  		"service_name": "inventory_service",
  		"url": "https://inventory.example.com/health",
  		"success_codes": ["200-299", "500-599"],
  		"time_interval": 15000000000
	}`

	emptyNotifications = `
	{
  		"id": "550e8400-e29b-41d4-a716-446655440000",
  		"service_name": "empty_notifications",
  		"url": "https://example.com",
  		"notification_services": [],
  		"time_interval": 30000000000
	}`

	nullFields = `
	{
	  	"id": "550e8400-e29b-41d4-a716-446655440000",
	  	"service_name": "null_fields_service",
	  	"url": "https://example.com",
	  	"success_codes": null,
	  	"notification_services": null,
	  	"time_interval": 30000000000
	}`
	incorrectEndpoint = `
	{
	  	"ids": ["550e8400-e29b-41d4-a716-446655440000"]
	}`
	incorrectEndpointType = `
	{
	  	"id": ["550e8400-e29b-41d4-a716-446655440000"]
	}`
)

func Test_endpointsJSON(t *testing.T) {

	testingTable := []struct {
		name        string
		data        string
		expectError bool
	}{
		{
			name: "required fields only",
			data: fmt.Sprintf("[%s]", requiredFieldsOnly),
		},
		{
			name: "all fields",
			data: fmt.Sprintf("[%s]", allFields),
		},
		{
			name: "null fields",
			data: fmt.Sprintf("[%s]", nullFields),
		},
		{
			name:        "fail structure",
			data:        incorrectEndpoint,
			expectError: true,
		},
		{
			name:        "fail type",
			data:        fmt.Sprintf("[%s]", incorrectEndpointType),
			expectError: true,
		},
		{
			name: "empty endpoints list",
			data: fmt.Sprintf("[%s]", ""),
		},
		{
			name: "pack of endpoints",
			data: fmt.Sprintf("[%s,%s,%s]", requiredFieldsOnly, allFields, nullFields),
		},
	}

	for _, tt := range testingTable {
		t.Run(tt.name, func(t *testing.T) {
			var endpoints Endpoints
			err := json.Unmarshal([]byte(tt.data), &endpoints)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_parseRanges(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		want      []int
		wantError bool
	}{
		{
			name:  "multiple ranges",
			input: []string{"200-202", "205-207"},
			want:  []int{200, 201, 202, 205, 206, 207},
		},
		{
			name:  "single value ranges",
			input: []string{"404-404", "500-500"},
			want:  []int{404, 500},
		},
		{
			name:      "invalid range",
			input:     []string{"200-abc"},
			wantError: true,
		},
		{
			name:  "empty input",
			input: []string{},
			want:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRanges(tt.input)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToServiceEndpoint(t *testing.T) {
	// Common test endpoint
	baseEndpoint := Endpoint{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		ServiceName: "test-service",
		URL:         "https://example.com",
		Interval:    30 * time.Second,
	}

	tests := []struct {
		name          string
		input         Endpoint
		expected      *models.Endpoint
		expectError   bool
		errorContains string
	}{
		{
			name: "success with only success codes",
			input: Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []string{"200", "201", "204"},
				Interval:     baseEndpoint.Interval,
			},
			expected: &models.Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []int{200, 201, 204},
				Interval:     baseEndpoint.Interval,
			},
		},
		{
			name: "success with only ranges",
			input: Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []string{"200-202", "404-404"},
				Interval:     baseEndpoint.Interval,
			},
			expected: &models.Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []int{200, 201, 202, 404},
				Interval:     baseEndpoint.Interval,
			},
		},
		{
			name: "success with codes and ranges",
			input: Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []string{"200", "500", "201-203", "404-405"},
				Interval:     baseEndpoint.Interval,
			},
			expected: &models.Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []int{200, 201, 202, 203, 404, 405, 500},
				Interval:     baseEndpoint.Interval,
			},
		},
		{
			name: "success with duplicate codes",
			input: Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []string{"200", "200", "201", "201-203"},
				Interval:     baseEndpoint.Interval,
			},
			expected: &models.Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []int{200, 201, 202, 203},
				Interval:     baseEndpoint.Interval,
			},
		},
		{
			name: "empty success codes and ranges",
			input: Endpoint{
				ID:          baseEndpoint.ID,
				ServiceName: baseEndpoint.ServiceName,
				URL:         baseEndpoint.URL,
				Interval:    baseEndpoint.Interval,
			},
			expected: &models.Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []int{},
				Interval:     baseEndpoint.Interval,
			},
		},
		{
			name: "invalid range format",
			input: Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []string{"200-abc"},
				Interval:     baseEndpoint.Interval,
			},
			expected: &models.Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []int{},
				Interval:     baseEndpoint.Interval,
			},
		},
		{
			name: "invalid range (start > end)",
			input: Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []string{"300-200"},
				Interval:     baseEndpoint.Interval,
			},
			expected: &models.Endpoint{
				ID:           baseEndpoint.ID,
				ServiceName:  baseEndpoint.ServiceName,
				URL:          baseEndpoint.URL,
				SuccessCodes: []int{},
				Interval:     baseEndpoint.Interval,
			},
		},
		{
			name: "with notification services",
			input: Endpoint{
				ID:                   baseEndpoint.ID,
				ServiceName:          baseEndpoint.ServiceName,
				URL:                  baseEndpoint.URL,
				SuccessCodes:         []string{"200"},
				NotificationServices: []string{"slack", "email"},
				Interval:             baseEndpoint.Interval,
			},
			expected: &models.Endpoint{
				ID:                   baseEndpoint.ID,
				ServiceName:          baseEndpoint.ServiceName,
				URL:                  baseEndpoint.URL,
				SuccessCodes:         []int{200},
				NotificationServices: []string{"slack", "email"},
				Interval:             baseEndpoint.Interval,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToServiceEndpoint(tt.input)

			if tt.expectError {
				// If your function doesn't currently return errors, you might want to:
				// 1. Change the function signature to return error
				// 2. Add logging for the error case
				// 3. Add assertions based on whatever error handling you implement
				t.Log("Note: Currently the function doesn't return errors, but invalid ranges should be handled")
				return
			}

			require.NotNil(t, result, "result should not be nil")
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.ServiceName, result.ServiceName)
			assert.Equal(t, tt.expected.URL, result.URL)
			assert.ElementsMatch(t, tt.expected.SuccessCodes, result.SuccessCodes)
			assert.Equal(t, tt.expected.NotificationServices, result.NotificationServices)
			assert.Equal(t, tt.expected.Interval, result.Interval)
		})
	}
}

func Test_codesRanges(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []int{},
			expected: []string{},
		},
		{
			name:     "single number",
			input:    []int{5},
			expected: []string{"5"},
		},
		{
			name:     "consecutive numbers",
			input:    []int{1, 2, 3, 4, 5},
			expected: []string{"1-5"},
		},
		{
			name:     "non-consecutive numbers",
			input:    []int{1, 3, 5, 7},
			expected: []string{"1", "3", "5", "7"},
		},
		{
			name:     "mixed consecutive and non-consecutive",
			input:    []int{1, 2, 3, 5, 7, 8, 9, 10, 15},
			expected: []string{"1-3", "5", "7-10", "15"},
		},
		{
			name:     "multiple ranges",
			input:    []int{10, 11, 12, 14, 15, 16, 20, 21, 22},
			expected: []string{"10-12", "14-16", "20-22"},
		},
		{
			name:     "negative numbers",
			input:    []int{-3, -2, -1, 0, 1, 5},
			expected: []string{"-3-1", "5"},
		},
		{
			name:     "duplicate numbers",
			input:    []int{1, 1, 2, 2, 3, 5, 5},
			expected: []string{"1-3", "5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := codesRanges(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFromServiceEndpoint(t *testing.T) {
	baseModel := &models.Endpoint{
		ID:                   "550e8400-e29b-41d4-a716-446655440000",
		ServiceName:          "test-service",
		URL:                  "https://example.com",
		NotificationServices: []string{"slack", "email"},
		Interval:             30 * time.Second,
	}

	tests := []struct {
		name     string
		input    *models.Endpoint
		expected Endpoint
	}{
		{
			name: "with consecutive success codes",
			input: &models.Endpoint{
				ID:                   baseModel.ID,
				ServiceName:          baseModel.ServiceName,
				URL:                  baseModel.URL,
				SuccessCodes:         []int{200, 201, 202},
				NotificationServices: baseModel.NotificationServices,
				Interval:             baseModel.Interval,
			},
			expected: Endpoint{
				ID:                   baseModel.ID,
				ServiceName:          baseModel.ServiceName,
				URL:                  baseModel.URL,
				SuccessCodes:         []string{"200-202"},
				NotificationServices: baseModel.NotificationServices,
				Interval:             baseModel.Interval,
			},
		},
		{
			name: "with non-consecutive success codes",
			input: &models.Endpoint{
				ID:                   baseModel.ID,
				ServiceName:          baseModel.ServiceName,
				URL:                  baseModel.URL,
				SuccessCodes:         []int{200, 404, 500},
				NotificationServices: baseModel.NotificationServices,
				Interval:             baseModel.Interval,
			},
			expected: Endpoint{
				ID:                   baseModel.ID,
				ServiceName:          baseModel.ServiceName,
				URL:                  baseModel.URL,
				SuccessCodes:         []string{"200", "404", "500"},
				NotificationServices: baseModel.NotificationServices,
				Interval:             baseModel.Interval,
			},
		},
		{
			name: "with mixed success codes",
			input: &models.Endpoint{
				ID:                   baseModel.ID,
				ServiceName:          baseModel.ServiceName,
				URL:                  baseModel.URL,
				SuccessCodes:         []int{500, 200, 201, 204, 205},
				NotificationServices: baseModel.NotificationServices,
				Interval:             baseModel.Interval,
			},
			expected: Endpoint{
				ID:                   baseModel.ID,
				ServiceName:          baseModel.ServiceName,
				URL:                  baseModel.URL,
				SuccessCodes:         []string{"200-201", "204-205", "500"},
				NotificationServices: baseModel.NotificationServices,
				Interval:             baseModel.Interval,
			},
		},
		{
			name: "with empty success codes",
			input: &models.Endpoint{
				ID:                   baseModel.ID,
				ServiceName:          baseModel.ServiceName,
				URL:                  baseModel.URL,
				SuccessCodes:         []int{},
				NotificationServices: baseModel.NotificationServices,
				Interval:             baseModel.Interval,
			},
			expected: Endpoint{
				ID:                   baseModel.ID,
				ServiceName:          baseModel.ServiceName,
				URL:                  baseModel.URL,
				SuccessCodes:         []string{},
				NotificationServices: baseModel.NotificationServices,
				Interval:             baseModel.Interval,
			},
		},
		{
			name: "with duplicate success codes",
			input: &models.Endpoint{
				ID:                   baseModel.ID,
				ServiceName:          baseModel.ServiceName,
				URL:                  baseModel.URL,
				SuccessCodes:         []int{200, 200, 201, 201, 202},
				NotificationServices: baseModel.NotificationServices,
				Interval:             baseModel.Interval,
			},
			expected: Endpoint{
				ID:                   baseModel.ID,
				ServiceName:          baseModel.ServiceName,
				URL:                  baseModel.URL,
				SuccessCodes:         []string{"200-202"},
				NotificationServices: baseModel.NotificationServices,
				Interval:             baseModel.Interval,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromServiceEndpoint(tt.input)
			assert.Equal(t, tt.expected, result, "FromServiceEndpoint() returned unexpected result")
		})
	}
}

func Benchmark_CodesRanges(b *testing.B) {
	codes := generateSequence(100, 299)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = codesRanges(codes)
	}
}

func Benchmark_FromServiceEndpoint(b *testing.B) {
	model := &models.Endpoint{
		SuccessCodes: generateSequence(100, 299),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FromServiceEndpoint(model)
	}
}

// Helper function to generate a sequence of integers
func generateSequence(start, end int) []int {
	result := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}
