package models

import "fmt"

type User struct {
	Nickname     string
	Email        string
	ID           string `json:"-"`
	PasswordHash []byte `json:"-"`
	IsAdmin      bool
}

func (user User) GetID() string {
	return user.ID
}

func (user User) GetEmail() string {
	return user.Email
}

func UserCacheKey(id string) string {
	return fmt.Sprintf("user:%s", id)
}
