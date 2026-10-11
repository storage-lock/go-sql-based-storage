package sql_based_storage

import (
	"context"
	"testing"

	"github.com/storage-lock/go-storage"
	"github.com/stretchr/testify/assert"
)

func TestNewSql92Provider_ReturnsNonNil(t *testing.T) {
	provider := NewSql92Provider()
	assert.NotNil(t, provider)
}

func TestSql92Provider_CreateTableSql(t *testing.T) {
	provider := NewSql92Provider()
	ctx := context.Background()
	sql, params := provider.CreateTableSql(ctx, "test_table")

	assert.Contains(t, sql, "CREATE TABLE IF NOT EXISTS")
	assert.Contains(t, sql, "test_table")
	assert.Contains(t, sql, "lock_id")
	assert.Contains(t, sql, "owner_id")
	assert.Contains(t, sql, "version")
	assert.Contains(t, sql, "lock_information_json_string")
	assert.Empty(t, params)
}

func TestSql92Provider_CreateTableSql_DefaultTableName(t *testing.T) {
	provider := NewSql92Provider()
	ctx := context.Background()
	sql, _ := provider.CreateTableSql(ctx, storage.DefaultStorageTableName)

	assert.Contains(t, sql, storage.DefaultStorageTableName)
}

func TestSql92Provider_UpdateWithVersionSql(t *testing.T) {
	provider := NewSql92Provider()
	ctx := context.Background()

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   1,
		LockCount: 1,
	}

	sql, params := provider.UpdateWithVersionSql(ctx, "test_table", "lock-1", 1, 2, lockInfo)

	assert.Contains(t, sql, "UPDATE")
	assert.Contains(t, sql, "test_table")
	assert.Contains(t, sql, "WHERE lock_id = ? AND version = ?")
	assert.Len(t, params, 4)
	assert.Equal(t, storage.Version(2), params[0])
	assert.Equal(t, lockInfo.ToJsonString(), params[1])
	assert.Equal(t, "lock-1", params[2])
	assert.Equal(t, storage.Version(1), params[3])
}

func TestSql92Provider_CreateWithVersionSql(t *testing.T) {
	provider := NewSql92Provider()
	ctx := context.Background()

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   1,
		LockCount: 1,
	}

	sql, params := provider.CreateWithVersionSql(ctx, "test_table", "lock-1", 1, lockInfo)

	assert.Contains(t, sql, "INSERT INTO")
	assert.Contains(t, sql, "test_table")
	assert.Contains(t, sql, "lock_id, owner_id, version, lock_information_json_string")
	assert.Len(t, params, 4)
	assert.Equal(t, "lock-1", params[0])
	assert.Equal(t, "owner-1", params[1])
	assert.Equal(t, storage.Version(1), params[2])
	assert.Equal(t, lockInfo.ToJsonString(), params[3])
}

func TestSql92Provider_DeleteWithVersionSql(t *testing.T) {
	provider := NewSql92Provider()
	ctx := context.Background()

	lockInfo := &storage.LockInformation{
		LockId:    "lock-1",
		OwnerId:   "owner-1",
		Version:   1,
		LockCount: 1,
	}

	sql, params := provider.DeleteWithVersionSql(ctx, "test_table", "lock-1", 1, lockInfo)

	assert.Contains(t, sql, "DELETE FROM")
	assert.Contains(t, sql, "test_table")
	assert.Contains(t, sql, "WHERE lock_id = ? AND owner_id = ? AND version = ?")
	assert.Len(t, params, 3)
	assert.Equal(t, "lock-1", params[0])
	assert.Equal(t, "owner-1", params[1])
	assert.Equal(t, storage.Version(1), params[2])
}

func TestSql92Provider_NowTimestampSql(t *testing.T) {
	provider := NewSql92Provider()
	ctx := context.Background()

	sql, params := provider.NowTimestampSql(ctx, "test_table")

	assert.Equal(t, "SELECT UNIX_TIMESTAMP(NOW())", sql)
	assert.Empty(t, params)
}

func TestSql92Provider_FindLockInformationJsonStringByIdSql(t *testing.T) {
	provider := NewSql92Provider()
	ctx := context.Background()

	sql, params := provider.FindLockInformationJsonStringByIdSql(ctx, "test_table", "lock-1")

	assert.Contains(t, sql, "SELECT lock_information_json_string")
	assert.Contains(t, sql, "FROM test_table")
	assert.Contains(t, sql, "WHERE lock_id = ?")
	assert.Len(t, params, 1)
	assert.Equal(t, "lock-1", params[0])
}

func TestSql92Provider_ListLockInformationJsonStringSql(t *testing.T) {
	provider := NewSql92Provider()
	ctx := context.Background()

	sql, params := provider.ListLockInformationJsonStringSql(ctx, "test_table")

	assert.Contains(t, sql, "SELECT lock_information_json_string")
	assert.Contains(t, sql, "FROM test_table")
	assert.Empty(t, params)
}

func TestSql92Provider_Interface_Satisfied(t *testing.T) {
	// 验证 Sql92Provider 实现了 SqlProvider 接口
	var _ SqlProvider = NewSql92Provider()
}
