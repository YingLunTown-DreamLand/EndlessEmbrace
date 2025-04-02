package database

import "go.etcd.io/bbolt"

// 返回当前存储桶中已有键的名称的集合。
// 如果试图在根目录上调用此函数，
// 或当前存储桶已经无效，
// 或当前存储桶不具备读权限，
// 则永远返回非 nil 的空切片
func (b *Bucket) GetKeyMapping() (mapping [][]byte) {
	b.locker.RLock()
	defer b.locker.RUnlock()

	mapping = make([][]byte, 0)

	if len(b.path) == 1 || b.permission&BucketPermissionReadOnly == 0 {
		return
	}

	_ = b.attachBucket(
		"GetKeyMapping",
		func(b *bbolt.Bucket) error {
			return b.ForEach(func(k, v []byte) error {
				mapping = append(mapping, k)
				return nil
			})
		},
		DatabaseFlagReadOnly,
	)
	return
}

// 返回当前存储桶中已有存储桶的名称的集合。
// 如果当前存储桶已经无效，
// 或数据库已被关闭，
// 或当前存储桶不具备读权限，
// 则永远返回非 nil 的空切片
func (b *Bucket) GetBucketMapping() (mapping [][]byte) {
	b.locker.RLock()
	defer b.locker.RUnlock()

	mapping = make([][]byte, 0)

	if b.db == nil || *b.db == nil || b.permission&BucketPermissionReadOnly == 0 {
		return
	}

	if len(b.path) == 1 {
		_ = (*b.db).View(func(tx *bbolt.Tx) error {
			return tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
				mapping = append(mapping, name)
				return nil
			})
		})
		return
	}

	_ = b.attachBucket(
		"GetBucketMapping",
		func(b *bbolt.Bucket) error {
			return b.ForEachBucket(func(k []byte) error {
				mapping = append(mapping, k)
				return nil
			})
		},
		DatabaseFlagReadOnly,
	)
	return
}

// 确定数据库中是否存在名为 name 的键。
// 如果当前存储桶已经无效，
// 或尝试在根目录上调用此函数，
// 或当前存储桶不具备读权限，
// 亦会永远返回假
func (b *Bucket) HasKey(name []byte) (has bool) {
	b.locker.RLock()
	defer b.locker.RUnlock()

	if len(b.path) == 1 || b.permission&BucketPermissionReadOnly == 0 {
		return false
	}

	_ = b.attachBucket(
		"HasKey",
		func(b *bbolt.Bucket) error {
			has = (b.Bucket(name) != nil)
			return nil
		},
		DatabaseFlagReadOnly,
	)
	return
}

// 确定数据库中是否存在名为 name 的存储桶。
// 如果当前存储桶已经无效，
// 或数据库已被关闭，
// 或当前存储桶不具备读权限，
// 则永远返回假
func (b *Bucket) HasBucket(name []byte) (has bool) {
	b.locker.RLock()
	defer b.locker.RUnlock()

	if b.db == nil || *b.db == nil || b.permission&BucketPermissionReadOnly == 0 {
		return
	}

	if len(b.path) == 1 {
		_ = (*b.db).View(func(tx *bbolt.Tx) error {
			has = (tx.Bucket(name) != nil)
			return nil
		})
		return
	}

	_ = b.attachBucket(
		"HasBucket",
		func(b *bbolt.Bucket) error {
			has = (b.Bucket(name) != nil)
			return nil
		},
		DatabaseFlagReadOnly,
	)
	return
}
