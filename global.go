package sqlx_locks

import (
	"github.com/jmoiron/sqlx"
)

var GlobalSqlxLockFactory *SqlxLockFactory

// InitGlobalSqlxLockFactory 初始化全局的SqlxLockFactory
func InitGlobalSqlxLockFactory(db *sqlx.DB) error {
	factory, err := NewSqlxLockFactory(db)
	if err != nil {
		return err
	}
	GlobalSqlxLockFactory = factory
	return nil
}
