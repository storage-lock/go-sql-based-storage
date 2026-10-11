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

func (x *FooStorage) Capabilities() []storage.StorageCapability {
	return x.SqlBasedStorage.Capabilities()
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

func run(manager storage.ConnectionManager[*sql.DB]) (string, error) {
	fooStorage, err := NewFooStorage(manager)
	if err != nil {
		return "", err
	}
	_ = fooStorage.Capabilities()
	sqlText, args := (&FooSqlProvider{}).NowTimestampSql(context.Background(), "locks")
	_, _ = sqlText, args
	return fooStorage.GetName(), nil
}

func start(open func(string, string) (*sql.DB, error), driverName, dsn string) (string, error) {
	db, err := open(driverName, dsn)
	if err != nil {
		return "", err
	}
	return run(storage.NewFixedSqlDBConnectionManager(db))
}

var openDatabase = sql.Open

func main() {
	name, err := start(openDatabase, "mysql", "xxx")
	if err != nil {
		panic(err)
	}
	fmt.Println(name)
}
