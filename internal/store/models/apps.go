package models

import "fmt"

type App struct {
	Name   string
	Secret string
	ID     string `json:"-"`
}

func (app App) GetID() string {
	return app.ID
}

func (app App) GetSecret() []byte {
	return []byte(app.Secret)
}

func AppCacheKey(id string) string {
	return fmt.Sprintf("app:%s", id)
}
