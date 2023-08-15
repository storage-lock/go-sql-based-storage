package sql_based_storage

import "fmt"

var (
	ErrSqlProviderCanNotNil       = fmt.Errorf("SqlProvider can not nil")
	ErrConnectionManagerCanNotNil = fmt.Errorf("ConnectionManager can not nil")
	ErrTableFullNameCanNotEmpty   = fmt.Errorf("TableFullName can not empty")
)
