package models

import "time"

type Endpoint struct {
	ID                   string `db:"id"`
	ServiceName          string `db:"service_name"`
	URL                  string `db:"url"`
	SuccessCodes         []int
	NotificationServices []string
	Interval             time.Duration `db:"interval"`
}
