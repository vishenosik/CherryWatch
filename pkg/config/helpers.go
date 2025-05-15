package config

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Server struct {
	Host    string
	Port    uint16 `validate:"gte=1,lte=65535"`
	Timeout time.Duration
}

func (srv Server) Validate() error {
	if err := validator.New().Struct(srv); err != nil {
		return err
	}
	return nil
}

type Credentials struct {
	User     string
	Password string
}
