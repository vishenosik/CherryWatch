package embed

import "embed"

var (
	//go:embed migrations
	Migrations embed.FS
)
