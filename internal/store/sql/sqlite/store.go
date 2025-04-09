package sqlite

import (
	"database/sql"
	"path"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	embed "github.com/vishenosik/CherryWatch"
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

	// Stores migration
	goose.SetBaseFS(embed.Migrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return nil, errors.Wrap(err, "could not set dialect "+"sqlite")
	}

	if err := goose.Up(db.DB, path.Join("migrations", "sqlite")); err != nil {
		return nil, errors.Wrap(err, "could not run migrations up")
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

func (str *Store) Close() error {
	return str.db.Close()
}
