package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	srv_models "github.com/vishenosik/CherryWatch/internal/services/models"
)

func TestFromServiceEndpoints(t *testing.T) {
	tests := []struct {
		name     string
		input    srv_models.Endpoints
		expected Endpoints
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty input",
			input:    srv_models.Endpoints{},
			expected: nil,
		},
		{
			name: "single endpoint",
			input: srv_models.Endpoints{
				&srv_models.Endpoint{
					ID:                   "1",
					ServiceName:          "test",
					URL:                  "http://example.com",
					Interval:             10 * time.Second,
					SuccessCodes:         []int{200, 201},
					NotificationServices: []string{"slack", "email"},
				},
			},
			expected: Endpoints{
				&Endpoint{
					ID:          "1",
					ServiceName: "test",
					URL:         "http://example.com",
					Interval:    10 * time.Second,
					SuccessCodes: SuccessCodes{
						&SuccessCode{ID: "1", Code: 200},
						&SuccessCode{ID: "1", Code: 201},
					},
					NotificationServices: NotificationServices{
						&NotificationService{ID: "1", ServiceName: "slack"},
						&NotificationService{ID: "1", ServiceName: "email"},
					},
				},
			},
		},
		{
			name: "multiple endpoints",
			input: srv_models.Endpoints{
				&srv_models.Endpoint{
					ID:                   "1",
					ServiceName:          "service1",
					URL:                  "http://service1.com",
					Interval:             5 * time.Second,
					SuccessCodes:         []int{200},
					NotificationServices: []string{"slack"},
				},
				&srv_models.Endpoint{
					ID:                   "2",
					ServiceName:          "service2",
					URL:                  "http://service2.com",
					Interval:             15 * time.Second,
					SuccessCodes:         []int{200, 204},
					NotificationServices: []string{"email"},
				},
			},
			expected: Endpoints{
				&Endpoint{
					ID:          "1",
					ServiceName: "service1",
					URL:         "http://service1.com",
					Interval:    5 * time.Second,
					SuccessCodes: SuccessCodes{
						&SuccessCode{ID: "1", Code: 200},
					},
					NotificationServices: NotificationServices{
						&NotificationService{ID: "1", ServiceName: "slack"},
					},
				},
				&Endpoint{
					ID:          "2",
					ServiceName: "service2",
					URL:         "http://service2.com",
					Interval:    15 * time.Second,
					SuccessCodes: SuccessCodes{
						&SuccessCode{ID: "2", Code: 200},
						&SuccessCode{ID: "2", Code: 204},
					},
					NotificationServices: NotificationServices{
						&NotificationService{ID: "2", ServiceName: "email"},
					},
				},
			},
		},
		{
			name: "endpoint with nil slices",
			input: srv_models.Endpoints{
				&srv_models.Endpoint{
					ID:          "1",
					ServiceName: "test",
					URL:         "http://example.com",
					Interval:    10 * time.Second,
				},
			},
			expected: Endpoints{
				&Endpoint{
					ID:                   "1",
					ServiceName:          "test",
					URL:                  "http://example.com",
					Interval:             10 * time.Second,
					SuccessCodes:         SuccessCodes{},
					NotificationServices: NotificationServices{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromServiceEndpoints(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToServiceEndpoints(t *testing.T) {
	tests := []struct {
		name     string
		input    Endpoints
		expected srv_models.Endpoints
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty input",
			input:    Endpoints{},
			expected: nil,
		},
		{
			name: "single endpoint",
			input: Endpoints{
				&Endpoint{
					ID:          "1",
					ServiceName: "test",
					URL:         "http://example.com",
					Interval:    10 * time.Second,
					SuccessCodes: SuccessCodes{
						&SuccessCode{ID: "1", Code: 200},
						&SuccessCode{ID: "1", Code: 201},
					},
					NotificationServices: NotificationServices{
						&NotificationService{ID: "1", ServiceName: "slack"},
						&NotificationService{ID: "1", ServiceName: "email"},
					},
				},
			},
			expected: srv_models.Endpoints{
				&srv_models.Endpoint{
					ID:                   "1",
					ServiceName:          "test",
					URL:                  "http://example.com",
					Interval:             10 * time.Second,
					SuccessCodes:         []int{200, 201},
					NotificationServices: []string{"slack", "email"},
				},
			},
		},
		{
			name: "endpoint with empty slices",
			input: Endpoints{
				&Endpoint{
					ID:                   "1",
					ServiceName:          "test",
					URL:                  "http://example.com",
					Interval:             10 * time.Second,
					SuccessCodes:         SuccessCodes{},
					NotificationServices: NotificationServices{},
				},
			},
			expected: srv_models.Endpoints{
				&srv_models.Endpoint{
					ID:                   "1",
					ServiceName:          "test",
					URL:                  "http://example.com",
					Interval:             10 * time.Second,
					SuccessCodes:         []int{},
					NotificationServices: []string{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToServiceEndpoints(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	original := srv_models.Endpoints{
		&srv_models.Endpoint{
			ID:                   "1",
			ServiceName:          "test",
			URL:                  "http://example.com",
			Interval:             10 * time.Second,
			SuccessCodes:         []int{200, 201},
			NotificationServices: []string{"slack", "email"},
		},
		&srv_models.Endpoint{
			ID:                   "2",
			ServiceName:          "another",
			URL:                  "http://another.com",
			Interval:             5 * time.Second,
			SuccessCodes:         []int{200},
			NotificationServices: []string{"sms"},
		},
	}

	// Convert to our model and back
	converted := FromServiceEndpoints(original)
	roundTrip := ToServiceEndpoints(converted)

	assert.Equal(t, original, roundTrip)
}

func generateTestEndpoints(count int) srv_models.Endpoints {
	endpoints := make(srv_models.Endpoints, count)
	for i := 0; i < count; i++ {
		endpoints[i] = &srv_models.Endpoint{
			ID:                   string(rune('a' + i)),
			ServiceName:          "service-" + string(rune('a'+i)),
			URL:                  "http://" + string(rune('a'+i)) + ".com",
			Interval:             time.Duration(i+1) * time.Second,
			SuccessCodes:         []int{200, 201, 204},
			NotificationServices: []string{"slack", "email", "sms"},
		}
	}
	return endpoints
}

func BenchmarkFromServiceEndpoints(b *testing.B) {
	testCases := []struct {
		name  string
		count int
	}{
		{"1 endpoint", 1},
		{"10 endpoints", 10},
		{"100 endpoints", 100},
		{"1000 endpoints", 1000},
	}

	for _, tc := range testCases {
		endpoints := generateTestEndpoints(tc.count)
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = FromServiceEndpoints(endpoints)
			}
		})
	}
}

func BenchmarkToServiceEndpoints(b *testing.B) {
	testCases := []struct {
		name  string
		count int
	}{
		{"1 endpoint", 1},
		{"10 endpoints", 10},
		{"100 endpoints", 100},
		{"1000 endpoints", 1000},
	}

	for _, tc := range testCases {
		serviceEndpoints := generateTestEndpoints(tc.count)
		endpoints := FromServiceEndpoints(serviceEndpoints)
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = ToServiceEndpoints(endpoints)
			}
		})
	}
}

func BenchmarkRoundTripConversion(b *testing.B) {
	testCases := []struct {
		name  string
		count int
	}{
		{"1 endpoint", 1},
		{"10 endpoints", 10},
		{"100 endpoints", 100},
	}

	for _, tc := range testCases {
		endpoints := generateTestEndpoints(tc.count)
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				converted := FromServiceEndpoints(endpoints)
				_ = ToServiceEndpoints(converted)
			}
		})
	}
}
