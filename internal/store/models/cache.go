package models

import (
	"context"
	"time"
)

type CacheProvider interface {
	Set(
		ctx context.Context,
		key string,
		value any,
		expiration time.Duration,
	) error

	Get(
		ctx context.Context,
		key string,
	) (string, error)

	Delete(
		ctx context.Context,
		key string,
	) error
}
