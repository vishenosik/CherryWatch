package postgresql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Endpoint struct {
	ID                   string
	ServiceName          string
	URL                  string
	SuccessCodes         []int
	NotificationServices []string
	Interval             time.Duration
}

type EndpointStorage struct {
	db *sql.DB
}

func NewEndpointStorage(db *sql.DB) *EndpointStorage {
	return &EndpointStorage{db: db}
}

// CreateEndpoint inserts a new endpoint into the database
func (s *EndpointStorage) CreateEndpoint(e *Endpoint) error {
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

	// Batch insert success codes
	if len(e.SuccessCodes) > 0 {
		valueStrings := make([]string, 0, len(e.SuccessCodes))
		valueArgs := make([]interface{}, 0, len(e.SuccessCodes)*2)
		for i, code := range e.SuccessCodes {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			valueArgs = append(valueArgs, id, code)
		}
		stmt := fmt.Sprintf(`
			INSERT INTO endpoint_success_codes (endpoint_id, code)
			VALUES %s
		`, strings.Join(valueStrings, ","))
		_, err = tx.Exec(stmt, valueArgs...)
		if err != nil {
			return fmt.Errorf("failed to insert success codes: %v", err)
		}
	}

	// Batch insert notification services
	if len(e.NotificationServices) > 0 {
		valueStrings := make([]string, 0, len(e.NotificationServices))
		valueArgs := make([]interface{}, 0, len(e.NotificationServices)*2)
		for i, service := range e.NotificationServices {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			valueArgs = append(valueArgs, id, service)
		}
		stmt := fmt.Sprintf(`
			INSERT INTO endpoint_notification_services (endpoint_id, service_name)
			VALUES %s
		`, strings.Join(valueStrings, ","))
		_, err = tx.Exec(stmt, valueArgs...)
		if err != nil {
			return fmt.Errorf("failed to insert notification services: %v", err)
		}
	}

	return tx.Commit()
}

// GetEndpoint retrieves an endpoint by ID using a single query
func (s *EndpointStorage) GetEndpoint(id string) (*Endpoint, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %v", err)
	}

	var e Endpoint
	var intervalSeconds int
	// var successCodes []sql.NullInt64
	// var notificationServices []sql.NullString

	// Single query with LEFT JOINs and array aggregation
	row := s.db.QueryRow(`
		SELECT 
			e.id,
			e.service_name,
			e.url,
			e.interval_seconds,
			ARRAY(
				SELECT code 
				FROM endpoint_success_codes 
				WHERE endpoint_id = e.id
			) AS success_codes,
			ARRAY(
				SELECT service_name 
				FROM endpoint_notification_services 
				WHERE endpoint_id = e.id
			) AS notification_services
		FROM 
			endpoints e
		WHERE 
			e.id = $1
	`, uid)

	var codes []int64
	var services []string
	err = row.Scan(
		&e.ID,
		&e.ServiceName,
		&e.URL,
		&intervalSeconds,
		pq.Array(&codes),
		pq.Array(&services),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("endpoint not found")
		}
		return nil, fmt.Errorf("failed to get endpoint: %v", err)
	}

	e.Interval = time.Duration(intervalSeconds) * time.Second

	// Convert success codes
	e.SuccessCodes = make([]int, len(codes))
	for i, code := range codes {
		e.SuccessCodes[i] = int(code)
	}

	// Convert notification services
	e.NotificationServices = services

	return &e, nil
}

// GetAllEndpoints retrieves all endpoints using a single query
func (s *EndpointStorage) GetAllEndpoints() ([]Endpoint, error) {
	// Single query with array aggregation
	rows, err := s.db.Query(`
		SELECT 
			e.id,
			e.service_name,
			e.url,
			e.interval_seconds,
			ARRAY(
				SELECT code 
				FROM endpoint_success_codes 
				WHERE endpoint_id = e.id
			) AS success_codes,
			ARRAY(
				SELECT service_name 
				FROM endpoint_notification_services 
				WHERE endpoint_id = e.id
			) AS notification_services
		FROM 
			endpoints e
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints: %v", err)
	}
	defer rows.Close()

	var endpoints []Endpoint
	for rows.Next() {
		var e Endpoint
		var intervalSeconds int
		var codes []int64
		var services []string

		err := rows.Scan(
			&e.ID,
			&e.ServiceName,
			&e.URL,
			&intervalSeconds,
			pq.Array(&codes),
			pq.Array(&services),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan endpoint: %v", err)
		}

		e.Interval = time.Duration(intervalSeconds) * time.Second
		e.SuccessCodes = make([]int, len(codes))
		for i, code := range codes {
			e.SuccessCodes[i] = int(code)
		}
		e.NotificationServices = services

		endpoints = append(endpoints, e)
	}

	return endpoints, nil
}

// UpdateEndpoint modifies an existing endpoint
func (s *EndpointStorage) UpdateEndpoint(e *Endpoint) error {
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

	// Delete and re-insert success codes in one operation
	_, err = tx.Exec(`
		WITH deleted AS (
			DELETE FROM endpoint_success_codes 
			WHERE endpoint_id = $1
			RETURNING 1
		)
		INSERT INTO endpoint_success_codes (endpoint_id, code)
		SELECT $1, unnest($2::int[])
	`, id, pq.Array(e.SuccessCodes))
	if err != nil {
		return fmt.Errorf("failed to update success codes: %v", err)
	}

	// Delete and re-insert notification services in one operation
	_, err = tx.Exec(`
		WITH deleted AS (
			DELETE FROM endpoint_notification_services 
			WHERE endpoint_id = $1
			RETURNING 1
		)
		INSERT INTO endpoint_notification_services (endpoint_id, service_name)
		SELECT $1, unnest($2::text[])
	`, id, pq.Array(e.NotificationServices))
	if err != nil {
		return fmt.Errorf("failed to update notification services: %v", err)
	}

	return tx.Commit()
}

// DeleteEndpoint removes an endpoint from the database
func (s *EndpointStorage) DeleteEndpoint(id string) error {
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
