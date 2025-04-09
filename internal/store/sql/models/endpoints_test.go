package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	srv_models "github.com/vishenosik/CherryWatch/internal/services/models"
)

func TestSuccessCodesBatch(t *testing.T) {
	tests := []struct {
		name     string
		input    srv_models.Endpoints
		expected SuccessCodes
	}{
		{
			name:     "empty input",
			input:    srv_models.Endpoints{},
			expected: SuccessCodes{},
		},
		{
			name: "single endpoint with no success codes",
			input: srv_models.Endpoints{
				{
					ID:           "1",
					ServiceName:  "service1",
					SuccessCodes: []int{},
				},
			},
			expected: SuccessCodes{},
		},
		{
			name: "single endpoint with single success code",
			input: srv_models.Endpoints{
				{
					ID:           "1",
					ServiceName:  "service1",
					SuccessCodes: []int{200},
				},
			},
			expected: SuccessCodes{
				{ID: "1", Code: 200},
			},
		},
		{
			name: "single endpoint with multiple success codes",
			input: srv_models.Endpoints{
				{
					ID:           "1",
					ServiceName:  "service1",
					SuccessCodes: []int{200, 201, 204},
				},
			},
			expected: SuccessCodes{
				{ID: "1", Code: 200},
				{ID: "1", Code: 201},
				{ID: "1", Code: 204},
			},
		},
		{
			name: "multiple models.Endpoints with success codes",
			input: srv_models.Endpoints{
				{
					ID:           "1",
					ServiceName:  "service1",
					SuccessCodes: []int{200, 201},
				},
				{
					ID:           "2",
					ServiceName:  "service2",
					SuccessCodes: []int{204, 301},
				},
			},
			expected: SuccessCodes{
				{ID: "1", Code: 200},
				{ID: "1", Code: 201},
				{ID: "2", Code: 204},
				{ID: "2", Code: 301},
			},
		},
		{
			name: "mix of models.Endpoints with and without success codes",
			input: srv_models.Endpoints{
				{
					ID:           "1",
					ServiceName:  "service1",
					SuccessCodes: []int{200},
				},
				{
					ID:           "2",
					ServiceName:  "service2",
					SuccessCodes: []int{},
				},
				{
					ID:           "3",
					ServiceName:  "service3",
					SuccessCodes: []int{204, 304},
				},
			},
			expected: SuccessCodes{
				{ID: "1", Code: 200},
				{ID: "3", Code: 204},
				{ID: "3", Code: 304},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SuccessCodesBatch(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func BenchmarkSuccessCodesBatch(b *testing.B) {
	benchmarks := []struct {
		name  string
		input srv_models.Endpoints
	}{
		{
			name:  "empty",
			input: srv_models.Endpoints{},
		},
		{
			name: "small",
			input: srv_models.Endpoints{
				{
					ID:           "1",
					SuccessCodes: []int{200, 201},
				},
				{
					ID:           "2",
					SuccessCodes: []int{204},
				},
			},
		},
		{
			name:  "medium",
			input: generateEndpoints(100, 5),
		},
		{
			name:  "large",
			input: generateEndpoints(1000, 10),
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = SuccessCodesBatch(bm.input)
			}
		})
	}
}

// generateEndpoints helper function to create test data
func generateEndpoints(count, codesPerEndpoint int) srv_models.Endpoints {
	endpoints := make(srv_models.Endpoints, count)
	for i := 0; i < count; i++ {
		codes := make([]int, codesPerEndpoint)
		for j := 0; j < codesPerEndpoint; j++ {
			codes[j] = 200 + j
		}
		endpoints[i] = &srv_models.Endpoint{
			ID:           string(rune('a' + i%26)),
			ServiceName:  "service",
			SuccessCodes: codes,
		}
	}
	return endpoints
}

// Unit tests
func TestConverters(t *testing.T) {
	t.Run("TestConvertSliceToStructuredEndpoints", func(t *testing.T) {
		endpoints := srv_models.Endpoints{
			{
				ID:                   "123",
				ServiceName:          "service1",
				URL:                  "http://service1.com",
				SuccessCodes:         []int{200, 201},
				NotificationServices: []string{"slack", "email"},
				Interval:             30 * time.Second,
			},
			{
				ID:                   "456",
				ServiceName:          "service2",
				URL:                  "http://service2.com",
				SuccessCodes:         []int{200, 204},
				NotificationServices: []string{"pagerduty"},
				Interval:             60 * time.Second,
			},
		}

		structured := FromServiceEndpoints(endpoints)

		assert.NotNil(t, structured)
		assert.Len(t, structured.Infos, 2)
		assert.Len(t, structured.SuccessCodes, 4)
		assert.Len(t, structured.NotificationServices, 3)

		// Verify the first endpoint
		assert.Equal(t, "123", structured.Infos[0].ID)
		assert.Equal(t, "service1", structured.Infos[0].ServiceName)
		assert.Equal(t, "http://service1.com", structured.Infos[0].URL)
		assert.Equal(t, 30*time.Second, structured.Infos[0].Interval)

		// Verify success codes for first endpoint
		var codesForFirst []int
		for _, code := range structured.SuccessCodes {
			if code.ID == "123" {
				codesForFirst = append(codesForFirst, code.Code)
			}
		}
		assert.ElementsMatch(t, []int{200, 201}, codesForFirst)

		// Verify notification services for first endpoint
		var servicesForFirst []string
		for _, service := range structured.NotificationServices {
			if service.ID == "123" {
				servicesForFirst = append(servicesForFirst, service.ServiceName)
			}
		}
		assert.ElementsMatch(t, []string{"slack", "email"}, servicesForFirst)
	})

	t.Run("TestConvertSliceToStructuredEndpoints_EmptyInput", func(t *testing.T) {
		structured := FromServiceEndpoints(nil)
		assert.Nil(t, structured)

		structured = FromServiceEndpoints(srv_models.Endpoints{})
		assert.Nil(t, structured)
	})

	t.Run("TestConvertStructuredToSliceEndpoints", func(t *testing.T) {
		structured := &Endpoints{
			Infos: []*Info{
				{
					ID:          "123",
					ServiceName: "service1",
					URL:         "http://service1.com",
					Interval:    30 * time.Second,
				},
				{
					ID:          "456",
					ServiceName: "service2",
					URL:         "http://service2.com",
					Interval:    60 * time.Second,
				},
			},
			SuccessCodes: []*SuccessCode{
				{ID: "123", Code: 200},
				{ID: "123", Code: 201},
				{ID: "456", Code: 200},
				{ID: "456", Code: 204},
			},
			NotificationServices: []*NotificationService{
				{ID: "123", ServiceName: "slack"},
				{ID: "123", ServiceName: "email"},
				{ID: "456", ServiceName: "pagerduty"},
			},
		}

		endpoints := ToServiceEndpoints(structured)

		assert.NotNil(t, endpoints)
		assert.Len(t, endpoints, 2)

		// Verify first endpoint
		assert.Equal(t, "123", endpoints[0].ID)
		assert.Equal(t, "service1", endpoints[0].ServiceName)
		assert.Equal(t, "http://service1.com", endpoints[0].URL)
		assert.Equal(t, 30*time.Second, endpoints[0].Interval)
		assert.ElementsMatch(t, []int{200, 201}, endpoints[0].SuccessCodes)
		assert.ElementsMatch(t, []string{"slack", "email"}, endpoints[0].NotificationServices)

		// Verify second endpoint
		assert.Equal(t, "456", endpoints[1].ID)
		assert.Equal(t, "service2", endpoints[1].ServiceName)
		assert.Equal(t, "http://service2.com", endpoints[1].URL)
		assert.Equal(t, 60*time.Second, endpoints[1].Interval)
		assert.ElementsMatch(t, []int{200, 204}, endpoints[1].SuccessCodes)
		assert.ElementsMatch(t, []string{"pagerduty"}, endpoints[1].NotificationServices)
	})

	t.Run("TestConvertStructuredToSliceEndpoints_NilInput", func(t *testing.T) {
		endpoints := ToServiceEndpoints(nil)
		assert.Nil(t, endpoints)

		endpoints = ToServiceEndpoints(&Endpoints{Infos: []*Info{}})
		assert.Nil(t, endpoints)
	})

	t.Run("TestRoundTripConversion", func(t *testing.T) {
		original := srv_models.Endpoints{
			{
				ID:                   "789",
				ServiceName:          "service3",
				URL:                  "http://service3.com",
				SuccessCodes:         []int{200},
				NotificationServices: []string{"sms"},
				Interval:             15 * time.Second,
			},
		}

		// Convert to structured and back
		structured := FromServiceEndpoints(original)
		convertedBack := ToServiceEndpoints(structured)

		assert.Equal(t, original, convertedBack)
	})
}
