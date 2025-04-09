package sqlite

import (
	"database/sql"
	"path"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Store struct {
	db *sqlx.DB
	*endpoints
}

func MustInitSqlite(storePath string) *Store {
	str, err := NewSqliteStore(storePath)
	if err != nil {
		panic(err)
	}
	return str
}

func NewSqliteStore(storePath string) (*Store, error) {

	const op = "Store.sqlite.New"

	db, err := sqlx.Open("sqlite3", storePath)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	edps := newEndpoints(db)

	return &Store{
		db:        db,
		endpoints: edps,
	}, nil
}

func (str *Store) DB() *sql.DB {
	return str.db.DB
}

func (str *Store) Dialect() string {
	return "sqlite"
}

func (str *Store) MigrationsPath() string {
	return path.Join("migrations", str.Dialect())
}

func (str *Store) Stop() error {
	return str.db.Close()
}
