package sql_based_storage

import (
	"context"
	"github.com/storage-lock/go-storage"
)

// SqlProvider 用于提供各种用途的SQl方言
type SqlProvider interface {

	// CreateTableSql 返回创建存储锁的表的SQL，创建时请根据需要自行设置主键约束
	CreateTableSql(ctx context.Context, tableFullName string) (string, []any)

	// UpdateWithVersionSql 根据LockId、Version更新LockInformation的SQL
	UpdateWithVersionSql(ctx context.Context, tableFullName string, lockId string, exceptedVersion, newVersion storage.Version, lockInformation *storage.LockInformation) (string, []any)

	// CreateWithVersionSql 根据LockId、Version创建LockInformation的SQL
	CreateWithVersionSql(ctx context.Context, tableFullName string, lockId string, version storage.Version, lockInformation *storage.LockInformation) (string, []any)

	// DeleteWithVersionSql 根据LockId和Version删除锁信息的SQL
	DeleteWithVersionSql(ctx context.Context, tableFullName string, lockId string, exceptedVersion storage.Version, lockInformation *storage.LockInformation) (string, []any)

	// NowTimestampSql 获取Storage当前时间的Unix时间戳的SQL
	NowTimestampSql(ctx context.Context, tableFullName string) (string, []any)

	// FindLockInformationJsonStringByIdSql 根据LockId查询锁信息的JsonString的SQL
	FindLockInformationJsonStringByIdSql(ctx context.Context, tableFullName string, lockId string) (string, []any)

	// ListLockInformationJsonStringSql 返回能够列出所有的锁的JsonString的SQL
	ListLockInformationJsonStringSql(ctx context.Context, tableFullName string) (string, []any)
}
