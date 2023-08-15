# SQL Based Storage 

# 一、这是什么？

基于SQL的关系型数据库的`Storage`的通用实现。

# 二、安装依赖

```bash
go get -u github.com/storage-lock/go-sql-based-storage
```

# 三、组件介绍

## SqlBasedStorage

把基于`SQL`构建`Storage`的流程的公共部分抽象出来，这样对于所有支持`SQL`的存储设备来说就只需要实现不同的部分就可以了，不同的部分使用`SqlProvider`来封装不同的`SQL`方言。

## SqlBasedStorageOptions

创建`SqlBasedStorage`的时候需要提供一些上下文参数，是通过这个`Options`来传递的： 

```go
type SqlBasedStorageOptions struct {
	SqlProvider       SqlProvider
	ConnectionManager storage.ConnectionManager[*sql.DB]
	TableFullName     string
}
```

## SqlProvider

用于封装不同的`SQL`方言：

```go
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
```

## Sql97Provider

对`SqlProvider`的`sql97`实现。

# 四、使用示例

创建你自己的`Storage`，内嵌`*sql_based_storage.SqlBasedStorage`，只覆写必要的方法 ：

```go
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
```

然后实现自己的`sql`方言，可以从`SqlProvider`从头开始写，也可以内嵌`Sql97Provider`只覆写必要的方法：

```go
type FooSqlProvider struct {
	sql_based_storage.Sql97Provider
}

var _ sql_based_storage.SqlProvider = &FooSqlProvider{}

func (x *FooSqlProvider) NowTimestampSql(ctx context.Context, tableFullName string) (string, []any) {
	return "SELECT NOW()", nil
}
```

最后使用的话：

```go
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
```



