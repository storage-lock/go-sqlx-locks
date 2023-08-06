package sqlx_locks

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	sqldb_storage "github.com/storage-lock/go-sqldb-storage"
	"github.com/storage-lock/go-storage"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
)

type SqlxLockFactory struct {
	db *sqlx.DB
	*storage_lock_factory.StorageLockFactory[*sql.DB]
}

func NewSqlxLockFactory(db *sqlx.DB) (*SqlxLockFactory, error) {
	connectionManager := NewSqlxConnectionManager(db)

	storage, err := CreateStorageForSqlxDb(db, connectionManager)
	if err != nil {
		return nil, err
	}

	factory := storage_lock_factory.NewStorageLockFactory[*sql.DB](storage, connectionManager)

	return &SqlxLockFactory{
		db:                 db,
		StorageLockFactory: factory,
	}, nil
}

// CreateStorageForSqlxDb 尝试从sqlx创建Storage
func CreateStorageForSqlxDb(db *sqlx.DB, connectionManager storage.ConnectionManager[*sql.DB]) (storage.Storage, error) {

	// 先尝试根据驱动名称创建
	storage, err := sqldb_storage.NewStorageByDriverName(db.DriverName(), connectionManager)
	if storage != nil && err == nil {
		return storage, err
	}

	// 再然后根据识别出来的名称创建
	return sqldb_storage.NewStorageBySqlDb(db.DB, connectionManager)
}
