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
