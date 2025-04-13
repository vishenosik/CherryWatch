package sqlite

import (
	"fmt"
	"path"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	embed "github.com/vishenosik/CherryWatch"
	"github.com/vishenosik/CherryWatch/internal/store/sql/models"
)

const (
	dialect = "sqlite"
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

	db, err := connect(storePath)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &Store{
		db:        db,
		endpoints: newEndpoints(db),
	}, nil
}

func (str *Store) Close() error {
	return str.db.Close()
}

func connect(storePath string) (*sqlx.DB, error) {

	db, err := sqlx.Open("sqlite3", storePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to sqlite")
	}

	goose.SetBaseFS(embed.Migrations)

	if err := goose.SetDialect(dialect); err != nil {
		return nil, fmt.Errorf("failed to set dialect %s: %w", dialect, err)
	}

	if err := goose.Up(db.DB, path.Join(embed.MigrationsPath, dialect)); err != nil {
		return nil, errors.Wrap(err, "failed to run migrations up")
	}

	return db, nil
}

func uniqueError(err error, mapper map[string]string) error {

	if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return nil
	}

	for fld := range mapper {
		if strings.Contains(err.Error(), fld) {
			return errors.Wrap(models.ErrAlreadyExists, mapper[fld])
		}
	}

	return nil
}
