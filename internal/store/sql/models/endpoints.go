package models

import (
	"time"

	srv_models "github.com/vishenosik/CherryWatch/internal/services/models"
)

type Info struct {
	ID          string        `db:"id"`
	ServiceName string        `db:"service_name"`
	URL         string        `db:"url"`
	Interval    time.Duration `db:"interval"`
}

type Infos = []*Info

type Endpoints struct {
	Infos                Infos
	SuccessCodes         SuccessCodes
	NotificationServices NotificationServices
}

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
func FromServiceEndpoints(endpoints srv_models.Endpoints) *Endpoints {
	if len(endpoints) == 0 {
		return nil
	}

	infos := make([]*Info, 0, len(endpoints))
	successCodes := make([]*SuccessCode, 0)
	notificationServices := make([]*NotificationService, 0)

	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		infos = append(infos, &Info{
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

	return &Endpoints{
		Infos:                infos,
		SuccessCodes:         successCodes,
		NotificationServices: notificationServices,
	}
}

// Convert structured Endpoints back to []*Endpoint
func ToServiceEndpoints(structured *Endpoints) srv_models.Endpoints {
	if structured == nil || len(structured.Infos) == 0 {
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

	endpoints := make(srv_models.Endpoints, 0, len(structured.Infos))
	for _, info := range structured.Infos {
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
