package migrate

import (
	"database/sql"
	"io/fs"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

type Storage interface {
	DB() *sql.DB
	Dialect() string
	MigrationsPath() string
}

type migrator struct{}

func NewMigrator(
	logger goose.Logger,
	migrations fs.FS,
) *migrator {

	if logger != nil {
		goose.SetLogger(logger)
	}

	goose.SetBaseFS(migrations)

	return &migrator{}
}

func (m *migrator) MustMigrate(store Storage) {
	if err := m.Migrate(store); err != nil {
		panic(err)
	}
}

func (m *migrator) Migrate(store Storage) error {

	dialect := store.Dialect()

	if err := goose.SetDialect(dialect); err != nil {
		return errors.Wrap(err, "could not set dialect "+dialect)
	}

	if err := goose.Up(store.DB(), store.MigrationsPath()); err != nil {
		return errors.Wrap(err, "could not run migrations up")
	}
	return nil
}
