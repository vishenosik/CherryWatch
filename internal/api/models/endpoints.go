package models

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/vishenosik/CherryWatch/internal/services/models"
	devCol "github.com/vishenosik/CherryWatch/pkg/collections"
	"github.com/vishenosik/web-tools/collections"
)

type Endpoint struct {
	// Endpoint identifier (uuid4 only)
	ID string `json:"id"`
	// Name of checked service (ascii symbols only)
	ServiceName string `json:"service_name"`
	// URL string to trigger during checks
	URL string `json:"url"`
	// HTTP codes & code ranges which are considered successful
	// (codes should be in [100,599], ranges are strings like "200-345")
	SuccessCodes []string `json:"success_codes,omitempty"`
	// Services used to notify about check failure
	NotificationServices []string `json:"notification_services,omitempty"`
	// Time interval between checks
	Interval time.Duration `json:"time_interval"`
}

type Endpoints = []Endpoint

func ToServiceEndpoints(edps Endpoints) models.Endpoints {
	return devCol.ConvertFunc(edps, ToServiceEndpoint)
}

func ToServiceEndpoint(endpoint Endpoint) *models.Endpoint {
	ranges, err := parseRanges(endpoint.SuccessCodes)
	if err != nil {
		// TODO: Probably want to warn about error
	}
	return &models.Endpoint{
		ID:                   endpoint.ID,
		ServiceName:          endpoint.ServiceName,
		URL:                  endpoint.URL,
		SuccessCodes:         ranges,
		NotificationServices: endpoint.NotificationServices,
		Interval:             endpoint.Interval,
	}
}

func FromServiceEndpoints(edps models.Endpoints) Endpoints {
	return devCol.ConvertFunc(edps, FromServiceEndpoint)
}

func FromServiceEndpoint(endpoint *models.Endpoint) Endpoint {
	ranges := codesRanges(endpoint.SuccessCodes)
	return Endpoint{
		ID:                   endpoint.ID,
		ServiceName:          endpoint.ServiceName,
		URL:                  endpoint.URL,
		SuccessCodes:         ranges,
		NotificationServices: endpoint.NotificationServices,
		Interval:             endpoint.Interval,
	}
}

// IntsToRangeStrings converts []int to []string with ranges where possible
func codesRanges(codes []int) []string {
	if len(codes) == 0 {
		return []string{}
	}

	codes = collections.Unique(codes)
	// Sort the codes first
	sort.Ints(codes)

	var result []string
	start := codes[0]
	prev := codes[0]

	for i := 1; i < len(codes); i++ {
		current := codes[i]
		if current == prev+1 {
			prev = current
		} else {
			if start == prev {
				result = append(result, strconv.Itoa(start))
			} else {
				result = append(result, fmt.Sprintf("%d-%d", start, prev))
			}
			start = current
			prev = current
		}
	}

	// Add the last range or number
	if start == prev {
		result = append(result, strconv.Itoa(start))
	} else {
		result = append(result, fmt.Sprintf("%d-%d", start, prev))
	}

	return result
}

// ParseRanges converts slice of range strings to sorted unique []int
func parseRanges(rangeStrs []string) ([]int, error) {
	result := make([]int, 0, 0)
	for _, rangeStr := range rangeStrs {
		codes, err := parseRange(rangeStr)
		if err != nil {
			return nil, err
		}
		result = append(result, codes...)
	}
	return collections.Unique(result), nil
}

// ParseRange converts a string range like "200-345" to []int{200, 201, ..., 345}
func parseRange(rangeStr string) ([]int, error) {

	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		code, err := strconv.Atoi(rangeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid range format, expected 'start-end' or 'code'")
		}
		return []int{code}, nil
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid start value: %w", err)
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid end value: %w", err)
	}

	if start > end {
		return nil, fmt.Errorf("start cannot be greater than end")
	}

	result := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		result = append(result, i)
	}

	return result, nil
}
