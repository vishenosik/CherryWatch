package endpoints

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Endpoint struct {
	ID                   string
	ServiceName          string
	URL                  string
	SuccessCodes         []int
	NotificationServices []string
	Interval             time.Duration
}

type endpointStorage struct {
	db *sql.DB
}

func NewEndpointStorage(db *sql.DB) *endpointStorage {
	return &endpointStorage{db: db}
}

// CreateEndpoint inserts a new endpoint into the database
func (s *endpointStorage) CreateEndpoint(e *Endpoint) error {
	id, err := uuid.Parse(e.ID)
	if err != nil {
		return fmt.Errorf("invalid UUID: %v", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert main endpoint record
	_, err = tx.Exec(`
		INSERT INTO endpoints (id, service_name, url, interval_seconds)
		VALUES ($1, $2, $3, $4)
	`, id, e.ServiceName, e.URL, int(e.Interval.Seconds()))
	if err != nil {
		return fmt.Errorf("failed to insert endpoint: %v", err)
	}

	// Insert success codes
	for _, code := range e.SuccessCodes {
		_, err = tx.Exec(`
			INSERT INTO endpoint_success_codes (endpoint_id, code)
			VALUES ($1, $2)
		`, id, code)
		if err != nil {
			return fmt.Errorf("failed to insert success code: %v", err)
		}
	}

	// Insert notification services
	for _, service := range e.NotificationServices {
		_, err = tx.Exec(`
			INSERT INTO endpoint_notification_services (endpoint_id, service_name)
			VALUES ($1, $2)
		`, id, service)
		if err != nil {
			return fmt.Errorf("failed to insert notification service: %v", err)
		}
	}

	return tx.Commit()
}

// GetEndpoint retrieves an endpoint by ID
func (s *endpointStorage) GetEndpoint(id string) (*Endpoint, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %v", err)
	}

	var e Endpoint
	var intervalSeconds int

	// Get basic endpoint info
	err = s.db.QueryRow(`
		SELECT id, service_name, url, interval_seconds
		FROM endpoints
		WHERE id = $1
	`, uid).Scan(&e.ID, &e.ServiceName, &e.URL, &intervalSeconds)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("endpoint not found")
		}
		return nil, fmt.Errorf("failed to get endpoint: %v", err)
	}
	e.Interval = time.Duration(intervalSeconds) * time.Second

	// Get success codes
	rows, err := s.db.Query(`
		SELECT code
		FROM endpoint_success_codes
		WHERE endpoint_id = $1
	`, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get success codes: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var code int
		if err := rows.Scan(&code); err != nil {
			return nil, fmt.Errorf("failed to scan success code: %v", err)
		}
		e.SuccessCodes = append(e.SuccessCodes, code)
	}

	// Get notification services
	rows, err = s.db.Query(`
		SELECT service_name
		FROM endpoint_notification_services
		WHERE endpoint_id = $1
	`, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification services: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var service string
		if err := rows.Scan(&service); err != nil {
			return nil, fmt.Errorf("failed to scan notification service: %v", err)
		}
		e.NotificationServices = append(e.NotificationServices, service)
	}

	return &e, nil
}

// GetAllEndpoints retrieves all endpoints from the database
func (s *endpointStorage) GetAllEndpoints() ([]Endpoint, error) {
	var endpoints []Endpoint

	// Get all endpoint IDs first
	rows, err := s.db.Query("SELECT id FROM endpoints")
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan endpoint ID: %v", err)
		}

		endpoint, err := s.GetEndpoint(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get endpoint %s: %v", id, err)
		}
		endpoints = append(endpoints, *endpoint)
	}

	return endpoints, nil
}

// UpdateEndpoint modifies an existing endpoint
func (s *endpointStorage) UpdateEndpoint(e *Endpoint) error {
	id, err := uuid.Parse(e.ID)
	if err != nil {
		return fmt.Errorf("invalid UUID: %v", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update main endpoint record
	_, err = tx.Exec(`
		UPDATE endpoints
		SET service_name = $1, url = $2, interval_seconds = $3
		WHERE id = $4
	`, e.ServiceName, e.URL, int(e.Interval.Seconds()), id)
	if err != nil {
		return fmt.Errorf("failed to update endpoint: %v", err)
	}

	// Delete existing success codes
	_, err = tx.Exec("DELETE FROM endpoint_success_codes WHERE endpoint_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete old success codes: %v", err)
	}

	// Insert new success codes
	for _, code := range e.SuccessCodes {
		_, err = tx.Exec(`
			INSERT INTO endpoint_success_codes (endpoint_id, code)
			VALUES ($1, $2)
		`, id, code)
		if err != nil {
			return fmt.Errorf("failed to insert success code: %v", err)
		}
	}

	// Delete existing notification services
	_, err = tx.Exec("DELETE FROM endpoint_notification_services WHERE endpoint_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete old notification services: %v", err)
	}

	// Insert new notification services
	for _, service := range e.NotificationServices {
		_, err = tx.Exec(`
			INSERT INTO endpoint_notification_services (endpoint_id, service_name)
			VALUES ($1, $2)
		`, id, service)
		if err != nil {
			return fmt.Errorf("failed to insert notification service: %v", err)
		}
	}

	return tx.Commit()
}

// DeleteEndpoint removes an endpoint from the database
func (s *endpointStorage) DeleteEndpoint(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid UUID: %v", err)
	}

	_, err = s.db.Exec("DELETE FROM endpoints WHERE id = $1", uid)
	if err != nil {
		return fmt.Errorf("failed to delete endpoint: %v", err)
	}
	return nil
}

// JSON representation for API responses
func (e *Endpoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID                   string   `json:"id"`
		ServiceName          string   `json:"service_name"`
		URL                  string   `json:"url"`
		SuccessCodes         []int    `json:"success_codes"`
		NotificationServices []string `json:"notification_services"`
		Interval             string   `json:"interval"`
	}{
		ID:                   e.ID,
		ServiceName:          e.ServiceName,
		URL:                  e.URL,
		SuccessCodes:         e.SuccessCodes,
		NotificationServices: e.NotificationServices,
		Interval:             e.Interval.String(),
	})
}
