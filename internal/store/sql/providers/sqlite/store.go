package sqlite

import (
	"database/sql"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const (
	dialect string = "sqlite"
)

type Store struct {
	db *sql.DB
}

func MustInitSqlite(StorePath string) *Store {
	Store, err := NewSqliteStore(StorePath)
	if err != nil {
		panic(err)
	}
	return Store
}

func NewSqliteStore(StorePath string) (*Store, error) {

	const op = "Store.sqlite.New"

	db, err := sql.Open("sqlite3", StorePath)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &Store{
		db: db,
	}, nil
}

func (Store *Store) DB() *sql.DB {
	return Store.db
}

func (Store *Store) Dialect() string {
	return dialect
}

func (Store *Store) MigrationsPath() string {
	return path.Join("migrations", dialect)
}

func (Store *Store) Stop() error {
	return Store.db.Close()
}
