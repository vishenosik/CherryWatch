package sqlite

import (
	"context"
	"path"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	embed "github.com/vishenosik/CherryWatch"
	"github.com/vishenosik/gocherry/pkg/sql"
)

type cancelFunc = func()

func suite(t *testing.T) (*sqlx.DB, cancelFunc) {

	tempDir := t.TempDir() // Automatically cleaned up after test
	dbPath := filepath.Join(tempDir, "test.db")

	store, err := sql.NewSqliteStoreConfig(
		sql.SqliteConfig{
			StorePath: dbPath,
		},
		sql.WithMigration(
			embed.Migrations,
			path.Join(embed.MigrationsPath, "sqlite"),
		),
	)
	require.NoError(t, err)

	x, err := store.Open(context.TODO())
	require.NoError(t, err)

	return x, func() {
		store.Close(context.TODO())
	}
}
