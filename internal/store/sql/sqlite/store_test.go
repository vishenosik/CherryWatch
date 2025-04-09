package sqlite

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type cancelFunc = func()

func suite(t *testing.T) (*Store, cancelFunc) {

	tempDir := t.TempDir() // Automatically cleaned up after test
	dbPath := filepath.Join(tempDir, "test.db")

	store, err := NewSqliteStore(dbPath)
	require.NoError(t, err)

	return store, func() {
		store.Close()
	}

}
