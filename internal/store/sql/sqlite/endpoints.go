package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	srv_models "github.com/vishenosik/CherryWatch/internal/services/models"
	"github.com/vishenosik/CherryWatch/internal/store/sql/models"
)

type endpoints struct {
	db *sqlx.DB
}

func newEndpoints(db *sqlx.DB) *endpoints {
	return &endpoints{db: db}
}

// CreateEndpoint inserts a new endpoint into the database
func (s *endpoints) CreateEndpoints(ctx context.Context, edps srv_models.Endpoints) error {
	return createEndpoints(ctx, s.db, models.FromServiceEndpoints(edps))
}

// CreateEndpoint inserts a new endpoint into the database
func createEndpoints(ctx context.Context, db *sqlx.DB, edps *models.StoreEndpoints) error {

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if edps == nil || len(edps.Endpoints) == 0 {
		return errors.New("endpoints can't be empty")
	}

	_, err = tx.NamedExecContext(ctx, `
		INSERT INTO endpoints (id, service_name, url, interval)
    	VALUES (:id, :service_name, :url, :interval)`,
		edps.Endpoints,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return errors.Wrap(models.ErrAlreadyExists, "failed to insert endpoints")
		}
		return errors.Wrap(err, "failed to insert endpoints")
	}

	// Batch insert success codes
	if len(edps.SuccessCodes) != 0 {
		_, err = tx.NamedExecContext(ctx, `
			INSERT INTO endpoint_success_codes (endpoint_id, code) 
			VALUES (:endpoint_id, :code)`,
			edps.SuccessCodes,
		)
		if err != nil {
			return errors.Wrap(err, "failed to insert success codes")
		}
	}
	// Batch insert notification services
	if len(edps.NotificationServices) != 0 {
		_, err = tx.NamedExecContext(ctx, `
	INSERT INTO endpoint_notification_services (endpoint_id, service_name)  
	VALUES (:endpoint_id, :service_name)`,
			edps.NotificationServices,
		)
		if err != nil {
			return errors.Wrap(err, "failed to insert notification services")
		}
	}

	return tx.Commit()
}

func (s *endpoints) GetEndpoints(ids ...string) (srv_models.Endpoints, error) {

	var (
		edps *models.StoreEndpoints
		err  error
	)

	if len(ids) != 0 {
		edps, err = getEndpoints(s.db, ids...)
	} else {
		edps, err = getAllEndpoints(s.db)
	}

	if err != nil {
		return nil, err
	}
	return models.ToServiceEndpoints(edps), nil
}

// GetAllEndpoints retrieves all endpoints using a single query
func getAllEndpoints(db *sqlx.DB) (*models.StoreEndpoints, error) {
	// Enable WAL mode for better concurrent read performance
	_, _ = db.Exec("PRAGMA journal_mode=WAL")

	// Single query with GROUP_CONCAT for SQLite
	rows, err := db.Queryx(`
		SELECT 
			e.id,
			e.service_name,
			e.url,
			e.interval,
			(
				SELECT GROUP_CONCAT(code, ',') 
				FROM endpoint_success_codes 
				WHERE endpoint_id = e.id
			) AS success_codes,
			(
				SELECT GROUP_CONCAT(service_name, ',') 
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

	return handleRows(rows)
}

// GetAllEndpoints retrieves endpoints with specified ids otherwise all endpoints retrieved
func getEndpoints(db *sqlx.DB, ids ...string) (*models.StoreEndpoints, error) {

	if len(ids) == 0 {
		return nil, errors.New("nil ids")
	}

	query, args, err := sqlx.In(`
		SELECT
			e.id,
			e.service_name,
			e.url,
			e.interval,
			(
				SELECT GROUP_CONCAT(code, ',')
				FROM endpoint_success_codes
				WHERE endpoint_id = e.id
			) AS success_codes,
			(
				SELECT GROUP_CONCAT(service_name, ',')
				FROM endpoint_notification_services
				WHERE endpoint_id = e.id
			) AS notification_services
		FROM
			endpoints e
		WHERE
			e.id IN (?)
	`,
		ids,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to configure IN statement: %v", err)
	}

	rows, err := db.Queryx(sqlx.Rebind(sqlx.DOLLAR, query), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints: %v", err)
	}
	defer rows.Close()

	return handleRows(rows)
}

func handleRows(rows *sqlx.Rows) (*models.StoreEndpoints, error) {

	type row struct {
		models.Endpoint
		Codes    sql.NullString `db:"success_codes"`
		Services sql.NullString `db:"notification_services"`
	}

	endpoints := new(models.StoreEndpoints)

	for rows.Next() {

		var r row
		err := rows.StructScan(&r)
		if err != nil {
			return nil, fmt.Errorf("failed to scan endpoint: %v", err)
		}

		endpoints.Endpoints = append(endpoints.Endpoints, &r.Endpoint)

		if r.Codes.Valid {

			// Parse success codes
			codes := strings.Split(r.Codes.String, ",")

			for _, codeStr := range codes {

				code, err := strconv.Atoi(codeStr)
				if err != nil {
					continue
				}

				endpoints.SuccessCodes = append(endpoints.SuccessCodes, &models.SuccessCode{
					ID:   r.ID,
					Code: code,
				})
			}
		}

		// Parse notification services
		if r.Services.Valid {

			srvs := strings.Split(r.Services.String, ",")

			for _, srv := range srvs {
				endpoints.NotificationServices = append(endpoints.NotificationServices, &models.NotificationService{
					ID:          r.ID,
					ServiceName: srv,
				})
			}
		}
	}

	return endpoints, nil
}

// // UpdateEndpoint modifies an existing endpoint
// func (s *endpoints) UpdateEndpoint(e *Endpoint) error {
// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	// Update main endpoint record
// 	_, err = tx.Exec(`
// 		UPDATE endpoints
// 		SET service_name = ?, url = ?, interval_seconds = ?
// 		WHERE id = ?
// 	`, e.ServiceName, e.URL, int(e.Interval.Seconds()), e.ID)
// 	if err != nil {
// 		return fmt.Errorf("failed to update endpoint: %v", err)
// 	}

// 	// Delete existing success codes
// 	_, err = tx.Exec("DELETE FROM endpoint_success_codes WHERE endpoint_id = ?", e.ID)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete old success codes: %v", err)
// 	}

// 	// Insert new success codes in batch
// 	if len(e.SuccessCodes) > 0 {
// 		stmt, err := tx.Prepare("INSERT INTO endpoint_success_codes (endpoint_id, code) VALUES (?, ?)")
// 		if err != nil {
// 			return fmt.Errorf("failed to prepare success codes statement: %v", err)
// 		}
// 		defer stmt.Close()

// 		for _, code := range e.SuccessCodes {
// 			if _, err := stmt.Exec(e.ID, code); err != nil {
// 				return fmt.Errorf("failed to insert success code: %v", err)
// 			}
// 		}
// 	}

// 	// Delete existing notification services
// 	_, err = tx.Exec("DELETE FROM endpoint_notification_services WHERE endpoint_id = ?", e.ID)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete old notification services: %v", err)
// 	}

// 	// Insert new notification services in batch
// 	if len(e.NotificationServices) > 0 {
// 		stmt, err := tx.Prepare("INSERT INTO endpoint_notification_services (endpoint_id, service_name) VALUES (?, ?)")
// 		if err != nil {
// 			return fmt.Errorf("failed to prepare notification services statement: %v", err)
// 		}
// 		defer stmt.Close()

// 		for _, service := range e.NotificationServices {
// 			if _, err := stmt.Exec(e.ID, service); err != nil {
// 				return fmt.Errorf("failed to insert notification service: %v", err)
// 			}
// 		}
// 	}

// 	return tx.Commit()
// }

// // UpdateEndpoints modifies multiple endpoints in a single transaction
// func (s *endpoints) UpdateEndpoints(endpoints []*Endpoint) error {
// 	if len(endpoints) == 0 {
// 		return nil
// 	}

// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	// Prepare all statements in advance
// 	updateEndpointStmt, err := tx.Prepare(`
//         UPDATE endpoints
//         SET service_name = ?, url = ?, interval_seconds = ?
//         WHERE id = ?
//     `)
// 	if err != nil {
// 		return fmt.Errorf("failed to prepare endpoint update statement: %v", err)
// 	}
// 	defer updateEndpointStmt.Close()

// 	deleteCodesStmt, err := tx.Prepare(`
//         DELETE FROM endpoint_success_codes
//         WHERE endpoint_id = ?
//     `)
// 	if err != nil {
// 		return fmt.Errorf("failed to prepare delete codes statement: %v", err)
// 	}
// 	defer deleteCodesStmt.Close()

// 	insertCodeStmt, err := tx.Prepare(`
//         INSERT INTO endpoint_success_codes (endpoint_id, code)
//         VALUES (?, ?)
//     `)
// 	if err != nil {
// 		return fmt.Errorf("failed to prepare insert code statement: %v", err)
// 	}
// 	defer insertCodeStmt.Close()

// 	deleteServicesStmt, err := tx.Prepare(`
//         DELETE FROM endpoint_notification_services
//         WHERE endpoint_id = ?
//     `)
// 	if err != nil {
// 		return fmt.Errorf("failed to prepare delete services statement: %v", err)
// 	}
// 	defer deleteServicesStmt.Close()

// 	insertServiceStmt, err := tx.Prepare(`
//         INSERT INTO endpoint_notification_services (endpoint_id, service_name)
//         VALUES (?, ?)
//     `)
// 	if err != nil {
// 		return fmt.Errorf("failed to prepare insert service statement: %v", err)
// 	}
// 	defer insertServiceStmt.Close()

// 	// Process each endpoint
// 	for _, e := range endpoints {
// 		// Update main endpoint record
// 		_, err = updateEndpointStmt.Exec(
// 			e.ServiceName,
// 			e.URL,
// 			int(e.Interval.Seconds()),
// 			e.ID,
// 		)
// 		if err != nil {
// 			return fmt.Errorf("failed to update endpoint %s: %v", e.ID, err)
// 		}

// 		// Update success codes
// 		_, err = deleteCodesStmt.Exec(e.ID)
// 		if err != nil {
// 			return fmt.Errorf("failed to delete old codes for endpoint %s: %v", e.ID, err)
// 		}

// 		for _, code := range e.SuccessCodes {
// 			_, err = insertCodeStmt.Exec(e.ID, code)
// 			if err != nil {
// 				return fmt.Errorf("failed to insert code %d for endpoint %s: %v", code, e.ID, err)
// 			}
// 		}

// 		// Update notification services
// 		_, err = deleteServicesStmt.Exec(e.ID)
// 		if err != nil {
// 			return fmt.Errorf("failed to delete old services for endpoint %s: %v", e.ID, err)
// 		}

// 		for _, service := range e.NotificationServices {
// 			_, err = insertServiceStmt.Exec(e.ID, service)
// 			if err != nil {
// 				return fmt.Errorf("failed to insert service %s for endpoint %s: %v", service, e.ID, err)
// 			}
// 		}
// 	}

// 	return tx.Commit()
// }

// // DeleteEndpoint removes an endpoint from the database
// func (s *endpoints) DeleteEndpoint(id string) error {
// 	_, err := s.db.Exec("DELETE FROM endpoints WHERE id = ?", id)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete endpoint: %v", err)
// 	}
// 	return nil
// }

// // JSON representation for API responses
// func (e *Endpoint) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(struct {
// 		ID                   string   `json:"id"`
// 		ServiceName          string   `json:"service_name"`
// 		URL                  string   `json:"url"`
// 		SuccessCodes         []int    `json:"success_codes"`
// 		NotificationServices []string `json:"notification_services"`
// 		Interval             string   `json:"interval"`
// 	}{
// 		ID:                   e.ID,
// 		ServiceName:          e.ServiceName,
// 		URL:                  e.URL,
// 		SuccessCodes:         e.SuccessCodes,
// 		NotificationServices: e.NotificationServices,
// 		Interval:             e.Interval.String(),
// 	})
// }
