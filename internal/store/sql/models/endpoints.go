package models

import (
	"time"

	srv_models "github.com/vishenosik/CherryWatch/internal/services/models"
)

type Endpoint struct {
	ID                   string        `db:"id"`
	ServiceName          string        `db:"service_name"`
	URL                  string        `db:"url"`
	Interval             time.Duration `db:"interval"`
	SuccessCodes         SuccessCodes
	NotificationServices NotificationServices
}

type Endpoints = []*Endpoint

type SuccessCode struct {
	ID   string `db:"endpoint_id"`
	Code int    `db:"code"`
}

type SuccessCodes = []*SuccessCode

type NotificationService struct {
	ID          string `db:"endpoint_id"`
	ServiceName string `db:"service_name"`
}

type NotificationServices = []*NotificationService

// Convert []*Endpoint to structured Endpoints model
func FromServiceEndpoints(edps srv_models.Endpoints) Endpoints {

	if len(edps) == 0 {
		return nil
	}

	converted := make(Endpoints, 0, len(edps))

	for _, edp := range edps {

		if edp == nil {
			continue
		}

		successCodes := make(SuccessCodes, 0, len(edp.SuccessCodes))

		for _, code := range edp.SuccessCodes {
			successCodes = append(successCodes, &SuccessCode{
				ID:   edp.ID,
				Code: code,
			})
		}

		notificationServices := make(NotificationServices, 0, len(edp.NotificationServices))
		for _, service := range edp.NotificationServices {
			notificationServices = append(notificationServices, &NotificationService{
				ID:          edp.ID,
				ServiceName: service,
			})
		}

		converted = append(converted, &Endpoint{
			ID:                   edp.ID,
			ServiceName:          edp.ServiceName,
			URL:                  edp.URL,
			Interval:             edp.Interval,
			SuccessCodes:         successCodes,
			NotificationServices: notificationServices,
		})

	}

	return converted
}

// Convert structured Endpoints back to []*Endpoint
func ToServiceEndpoints(edps Endpoints) srv_models.Endpoints {

	if len(edps) == 0 {
		return nil
	}

	converted := make(srv_models.Endpoints, 0, len(edps))
	for _, edp := range edps {

		srv := &srv_models.Endpoint{
			ID:                   edp.ID,
			ServiceName:          edp.ServiceName,
			URL:                  edp.URL,
			Interval:             edp.Interval,
			SuccessCodes:         make([]int, 0, len(edp.SuccessCodes)),
			NotificationServices: make([]string, 0, len(edp.NotificationServices)),
		}

		for _, code := range edp.SuccessCodes {
			srv.SuccessCodes = append(srv.SuccessCodes, code.Code)
		}

		for _, service := range edp.NotificationServices {
			srv.NotificationServices = append(srv.NotificationServices, service.ServiceName)
		}

		converted = append(converted, srv)
	}

	return converted
}
