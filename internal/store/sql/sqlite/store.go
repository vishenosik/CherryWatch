package sqlite

import (
	"database/sql"
	"path"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type store struct {
	db *sqlx.DB
	*endpoints
}

func MustInitSqlite(storePath string) *store {
	Store, err := NewSqliteStore(storePath)
	if err != nil {
		panic(err)
	}
	return Store
}

func NewSqliteStore(storePath string) (*store, error) {

	const op = "Store.sqlite.New"

	db, err := sqlx.Open("sqlite3", storePath)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	edps := newEndpoints(db)

	return &store{
		db:        db,
		endpoints: edps,
	}, nil
}

func (str *store) DB() *sql.DB {
	return str.db.DB
}

func (str *store) Dialect() string {
	return "sqlite"
}

func (str *store) MigrationsPath() string {
	return path.Join("migrations", str.Dialect())
}

func (str *store) Stop() error {
	return str.db.Close()
}
