package database

import (
	"sync"

	"go.etcd.io/bbolt"
)

const (
	DatabaseFlagReadOnly = iota
	DatabaseFlagReadWrite
)

const (
	BucketPermissionReadOnly = 1 << iota
	BucketPermissionWriteOnly
	BucketPermissionCloseDatabase
)

// bucket 描述数据库中的单个存储桶。
//
// 特别地，将数据库的根目录
// 也视作为一个存储桶，
// 但不具备检索和存放键值的功能。
//
// bucket 是一个内部实现细节，
// 不应被外部使用
type bucket struct {
	locker     *sync.RWMutex // 全局互斥锁
	db         **bbolt.DB    // 底层(上层)数据库
	path       [][]byte      // 当前子桶的绝对路径
	permission int           // 当前子桶的权限级别
}
