package database

import (
	"fmt"
	"sync"

	"go.etcd.io/bbolt"
)

// 打开 path 处的数据库。
// 如果数据库不存在，则将会创建一个
func OpenOrCreateDatabase(path string) (
	database *Bucket,
	err error,
) {
	db, err := bbolt.Open(
		path,
		0600,
		&bbolt.Options{
			Timeout:      0,
			NoGrowSync:   false,
			FreelistType: bbolt.FreelistMapType,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("OpenOrCreateDatabase: %v", err)
	}
	return &Bucket{
		locker:     new(sync.RWMutex),
		db:         &db,
		path:       [][]byte{[]byte("root")},
		permission: BucketPermissionReadOnly | BucketPermissionWriteOnly | BucketPermissionCloseDatabase,
	}, nil
}

// CloseDatabase 将数据库关闭。
//
// CloseDatabase 可以从任意具
// 有关闭数据库权限的子桶上调用。
//
// 可以存在多个协程调用此函数，
// 这不会对数据库造成错误的影响
func CloseDatabase(database *Bucket) error {
	database.locker.Lock()
	defer database.locker.Unlock()

	if database.permission&BucketPermissionCloseDatabase == 0 {
		return fmt.Errorf("CloseDatabase: Permission denial")
	}

	if database.db == nil || *database.db == nil {
		return nil
	}

	if !database.stillAlive() {
		return fmt.Errorf("CloseDatabase: Current bucket is dead")
	}

	err := (*database.db).Close()
	if err != nil {
		return fmt.Errorf("CloseDatabase: %v", err)
	}
	*database.db = nil

	return nil
}
