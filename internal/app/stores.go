package app

import (
	"context"

	embed "github.com/vishenosik/CherryWatch"
	appctx "github.com/vishenosik/CherryWatch/internal/app/context"
	sqlstore "github.com/vishenosik/CherryWatch/internal/store/sql"
	"github.com/vishenosik/CherryWatch/internal/store/sql/providers/sqlite"
	std "github.com/vishenosik/web-tools/log"
	"github.com/vishenosik/web-tools/migrate"
)

func loadSqlStore(ctx context.Context) (*sqlstore.Store, error) {

	appContext := appctx.AppCtx(ctx)

	// Stores init
	sqliteStore := sqlite.MustInitSqlite(appContext.Config.StorePath)
	store := sqlstore.NewStore(sqliteStore)

	// Stores migration
	migrate.NewMigrator(
		std.NewStdLogger(appContext.Logger),
		embed.Migrations,
	).MustMigrate(sqliteStore)

	return store, nil
}
