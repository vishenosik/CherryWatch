package embed

import "embed"

const (
	MigrationsPath = "migrations"
)

var (
	//go:embed migrations
	Migrations embed.FS
)
