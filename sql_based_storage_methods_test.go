package sql_based_storage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/storage-lock/go-storage"
	"github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
)

func newTestStorage(t *testing.T) (*SqlBasedStorage, *sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建 sqlmock 失败: %v", err)
	}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(storage.NewFixedSqlDBConnectionManager(db)).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	if err != nil {
		t.Fatalf("创建 SqlBasedStorage 失败: %v", err)
	}

	return s, db, mock
}

func TestSqlBasedStorage_Init_CreatesTable(t *testing.T) {
	s, _, mock := newTestStorage(t)

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err := s.Init(ctx)
	assert.NoError(t, err)
}

func TestSqlBasedStorage_Init_TakeError(t *testing.T) {
	cm := &mockTakeErrorConnectionManager{}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	err = s.Init(ctx)
	assert.Error(t, err)
}

func TestSqlBasedStorage_Init_ExecContextError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnError(assert.AnError)

	ctx := context.Background()
	err := s.Init(ctx)
	assert.Error(t, err)
}

func TestSqlBasedStorage_Init_ReturnError(t *testing.T) {
	// 使用 sqlmock 创建可用的 DB 实例，避免 &sql.DB{} 在 ExecContext 时 panic
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(1, 1))

	cm := &mockErrorReturnConnectionManagerWithDB{db: sqlDB}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	err = s.Init(ctx)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestSqlBasedStorage_UpdateWithVersion_ReturnError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// 预期 UPDATE 成功，然后在 Return 时返回错误
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))

	cm := &mockErrorReturnConnectionManagerWithDB{db: sqlDB}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   2,
		LockCount: 2,
	}

	err = s.UpdateWithVersion(ctx, "lock-1", 1, 2, lockInfo)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestSqlBasedStorage_DeleteWithVersion_ReturnError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// 预期 DELETE 成功，然后在 Return 时返回错误
	mock.ExpectExec("DELETE FROM").WillReturnResult(sqlmock.NewResult(1, 1))

	cm := &mockErrorReturnConnectionManagerWithDB{db: sqlDB}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   1,
		LockCount: 1,
	}

	err = s.DeleteWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestSqlBasedStorage_CreateWithVersion_ReturnError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// 预期 INSERT 成功，然后在 Return 时返回错误
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))

	cm := &mockErrorReturnConnectionManagerWithDB{db: sqlDB}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	lockInfo := &storage.LockInformation{
		LockId:          "lock-1",
		OwnerId:         "owner-1",
		Version:         1,
		LockCount:       1,
		LockBeginTime:   time.Now(),
		LeaseExpireTime: time.Now().Add(time.Hour),
	}

	err = s.CreateWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestSqlBasedStorage_CreateWithVersion_Success(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{
		LockId:          "lock-1",
		OwnerId:         "owner-1",
		Version:         1,
		LockCount:       1,
		LockBeginTime:   time.Now(),
		LeaseExpireTime: time.Now().Add(time.Hour),
	}

	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err := s.CreateWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.NoError(t, err)
}

func TestSqlBasedStorage_CreateWithVersion_VersionMiss(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   1,
		LockCount: 1,
	}

	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(0, 0))

	ctx := context.Background()
	err := s.CreateWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.ErrorIs(t, err, storage_lock.ErrVersionMiss)
}

type errorResult struct{}

func (e *errorResult) LastInsertId() (int64, error) {
	return 0, assert.AnError
}

func (e *errorResult) RowsAffected() (int64, error) {
	return 0, assert.AnError
}

func TestSqlBasedStorage_CreateWithVersion_RowsAffectedError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   1,
		LockCount: 1,
	}

	mock.ExpectExec("INSERT INTO").WillReturnResult(&errorResult{})

	ctx := context.Background()
	err := s.CreateWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.Error(t, err)
}

func TestSqlBasedStorage_CreateWithVersion_TakeError(t *testing.T) {
	cm := &mockTakeErrorConnectionManager{}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	lockInfo := &storage.LockInformation{LockId: "lock-1"}
	err = s.CreateWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.Error(t, err)
}

func TestSqlBasedStorage_CreateWithVersion_ExecContextError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{LockId: "lock-1"}
	mock.ExpectExec("INSERT INTO").WillReturnError(assert.AnError)

	ctx := context.Background()
	err := s.CreateWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.Error(t, err)
}

func TestSqlBasedStorage_UpdateWithVersion_Success(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   2,
		LockCount: 2,
	}

	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err := s.UpdateWithVersion(ctx, "lock-1", 1, 2, lockInfo)
	assert.NoError(t, err)
}

func TestSqlBasedStorage_UpdateWithVersion_VersionMiss(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   2,
		LockCount: 2,
	}

	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(0, 0))

	ctx := context.Background()
	err := s.UpdateWithVersion(ctx, "lock-1", 1, 2, lockInfo)
	assert.ErrorIs(t, err, storage_lock.ErrVersionMiss)
}

func TestSqlBasedStorage_UpdateWithVersion_RowsAffectedError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   2,
		LockCount: 2,
	}

	mock.ExpectExec("UPDATE").WillReturnResult(&errorResult{})

	ctx := context.Background()
	err := s.UpdateWithVersion(ctx, "lock-1", 1, 2, lockInfo)
	assert.Error(t, err)
}

func TestSqlBasedStorage_UpdateWithVersion_TakeError(t *testing.T) {
	cm := &mockTakeErrorConnectionManager{}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	lockInfo := &storage.LockInformation{LockId: "lock-1"}
	err = s.UpdateWithVersion(ctx, "lock-1", 1, 2, lockInfo)
	assert.Error(t, err)
}

func TestSqlBasedStorage_UpdateWithVersion_ExecContextError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{LockId: "lock-1"}
	mock.ExpectExec("UPDATE").WillReturnError(assert.AnError)

	ctx := context.Background()
	err := s.UpdateWithVersion(ctx, "lock-1", 1, 2, lockInfo)
	assert.Error(t, err)
}

func TestSqlBasedStorage_DeleteWithVersion_Success(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   1,
		LockCount: 1,
	}

	mock.ExpectExec("DELETE FROM").WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err := s.DeleteWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.NoError(t, err)
}

func TestSqlBasedStorage_DeleteWithVersion_VersionMiss(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   1,
		LockCount: 1,
	}

	mock.ExpectExec("DELETE FROM").WillReturnResult(sqlmock.NewResult(0, 0))

	ctx := context.Background()
	err := s.DeleteWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.ErrorIs(t, err, storage_lock.ErrVersionMiss)
}

func TestSqlBasedStorage_DeleteWithVersion_RowsAffectedError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   1,
		LockCount: 1,
	}

	mock.ExpectExec("DELETE FROM").WillReturnResult(&errorResult{})

	ctx := context.Background()
	err := s.DeleteWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.Error(t, err)
}

func TestSqlBasedStorage_DeleteWithVersion_TakeError(t *testing.T) {
	cm := &mockTakeErrorConnectionManager{}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	lockInfo := &storage.LockInformation{LockId: "lock-1"}
	err = s.DeleteWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.Error(t, err)
}

func TestSqlBasedStorage_DeleteWithVersion_ExecContextError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{LockId: "lock-1"}
	mock.ExpectExec("DELETE FROM").WillReturnError(assert.AnError)

	ctx := context.Background()
	err := s.DeleteWithVersion(ctx, "lock-1", 1, lockInfo)
	assert.Error(t, err)
}

func TestSqlBasedStorage_Get_Success(t *testing.T) {
	s, _, mock := newTestStorage(t)

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   1,
		LockCount: 1,
	}

	jsonStr := lockInfo.ToJsonString()
	mock.ExpectQuery("SELECT lock_information_json_string").WillReturnRows(
		sqlmock.NewRows([]string{"lock_information_json_string"}).AddRow(jsonStr),
	)

	ctx := context.Background()
	result, err := s.Get(ctx, "lock-1")
	assert.NoError(t, err)
	assert.Equal(t, jsonStr, result)
}

func TestSqlBasedStorage_Get_QueryError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	mock.ExpectQuery("SELECT lock_information_json_string").WillReturnError(assert.AnError)

	ctx := context.Background()
	_, err := s.Get(ctx, "lock-1")
	assert.Error(t, err)
}

func TestSqlBasedStorage_Get_TakeError(t *testing.T) {
	cm := &mockTakeErrorConnectionManager{}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = s.Get(ctx, "lock-1")
	assert.Error(t, err)
}

func TestSqlBasedStorage_Get_NotFound(t *testing.T) {
	s, _, mock := newTestStorage(t)

	mock.ExpectQuery("SELECT lock_information_json_string").WillReturnRows(
		sqlmock.NewRows([]string{"lock_information_json_string"}),
	)

	ctx := context.Background()
	_, err := s.Get(ctx, "lock-1")
	assert.ErrorIs(t, err, storage_lock.ErrLockNotFound)
}

func TestSqlBasedStorage_Get_ScanError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	mock.ExpectQuery("SELECT lock_information_json_string").WillReturnRows(
		sqlmock.NewRows([]string{"lock_information_json_string", "extra_column"}).AddRow("invalid-json", "extra"),
	)

	ctx := context.Background()
	_, err := s.Get(ctx, "lock-1")
	assert.Error(t, err)
}

func TestSqlBasedStorage_Get_ReturnError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT lock_information_json_string").WillReturnRows(
		sqlmock.NewRows([]string{"lock_information_json_string"}).AddRow("{}"),
	)

	cm := &mockErrorReturnConnectionManagerWithDB{db: sqlDB}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = s.Get(ctx, "lock-1")
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestSqlBasedStorage_GetTime_Success(t *testing.T) {
	s, _, mock := newTestStorage(t)

	mock.ExpectQuery("SELECT UNIX_TIMESTAMP").WillReturnRows(
		sqlmock.NewRows([]string{"timestamp"}).AddRow(int64(1234567890)),
	)

	ctx := context.Background()
	result, err := s.GetTime(ctx)
	assert.NoError(t, err)
	assert.Equal(t, time.Unix(1234567890, 0), result)
}

func TestSqlBasedStorage_GetTime_QueryError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	mock.ExpectQuery("SELECT UNIX_TIMESTAMP").WillReturnError(assert.AnError)

	ctx := context.Background()
	_, err := s.GetTime(ctx)
	assert.Error(t, err)
}

func TestSqlBasedStorage_GetTime_NoRows(t *testing.T) {
	s, _, mock := newTestStorage(t)

	mock.ExpectQuery("SELECT UNIX_TIMESTAMP").WillReturnRows(
		sqlmock.NewRows([]string{"timestamp"}),
	)

	ctx := context.Background()
	_, err := s.GetTime(ctx)
	assert.Error(t, err)
}

func TestSqlBasedStorage_GetTime_ScanError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	mock.ExpectQuery("SELECT UNIX_TIMESTAMP").WillReturnRows(
		sqlmock.NewRows([]string{"timestamp"}).AddRow("not-a-timestamp"),
	)

	ctx := context.Background()
	_, err := s.GetTime(ctx)
	assert.Error(t, err)
}

func TestSqlBasedStorage_GetTime_ReturnError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT UNIX_TIMESTAMP").WillReturnRows(
		sqlmock.NewRows([]string{"timestamp"}).AddRow(int64(1234567890)),
	)

	cm := &mockErrorReturnConnectionManagerWithDB{db: sqlDB}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = s.GetTime(ctx)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestSqlBasedStorage_GetTime_TakeError(t *testing.T) {
	cm := &mockTakeErrorConnectionManager{}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = s.GetTime(ctx)
	assert.Error(t, err)
}

func TestSqlBasedStorage_List_Success(t *testing.T) {
	s, db, mock := newTestStorage(t)

	lockInfo1 := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   1,
		LockCount: 1,
	}
	lockInfo2 := &storage.LockInformation{
		LockId:    "lock-2",
		OwnerId:   "owner-2",
		Version:   1,
		LockCount: 1,
	}

	mock.ExpectQuery("SELECT lock_information_json_string").WillReturnRows(
		sqlmock.NewRows([]string{"lock_information_json_string"}).
			AddRow(lockInfo1.ToJsonString()).
			AddRow(lockInfo2.ToJsonString()),
	)

	ctx := context.Background()
	iter, err := s.List(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, iter)

	assert.True(t, iter.Next())
	item1 := iter.Value()
	assert.Equal(t, "lock-1", item1.LockId)

	assert.True(t, iter.Next())
	item2 := iter.Value()
	assert.Equal(t, "lock-2", item2.LockId)

	assert.False(t, iter.Next())

	_ = db.Close()
}

func TestSqlBasedStorage_List_QueryError(t *testing.T) {
	s, _, mock := newTestStorage(t)

	mock.ExpectQuery("SELECT lock_information_json_string").WillReturnError(assert.AnError)

	ctx := context.Background()
	_, err := s.List(ctx)
	assert.Error(t, err)
}

func TestSqlBasedStorage_List_TakeError(t *testing.T) {
	cm := &mockTakeErrorConnectionManager{}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = s.List(ctx)
	assert.Error(t, err)
}

func TestSqlBasedStorage_List_ReturnError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT lock_information_json_string").WillReturnRows(
		sqlmock.NewRows([]string{"lock_information_json_string"}),
	)

	cm := &mockErrorReturnConnectionManagerWithDB{db: sqlDB}

	opts := NewSqlBasedStorageOptions().
		SetSqlProvider(NewSql92Provider()).
		SetConnectionManager(cm).
		SetTableFullName("test_locks")

	s, err := NewSqlBasedStorage(opts)
	assert.NoError(t, err)

	ctx := context.Background()
	_, err = s.List(ctx)
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestSqlBasedStorage_List_Empty(t *testing.T) {
	s, db, mock := newTestStorage(t)

	mock.ExpectQuery("SELECT lock_information_json_string").WillReturnRows(
		sqlmock.NewRows([]string{"lock_information_json_string"}),
	)

	ctx := context.Background()
	iter, err := s.List(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, iter)

	assert.False(t, iter.Next())

	_ = db.Close()
}

func TestSqlBasedStorage_Interface_Satisfied(t *testing.T) {
	// 验证 SqlBasedStorage 实现了 Storage 接口
	var _ storage.Storage = &SqlBasedStorage{}
}

type mockTakeErrorConnectionManager struct{}

func (m *mockTakeErrorConnectionManager) Name() string {
	return "mock-take-error"
}

func (m *mockTakeErrorConnectionManager) Take(ctx context.Context) (*sql.DB, error) {
	return nil, assert.AnError
}

func (m *mockTakeErrorConnectionManager) Return(ctx context.Context, db *sql.DB) error {
	return nil
}

func (m *mockTakeErrorConnectionManager) Shutdown(ctx context.Context) error {
	return nil
}

type mockErrorReturnConnectionManagerWithDB struct {
	db *sql.DB
}

func (m *mockErrorReturnConnectionManagerWithDB) Name() string {
	return "mock-error-return-with-db"
}

func (m *mockErrorReturnConnectionManagerWithDB) Take(ctx context.Context) (*sql.DB, error) {
	return m.db, nil
}

func (m *mockErrorReturnConnectionManagerWithDB) Return(ctx context.Context, db *sql.DB) error {
	return assert.AnError
}

func (m *mockErrorReturnConnectionManagerWithDB) Shutdown(ctx context.Context) error {
	return nil
}

func TestSqlBasedStorageWatch(t *testing.T) {
	s, _, _ := newTestStorage(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := s.Watch(ctx, "watched")
	assert.NoError(t, err)
	assert.NotNil(t, ch)
	_, err = s.Watch(ctx, "")
	assert.Error(t, err)
}
