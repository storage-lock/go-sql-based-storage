package sql_based_storage

import (
	"context"
	"errors"
	"github.com/golang-infrastructure/go-iterator"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"time"
)

type SqlBasedStorage struct {
	options *SqlBasedStorageOptions
}

var _ storage.Storage = &SqlBasedStorage{}

func NewSqlBasedStorage(options *SqlBasedStorageOptions) (*SqlBasedStorage, error) {

	// 参数检查
	if err := options.Check(); err != nil {
		return nil, err
	}

	return &SqlBasedStorage{
		options: options,
	}, nil
}

const StorageName = "sql-based-storage"

func (x *SqlBasedStorage) GetName() string {
	return StorageName
}

func (x *SqlBasedStorage) Init(ctx context.Context) (returnError error) {

	db, err := x.options.ConnectionManager.Take(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err := x.options.ConnectionManager.Return(ctx, db)
		if returnError == nil {
			returnError = err
		}
	}()

	// 创建存储锁信息需要的表
	_, err = db.Exec(x.options.SqlProvider.CreateTableSql(ctx, x.options.TableFullName))
	if err != nil {
		return err
	}
	return nil
}

func (x *SqlBasedStorage) UpdateWithVersion(ctx context.Context, lockId string, exceptedVersion, newVersion storage.Version, lockInformation *storage.LockInformation) (returnError error) {

	db, err := x.options.ConnectionManager.Take(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err := x.options.ConnectionManager.Return(ctx, db)
		if returnError == nil {
			returnError = err
		}
	}()

	sql, params := x.options.SqlProvider.UpdateWithVersionSql(ctx, x.options.TableFullName, lockId, exceptedVersion, newVersion, lockInformation)
	execContext, err := db.ExecContext(ctx, sql, params...)
	if err != nil {
		return err
	}
	affected, err := execContext.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return storage_lock.ErrVersionMiss
	}
	return nil
}

func (x *SqlBasedStorage) CreateWithVersion(ctx context.Context, lockId string, version storage.Version, lockInformation *storage.LockInformation) (returnError error) {

	db, err := x.options.ConnectionManager.Take(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err := x.options.ConnectionManager.Return(ctx, db)
		if returnError == nil {
			returnError = err
		}
	}()

	sql, params := x.options.SqlProvider.CreateWithVersionSql(ctx, x.options.TableFullName, lockId, version, lockInformation)
	execContext, err := db.ExecContext(ctx, sql, params...)
	if err != nil {
		return err
	}
	affected, err := execContext.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return storage_lock.ErrVersionMiss
	}
	return nil
}

func (x *SqlBasedStorage) DeleteWithVersion(ctx context.Context, lockId string, exceptedVersion storage.Version, lockInformation *storage.LockInformation) (returnError error) {

	db, err := x.options.ConnectionManager.Take(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err := x.options.ConnectionManager.Return(ctx, db)
		if returnError == nil {
			returnError = err
		}
	}()

	sql, params := x.options.SqlProvider.DeleteWithVersionSql(ctx, x.options.TableFullName, lockId, exceptedVersion, lockInformation)
	execContext, err := db.ExecContext(ctx, sql, params...)
	if err != nil {
		return err
	}
	affected, err := execContext.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return storage_lock.ErrVersionMiss
	}
	return nil
}

func (x *SqlBasedStorage) Get(ctx context.Context, lockId string) (lockInformationJsonString string, returnError error) {

	db, err := x.options.ConnectionManager.Take(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		err := x.options.ConnectionManager.Return(ctx, db)
		if returnError == nil {
			returnError = err
		}
	}()

	sql, params := x.options.SqlProvider.SelectLockInformationJsonStringSql(ctx, x.options.TableFullName, lockId)
	rs, err := db.QueryContext(ctx, sql, params...)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = rs.Close()
	}()
	if !rs.Next() {
		return "", storage_lock.ErrLockNotFound
	}
	err = rs.Scan(&lockInformationJsonString)
	if err != nil {
		return "", err
	}
	return lockInformationJsonString, nil
}

func (x *SqlBasedStorage) GetTime(ctx context.Context) (now time.Time, returnError error) {

	db, err := x.options.ConnectionManager.Take(ctx)
	if err != nil {
		return time.Time{}, err
	}
	defer func() {
		err := x.options.ConnectionManager.Return(ctx, db)
		if returnError == nil {
			returnError = err
		}
	}()

	var zero time.Time
	// TODO 多实例的情况下可能会有问题，允许其能够比较方便的切换到NTP TimeProvider
	sql, params := x.options.SqlProvider.NowTimestampSql(ctx, x.options.TableFullName)
	rs, err := db.Query(sql, params...)
	if err != nil {
		return zero, err
	}
	defer func() {
		err := rs.Close()
		if returnError == nil {
			returnError = err
		}
	}()
	if !rs.Next() {
		return zero, errors.New("rs server time failed")
	}
	var databaseTimestamp uint64
	err = rs.Scan(&databaseTimestamp)
	if err != nil {
		return zero, err
	}

	// TODO 时区
	return time.Unix(int64(databaseTimestamp), 0), nil
}

func (x *SqlBasedStorage) Close(ctx context.Context) error {
	// 没有Storage级别的资源好回收的
	return nil
}

func (x *SqlBasedStorage) List(ctx context.Context) (iterator iterator.Iterator[*storage.LockInformation], returnError error) {

	db, err := x.options.ConnectionManager.Take(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := x.options.ConnectionManager.Return(ctx, db)
		if returnError == nil {
			returnError = err
		}
	}()

	sql, params := x.options.SqlProvider.ListLockInformationJsonStringSql(ctx, x.options.TableFullName)
	rows, err := db.QueryContext(ctx, sql, params...)
	if err != nil {
		return nil, err
	}
	return storage.NewSqlRowsIterator(rows), nil
}
