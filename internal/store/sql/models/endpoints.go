package models

import (
	"time"

	srv_models "github.com/vishenosik/CherryWatch/internal/services/models"
)

type Endpoint struct {
	ID          string        `db:"id"`
	ServiceName string        `db:"service_name"`
	URL         string        `db:"url"`
	Interval    time.Duration `db:"interval"`
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

type StoreEndpoints struct {
	Endpoints            Endpoints
	SuccessCodes         SuccessCodes
	NotificationServices NotificationServices
}

func SuccessCodesBatch(edps srv_models.Endpoints) SuccessCodes {

	successCodes := make(SuccessCodes, 0, len(edps))

	for _, edp := range edps {
		for _, sc := range edp.SuccessCodes {
			successCodes = append(successCodes, &SuccessCode{
				ID:   edp.ID,
				Code: sc,
			})
		}
	}

	return successCodes
}

func NotificationServicesBatch(edps srv_models.Endpoints) NotificationServices {

	successCodes := make(NotificationServices, 0, len(edps))

	for _, edp := range edps {
		for _, sc := range edp.NotificationServices {
			successCodes = append(successCodes, &NotificationService{
				ID:          edp.ID,
				ServiceName: sc,
			})
		}
	}

	return successCodes
}

// Convert []*Endpoint to structured Endpoints model
func FromServiceEndpoints(endpoints srv_models.Endpoints) *StoreEndpoints {
	if len(endpoints) == 0 {
		return nil
	}

	infos := make(Endpoints, 0, len(endpoints))
	successCodes := make(SuccessCodes, 0)
	notificationServices := make(NotificationServices, 0)

	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		infos = append(infos, &Endpoint{
			ID:          endpoint.ID,
			ServiceName: endpoint.ServiceName,
			URL:         endpoint.URL,
			Interval:    endpoint.Interval,
		})

		for _, code := range endpoint.SuccessCodes {
			successCodes = append(successCodes, &SuccessCode{
				ID:   endpoint.ID,
				Code: code,
			})
		}

		for _, service := range endpoint.NotificationServices {
			notificationServices = append(notificationServices, &NotificationService{
				ID:          endpoint.ID,
				ServiceName: service,
			})
		}
	}

	return &StoreEndpoints{
		Endpoints:            infos,
		SuccessCodes:         successCodes,
		NotificationServices: notificationServices,
	}
}

// Convert structured Endpoints back to []*Endpoint
func ToServiceEndpoints(structured *StoreEndpoints) srv_models.Endpoints {
	if structured == nil || len(structured.Endpoints) == 0 {
		return nil
	}

	// Create a map for faster lookup of success codes and notification services
	successCodesMap := make(map[string][]int)
	for _, code := range structured.SuccessCodes {
		successCodesMap[code.ID] = append(successCodesMap[code.ID], code.Code)
	}

	notificationServicesMap := make(map[string][]string)
	for _, service := range structured.NotificationServices {
		notificationServicesMap[service.ID] = append(notificationServicesMap[service.ID], service.ServiceName)
	}

	endpoints := make(srv_models.Endpoints, 0, len(structured.Endpoints))
	for _, info := range structured.Endpoints {
		endpoint := &srv_models.Endpoint{
			ID:                   info.ID,
			ServiceName:          info.ServiceName,
			URL:                  info.URL,
			Interval:             info.Interval,
			SuccessCodes:         successCodesMap[info.ID],
			NotificationServices: notificationServicesMap[info.ID],
		}
		endpoints = append(endpoints, endpoint)
	}

	return endpoints
}
