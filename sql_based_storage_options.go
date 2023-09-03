package sql_based_storage

import (
	"database/sql"
	"github.com/storage-lock/go-storage"
)

type SqlBasedStorageOptions struct {

	// 用于提供SQL
	SqlProvider SqlProvider

	// 用于管理连接
	ConnectionManager storage.ConnectionManager[*sql.DB]

	// 存储锁的表的名称
	TableFullName string
}

func NewSqlBasedStorageOptions() *SqlBasedStorageOptions {
	return &SqlBasedStorageOptions{
		SqlProvider:   NewSql92Provider(),
		TableFullName: storage.DefaultStorageTableName,
	}
}

func (x *SqlBasedStorageOptions) SetSqlProvider(sqlProvider SqlProvider) *SqlBasedStorageOptions {
	x.SqlProvider = sqlProvider
	return x
}

func (x *SqlBasedStorageOptions) SetConnectionManager(connectionManager storage.ConnectionManager[*sql.DB]) *SqlBasedStorageOptions {
	x.ConnectionManager = connectionManager
	return x
}

func (x *SqlBasedStorageOptions) SetTableFullName(tableFullName string) *SqlBasedStorageOptions {
	x.TableFullName = tableFullName
	return x
}

func (x *SqlBasedStorageOptions) Check() error {

	if x.SqlProvider == nil {
		return ErrSqlProviderCanNotNil
	}

	if x.ConnectionManager == nil {
		return ErrConnectionManagerCanNotNil
	}

	if x.TableFullName == "" {
		return ErrTableFullNameCanNotEmpty
	}

	return nil
}
