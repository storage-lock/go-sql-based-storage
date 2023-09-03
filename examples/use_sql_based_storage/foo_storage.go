package main

import (
	"context"
	"database/sql"
	"fmt"
	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
)

// ------------------------------------------------- --------------------------------------------------------------------

type FooStorage struct {
	*sql_based_storage.SqlBasedStorage
}

var _ storage.Storage = &FooStorage{}

func NewFooStorage(manager storage.ConnectionManager[*sql.DB]) (*FooStorage, error) {
	options := sql_based_storage.NewSqlBasedStorageOptions().SetConnectionManager(manager)
	storage, err := sql_based_storage.NewSqlBasedStorage(options)
	if err != nil {
		return nil, err
	}
	return &FooStorage{
		SqlBasedStorage: storage,
	}, nil
}

func (x *FooStorage) GetName() string {
	return "foo-storage"
}

// ------------------------------------------------- --------------------------------------------------------------------

type FooSqlProvider struct {
	sql_based_storage.Sql92Provider
}

var _ sql_based_storage.SqlProvider = &FooSqlProvider{}

func (x *FooSqlProvider) NowTimestampSql(ctx context.Context, tableFullName string) (string, []any) {
	return "SELECT NOW()", nil
}

// ------------------------------------------------- --------------------------------------------------------------------

func main() {
	db, err := sql.Open("mysql", "xxx")
	if err != nil {
		panic(err)
	}
	connectionManager := storage.NewFixedSqlDBConnectionManager(db)
	fooStorage, err := NewFooStorage(connectionManager)
	if err != nil {
		panic(err)
	}
	fmt.Println(fooStorage.GetName())
}
