package sqlstore

import (
	"database/sql"
)

type Store struct {
	provider StoreProvider
}

type StoreProvider interface {
	DB() *sql.DB
}

func NewStore(
	provider StoreProvider,
) *Store {
	return &Store{
		provider: provider,
	}
}

func (store *Store) Stop() error {
	return store.provider.DB().Close()
}
