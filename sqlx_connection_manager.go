package sqlx_locks

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/storage-lock/go-storage"
)

// SqlxConnectionManager 从sqlx（https://github.com/jmoiron/sqlx）中复用数据库连接，如果当前项目是引入的sqlx的话则可以与其复用数据库连接资源
// TODO 2023-8-4 01:32:17 单元测试
type SqlxConnectionManager struct {
	db *sqlx.DB
}

var _ storage.ConnectionManager[*sql.DB] = &SqlxConnectionManager{}

func NewSqlxConnectionManager(db *sqlx.DB) *SqlxConnectionManager {
	return &SqlxConnectionManager{
		db: db,
	}
}

const SqlxConnectionManagerName = "sqlx-connection-manager"

func (x *SqlxConnectionManager) Name() string {
	return SqlxConnectionManagerName
}

func (x *SqlxConnectionManager) Take(ctx context.Context) (*sql.DB, error) {
	return x.db.DB, nil
}

func (x *SqlxConnectionManager) Return(ctx context.Context, db *sql.DB) error {
	return nil
}

func (x *SqlxConnectionManager) Shutdown(ctx context.Context) error {
	if x.db != nil {
		return x.db.Close()
	}
	return nil
}
