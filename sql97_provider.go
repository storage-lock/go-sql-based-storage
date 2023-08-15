package sql_based_storage

import (
	"context"
	"fmt"
	"github.com/storage-lock/go-storage"
)

type Sql97Provider struct {
}

var _ SqlProvider = &Sql97Provider{}

func NewSql97Provider() *Sql97Provider {
	return &Sql97Provider{}
}

func (x *Sql97Provider) CreateTableSql(ctx context.Context, tableFullName string) (string, []any) {
	createTableSql := `CREATE TABLE IF NOT EXISTS %s (
    lock_id VARCHAR(255) NOT NULL PRIMARY KEY,
    owner_id VARCHAR(255) NOT NULL,
    version BIGINT NOT NULL,
    lock_information_json_string VARCHAR(255) NOT NULL
)`
	return fmt.Sprintf(createTableSql, tableFullName), nil
}

func (x *Sql97Provider) UpdateWithVersionSql(ctx context.Context, tableFullName string, lockId string, exceptedVersion, newVersion storage.Version, lockInformation *storage.LockInformation) (string, []any) {
	insertSql := fmt.Sprintf(`UPDATE %s SET version = ?, lock_information_json_string = ? WHERE lock_id = ? AND owner_id = ? AND version = ?`, tableFullName)
	return insertSql, []any{newVersion, lockInformation.ToJsonString(), lockId, lockInformation.OwnerId, exceptedVersion}
}

func (x *Sql97Provider) CreateWithVersionSql(ctx context.Context, tableFullName string, lockId string, version storage.Version, lockInformation *storage.LockInformation) (string, []any) {
	insertSql := fmt.Sprintf(`INSERT INTO %s (lock_id, owner_id, version, lock_information_json_string) VALUES (?, ?, ?, ?)`, tableFullName)
	return insertSql, []any{lockId, lockInformation.OwnerId, version, lockInformation.ToJsonString()}
}

func (x *Sql97Provider) DeleteWithVersionSql(ctx context.Context, tableFullName string, lockId string, exceptedVersion storage.Version, lockInformation *storage.LockInformation) (string, []any) {
	deleteSql := fmt.Sprintf(`DELETE FROM %s WHERE lock_id = ? AND owner_id = ? AND version = ?`, tableFullName)
	return deleteSql, []any{lockId, lockInformation.OwnerId, exceptedVersion}
}

func (x *Sql97Provider) NowTimestampSql(ctx context.Context, tableFullName string) (string, []any) {
	// 下面这个通用嘛？
	// "SELECT UNIX_TIMESTAMP(NOW())"
	//TODO implement me
	panic("implement me")
}

func (x *Sql97Provider) SelectLockInformationJsonStringSql(ctx context.Context, tableFullName string, lockId string) (string, []any) {
	getLockSql := fmt.Sprintf("SELECT lock_information_json_string FROM %s WHERE lock_id = ?", tableFullName)
	return getLockSql, []any{lockId}
}

func (x *Sql97Provider) ListLockInformationJsonStringSql(ctx context.Context, tableFullName string) (string, []any) {
	sql := fmt.Sprintf("SELECT lock_information_json_string FROM %s", tableFullName)
	return sql, nil
}
