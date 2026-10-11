package sqlx_locks

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func testSqlxDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", "file:sqlx_test.db?mode=memory")
	assert.Nil(t, err)
	t.Cleanup(func() { db.Close() })
	return sqlx.NewDb(db, "sqlite3")
}

func TestNewSqlxConnectionManager(t *testing.T) {
	db := testSqlxDB(t)
	m := NewSqlxConnectionManager(db)
	assert.NotNil(t, m)
	assert.Equal(t, SqlxConnectionManagerName, m.Name())

	ctx := context.Background()
	got, err := m.Take(ctx)
	assert.Nil(t, err)
	assert.Equal(t, db.DB, got)

	assert.Nil(t, m.Return(ctx, got))
	assert.Nil(t, m.Shutdown(ctx))
}

func TestSqlxConnectionManager_ShutdownNil(t *testing.T) {
	m := NewSqlxConnectionManager(nil)
	assert.Nil(t, m.Shutdown(context.Background()))
}

func TestNewSqlxLockFactory(t *testing.T) {
	db := testSqlxDB(t)
	factory, err := NewSqlxLockFactory(db)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 端到端加解锁
	ctx := context.Background()
	lock, err := factory.CreateLock("test-sqlx-lock")
	assert.Nil(t, err)
	assert.NotNil(t, lock)
	assert.Nil(t, lock.Lock(ctx, "owner-1"))
	assert.Nil(t, lock.UnLock(ctx, "owner-1"))
}

func TestCreateStorageForSqlxDb_Errors(t *testing.T) {
	// sqlmock 驱动：第一条路驱动名不对 + 回退驱动识别失败
	mockDB, _, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()
	mockSqlx := sqlx.NewDb(mockDB, "sqlmock")

	m := NewSqlxConnectionManager(mockSqlx)
	_, err = CreateStorageForSqlxDb(mockSqlx, m)
	assert.NotNil(t, err)

	_, err = NewSqlxLockFactory(mockSqlx)
	assert.NotNil(t, err)
}

func TestInitGlobalSqlxLockFactory(t *testing.T) {
	db := testSqlxDB(t)
	assert.Nil(t, InitGlobalSqlxLockFactory(db))
	assert.NotNil(t, GlobalSqlxLockFactory)
	GlobalSqlxLockFactory = nil

	mockDB, _, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()
	assert.NotNil(t, InitGlobalSqlxLockFactory(sqlx.NewDb(mockDB, "sqlmock")))
	assert.Nil(t, GlobalSqlxLockFactory)
}
