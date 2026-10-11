package main

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/storage-lock/go-storage"
)

func TestFooStorageExample(t *testing.T) {
	name, err := run(storage.NewFixedSqlDBConnectionManager(&sql.DB{}))
	if err != nil || name != "foo-storage" {
		t.Fatalf("name=%s err=%v", name, err)
	}
	if _, err := NewFooStorage(nil); err == nil {
		t.Fatal("nil manager accepted")
	}
	if _, err := run(nil); err == nil {
		t.Fatal("nil manager run accepted")
	}
	openErr := errors.New("open failed")
	if _, err := start(func(string, string) (*sql.DB, error) { return nil, openErr }, "mysql", "dsn"); !errors.Is(err, openErr) {
		t.Fatal(err)
	}
	name, err = start(func(string, string) (*sql.DB, error) { return &sql.DB{}, nil }, "mysql", "dsn")
	if err != nil || name != "foo-storage" {
		t.Fatalf("name=%s err=%v", name, err)
	}
	openDatabase = func(string, string) (*sql.DB, error) { return &sql.DB{}, nil }
	main()
	openDatabase = func(string, string) (*sql.DB, error) { return nil, errors.New("boom") }
	defer func() { _ = recover() }()
	main()
}
