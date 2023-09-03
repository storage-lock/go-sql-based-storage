package sql_based_storage

import (
	"context"
	"fmt"
	"github.com/storage-lock/go-storage"
)

// Sql92Provider SQL Provider的一个实现，SQL92标准
type Sql92Provider struct {
}

var _ SqlProvider = &Sql92Provider{}

func NewSql92Provider() *Sql92Provider {
	return &Sql92Provider{}
}

func (x *Sql92Provider) CreateTableSql(ctx context.Context, tableFullName string) (string, []any) {
	createTableSql := `CREATE TABLE IF NOT EXISTS %s (
    lock_id VARCHAR(255) NOT NULL PRIMARY KEY,
    owner_id VARCHAR(255) NOT NULL,
    version BIGINT NOT NULL,
    lock_information_json_string VARCHAR(255) NOT NULL
)`
	return fmt.Sprintf(createTableSql, tableFullName), nil
}

func (x *Sql92Provider) UpdateWithVersionSql(ctx context.Context, tableFullName string, lockId string, exceptedVersion, newVersion storage.Version, lockInformation *storage.LockInformation) (string, []any) {
	insertSql := fmt.Sprintf(`UPDATE %s SET version = ?, lock_information_json_string = ? WHERE lock_id = ? AND version = ?`, tableFullName)
	return insertSql, []any{newVersion, lockInformation.ToJsonString(), lockId, exceptedVersion}
}

func (x *Sql92Provider) CreateWithVersionSql(ctx context.Context, tableFullName string, lockId string, version storage.Version, lockInformation *storage.LockInformation) (string, []any) {
	insertSql := fmt.Sprintf(`INSERT INTO %s (lock_id, owner_id, version, lock_information_json_string) VALUES (?, ?, ?, ?)`, tableFullName)
	return insertSql, []any{lockId, lockInformation.OwnerId, version, lockInformation.ToJsonString()}
}

func (x *Sql92Provider) DeleteWithVersionSql(ctx context.Context, tableFullName string, lockId string, exceptedVersion storage.Version, lockInformation *storage.LockInformation) (string, []any) {
	deleteSql := fmt.Sprintf(`DELETE FROM %s WHERE lock_id = ? AND owner_id = ? AND version = ?`, tableFullName)
	return deleteSql, []any{lockId, lockInformation.OwnerId, exceptedVersion}
}

func (x *Sql92Provider) NowTimestampSql(ctx context.Context, tableFullName string) (string, []any) {
	return "SELECT UNIX_TIMESTAMP(NOW())", nil
}

func (x *Sql92Provider) FindLockInformationJsonStringByIdSql(ctx context.Context, tableFullName string, lockId string) (string, []any) {
	getLockSql := fmt.Sprintf("SELECT lock_information_json_string FROM %s WHERE lock_id = ?", tableFullName)
	return getLockSql, []any{lockId}
}

func (x *Sql92Provider) ListLockInformationJsonStringSql(ctx context.Context, tableFullName string) (string, []any) {
	sql := fmt.Sprintf("SELECT lock_information_json_string FROM %s", tableFullName)
	return sql, nil
}
