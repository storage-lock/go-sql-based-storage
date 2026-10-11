package sql_based_storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/storage-lock/go-storage"
	"github.com/stretchr/testify/assert"
)

func TestNewSqlBasedStorageOptions(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	assert.NotNil(t, opts)
	assert.NotEmpty(t, opts.TableFullName)
	assert.Nil(t, opts.ConnectionManager)
	assert.NotNil(t, opts.SqlProvider)
}

func TestNewSql92Provider(t *testing.T) {
	provider := NewSql92Provider()
	assert.NotNil(t, provider)
}

func TestNewSqlBasedStorage_NilConnectionManager(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	s, err := NewSqlBasedStorage(opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ConnectionManager")
	assert.Nil(t, s)
}

func TestSqlBasedStorage_GetName(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	opts.SetConnectionManager(storage.NewFixedSqlDBConnectionManager(&sql.DB{}))
	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, StorageName, s.GetName())
}

func TestSqlBasedStorage_Capabilities(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	opts.SetConnectionManager(storage.NewFixedSqlDBConnectionManager(&sql.DB{}))
	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	caps := s.Capabilities()
	assert.Contains(t, caps, storage.CapabilityCAS)
	assert.Contains(t, caps, storage.CapabilityReliableTime)
}

func TestSqlBasedStorage_Close(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	opts.SetConnectionManager(storage.NewFixedSqlDBConnectionManager(&sql.DB{}))
	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Close should not error even without initialization
	ctx := context.Background()
	err = s.Close(ctx)
	assert.NoError(t, err)
}
