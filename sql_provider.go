package sql_based_storage

import (
	"context"
	"github.com/storage-lock/go-storage"
)

// SqlProvider 用于提供各种SQl方言
type SqlProvider interface {

	// CreateTableSql 构造创建存储锁的表的sql语句以及参数
	CreateTableSql(ctx context.Context, tableFullName string) (string, []any)

	// UpdateWithVersionSql 构造根据版本更新锁信息的sql
	UpdateWithVersionSql(ctx context.Context, tableFullName string, lockId string, exceptedVersion, newVersion storage.Version, lockInformation *storage.LockInformation) (string, []any)

	// CreateWithVersionSql 构造根据版本创建锁信息的sql
	CreateWithVersionSql(ctx context.Context, tableFullName string, lockId string, version storage.Version, lockInformation *storage.LockInformation) (string, []any)

	// DeleteWithVersionSql 构造根据版本删除锁信息的sql
	DeleteWithVersionSql(ctx context.Context, tableFullName string, lockId string, exceptedVersion storage.Version, lockInformation *storage.LockInformation) (string, []any)

	// NowTimestampSql 构造获取锁的unix mil时间戳的sql
	NowTimestampSql(ctx context.Context, tableFullName string) (string, []any)

	// SelectLockInformationJsonStringSql 构造根据 lock_id 查询锁信息的sql
	SelectLockInformationJsonStringSql(ctx context.Context, tableFullName string, lockId string) (string, []any)

	// ListLockInformationJsonStringSql 提供列出所有的锁的json string的sql
	ListLockInformationJsonStringSql(ctx context.Context, tableFullName string) (string, []any)
}
