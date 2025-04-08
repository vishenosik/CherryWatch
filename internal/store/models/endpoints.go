package models

import (
	"iter"
	"time"
)

type Endpoint struct {
	ID                   string `db:"id"`
	ServiceName          string `db:"service_name"`
	URL                  string `db:"url"`
	SuccessCodes         []int
	NotificationServices []string
	Interval             time.Duration `db:"interval"`
}

type Endpoints = []Endpoint

type SuccessCode struct {
	ID   string `db:"endpoint_id"`
	Code int    `db:"code"`
}

type NotificationService struct {
	ID          string `db:"endpoint_id"`
	ServiceName string `db:"service_name"`
}

func SuccessCodesBatch(edps Endpoints) []SuccessCode {

	successCodes := make([]SuccessCode, 0, len(edps))

	for _, edp := range edps {
		for _, sc := range edp.SuccessCodes {
			successCodes = append(successCodes, SuccessCode{
				ID:   edp.ID,
				Code: sc,
			})
		}
	}

	return successCodes
}

func NotificationServicesBatch(edps Endpoints) []NotificationService {

	successCodes := make([]NotificationService, 0, len(edps))

	for _, edp := range edps {
		for _, sc := range edp.NotificationServices {
			successCodes = append(successCodes, NotificationService{
				ID:          edp.ID,
				ServiceName: sc,
			})
		}
	}

	return successCodes
}

func SuccessCodesIter(edp Endpoint) iter.Seq[int] {
	return Iter(edp.SuccessCodes)
}

func SCB(edps Endpoints) []SuccessCode {
	successCodes := make([]SuccessCode, 0, len(edps))

	// iter := func(yield func(Endpoint) bool) {
	// 	for _, edp := range edps {

	// 		for code := range SuccessCodesIter(edp) {
	// 			successCodes = append(successCodes, SuccessCode{
	// 				ID:   edp.ID,
	// 				Code: code,
	// 			})
	// 		}

	// 		if !yield(edp) {
	// 			return
	// 		}
	// 	}
	// }

	return successCodes

}

func Iter[S ~[]T, T any](slice S) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, elem := range slice {
			if !yield(elem) {
				return
			}
		}
	}
}
