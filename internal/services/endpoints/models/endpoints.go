package models

import "time"

type Endpoint struct {
	// Endpoint identifier
	ID string
	// Name of checked service
	ServiceName string
	// Protocol to use when checking (http, https, ws, etc.)
	Protocol string
	// URL path to trigger
	URL string
	// HTTP codes which are considered successful
	SuccessCodes []uint16
	// Time interval between checks
	Interval time.Duration
	// Services used to notify about check failure
	NotificationServices []string
}
