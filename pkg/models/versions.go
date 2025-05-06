package models

type Version interface {
	Parse_(string) (Version, error)
	In_(v1, v2 Version) bool
}
