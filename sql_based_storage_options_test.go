package sql_based_storage

import (
	"database/sql"
	"testing"

	"github.com/storage-lock/go-storage"
	"github.com/stretchr/testify/assert"
)

func TestNewSqlBasedStorageOptions_Defaults(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	assert.NotNil(t, opts)
	assert.NotNil(t, opts.SqlProvider)
	assert.Equal(t, storage.DefaultStorageTableName, opts.TableFullName)
	assert.Nil(t, opts.ConnectionManager)
}

func TestSqlBasedStorageOptions_SetSqlProvider_Chainable(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	result := opts.SetSqlProvider(NewSql92Provider())
	assert.Equal(t, opts, result)
	assert.NotNil(t, opts.SqlProvider)
}

func TestSqlBasedStorageOptions_SetConnectionManager_Chainable(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	cm := storage.NewFixedSqlDBConnectionManager(&sql.DB{})
	result := opts.SetConnectionManager(cm)
	assert.Equal(t, opts, result)
	assert.Equal(t, cm, opts.ConnectionManager)
}

func TestSqlBasedStorageOptions_SetTableFullName_Chainable(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	result := opts.SetTableFullName("custom_table")
	assert.Equal(t, opts, result)
	assert.Equal(t, "custom_table", opts.TableFullName)
}

func TestSqlBasedStorageOptions_Check_NilSqlProvider(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	opts.SqlProvider = nil
	err := opts.Check()
	assert.Error(t, err)
	assert.Equal(t, ErrSqlProviderCanNotNil, err)
}

func TestSqlBasedStorageOptions_Check_NilConnectionManager(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	opts.ConnectionManager = nil
	err := opts.Check()
	assert.Error(t, err)
	assert.Equal(t, ErrConnectionManagerCanNotNil, err)
}

func TestSqlBasedStorageOptions_Check_EmptyTableFullName(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	opts.ConnectionManager = storage.NewFixedSqlDBConnectionManager(&sql.DB{})
	opts.TableFullName = ""
	err := opts.Check()
	assert.Error(t, err)
	assert.Equal(t, ErrTableFullNameCanNotEmpty, err)
}

func TestSqlBasedStorageOptions_Check_ValidOptions(t *testing.T) {
	opts := NewSqlBasedStorageOptions()
	opts.ConnectionManager = storage.NewFixedSqlDBConnectionManager(&sql.DB{})
	err := opts.Check()
	assert.NoError(t, err)
}
