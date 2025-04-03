package database

import (
	"fmt"

	"go.etcd.io/bbolt"
)

// attachBucket 获取 b.path 处的存储桶，
// 并通过回调函数 function 检索此存储桶。
//
// 值得注意的是，attachBucket 是内部实现细节，
// 任何外部调用不应也不得调用此函数。
//
// 除此外，attachBucket 也不会进行权限检查，
// 也不确保一切操作是线程安全的，
// 这意味着锁操作需要由调用者控制。
//
// flag 指示操作类型，只可能为下列之一。
//   - 只读 (0)
//   - 读写 (1)
//
// callerName 指示调用者的名字，
// 这只用作 debug 的用途。
//
// 另外，attachBucket 不被允许访问数据库的根目录，
// 因为数据库的根不是存储桶
// (虽然它可以被视作一种特殊的存储桶)
//
// 如果一切正常，attachBucket 按原样返回
// function 所返回的错误
func (b *bucket) attachBucket(callerName string, function func(*bbolt.Bucket) error, flag int) error {
	if b.db == nil || *b.db == nil {
		return fmt.Errorf("%s: attachBucket: use of closed upstream database", callerName)
	}

	f := func(tx *bbolt.Tx) error {
		if len(b.path) == 1 {
			return fmt.Errorf("%s: attachBucket: Invalid operation (Try to access root as bucket)", callerName)
		}

		bucket := tx.Bucket(b.path[1])
		if bucket == nil {
			return fmt.Errorf("%s: attachBucket: Target bucket (path=%s) is not reachable", callerName, b.PathString())
		}

		for i := 2; i < len(b.path); i++ {
			bucket = bucket.Bucket(b.path[i])
			if bucket == nil {
				return fmt.Errorf("%s: attachBucket: Target bucket (path=%s) is not reachable", callerName, b.PathString())
			}
		}
		return function(bucket)
	}

	switch flag {
	case DatabaseFlagReadOnly:
		return (*b.db).View(f)
	case DatabaseFlagReadWrite:
		return (*b.db).Update(f)
	default:
		return fmt.Errorf("%s: attachBucket: Invalid flag %d was get, expected 0 (read only) or 1 (read/write)", callerName, flag)
	}
}

// StillAlive 检查当前存储桶是否仍然在数据库中有效。
// 它是不设有线程安全保护的内部实现细节
func (b *bucket) stillAlive() bool {
	if b.db == nil || *b.db == nil {
		return false
	}

	if len(b.path) == 1 {
		return true
	}

	err := b.attachBucket(
		"stillAlive",
		func(b *bbolt.Bucket) error {
			return nil
		},
		DatabaseFlagReadOnly,
	)
	return (err == nil)
}
