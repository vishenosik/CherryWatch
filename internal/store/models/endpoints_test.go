package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessCodesBatch(t *testing.T) {
	tests := []struct {
		name     string
		input    Endpoints
		expected []SuccessCode
	}{
		{
			name:     "empty input",
			input:    Endpoints{},
			expected: []SuccessCode{},
		},
		{
			name: "single endpoint with no success codes",
			input: Endpoints{
				{
					ID:           "1",
					ServiceName:  "service1",
					SuccessCodes: []int{},
				},
			},
			expected: []SuccessCode{},
		},
		{
			name: "single endpoint with single success code",
			input: Endpoints{
				{
					ID:           "1",
					ServiceName:  "service1",
					SuccessCodes: []int{200},
				},
			},
			expected: []SuccessCode{
				{ID: "1", Code: 200},
			},
		},
		{
			name: "single endpoint with multiple success codes",
			input: Endpoints{
				{
					ID:           "1",
					ServiceName:  "service1",
					SuccessCodes: []int{200, 201, 204},
				},
			},
			expected: []SuccessCode{
				{ID: "1", Code: 200},
				{ID: "1", Code: 201},
				{ID: "1", Code: 204},
			},
		},
		{
			name: "multiple endpoints with success codes",
			input: Endpoints{
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
			expected: []SuccessCode{
				{ID: "1", Code: 200},
				{ID: "1", Code: 201},
				{ID: "2", Code: 204},
				{ID: "2", Code: 301},
			},
		},
		{
			name: "mix of endpoints with and without success codes",
			input: Endpoints{
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
			expected: []SuccessCode{
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
		input Endpoints
	}{
		{
			name:  "empty",
			input: Endpoints{},
		},
		{
			name: "small",
			input: Endpoints{
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
func generateEndpoints(count, codesPerEndpoint int) Endpoints {
	endpoints := make(Endpoints, count)
	for i := 0; i < count; i++ {
		codes := make([]int, codesPerEndpoint)
		for j := 0; j < codesPerEndpoint; j++ {
			codes[j] = 200 + j
		}
		endpoints[i] = Endpoint{
			ID:           string(rune('a' + i%26)),
			ServiceName:  "service",
			SuccessCodes: codes,
		}
	}
	return endpoints
}
