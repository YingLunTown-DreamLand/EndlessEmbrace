package database

import (
	"fmt"

	"go.etcd.io/bbolt"
)

// 创建名为 name 的存储桶。
// 必须确保当前存储桶具有写权限。
// 如果名为 name 的存储桶已经存在，
// 则亦不会返回错误
func (b *Bucket) CreateBucket(name []byte) error {
	b.locker.Lock()
	defer b.locker.Unlock()

	if b.db == nil || *b.db == nil {
		return fmt.Errorf("CreateBucket: use of closed upstream database")
	}

	if b.permission&BucketPermissionWriteOnly == 0 {
		return fmt.Errorf("CreateBucket: Permission denial")
	}

	if len(b.path) == 1 {
		err := (*b.db).Update(func(tx *bbolt.Tx) error {
			if tx.Bucket(name) != nil {
				return nil
			}
			_, err := tx.CreateBucket(name)
			return err
		})
		if err != nil {
			return fmt.Errorf("CreateBucket: %v", err)
		}
		return nil
	}

	return b.attachBucket(
		"CreateBucket",
		func(b *bbolt.Bucket) error {
			if b.Bucket(name) != nil {
				return nil
			}
			_, err := b.CreateBucket(name)
			if err != nil {
				return fmt.Errorf("CreateBucket: %v", err)
			}
			return nil
		},
		DatabaseFlagReadWrite,
	)
}

// 删除名为 name 的存储桶。
// 必须确保当前存储桶具有写权限。
// 如果名为 name 的存储桶不存在，
// 则亦不会返回错误
func (b *Bucket) DeleteBucket(name []byte) error {
	b.locker.Lock()
	defer b.locker.Unlock()

	if b.db == nil || *b.db == nil {
		return fmt.Errorf("DeleteBucket: use of closed upstream database")
	}

	if b.permission&BucketPermissionWriteOnly == 0 {
		return fmt.Errorf("DeleteBucket: Permission denial")
	}

	if len(b.path) == 1 {
		err := (*b.db).Update(func(tx *bbolt.Tx) error {
			if tx.Bucket(name) == nil {
				return nil
			}
			return tx.DeleteBucket(name)
		})
		if err != nil {
			return fmt.Errorf("DeleteBucket: %v", err)
		}
		return nil
	}

	return b.attachBucket(
		"DeleteBucket",
		func(b *bbolt.Bucket) error {
			if b.Bucket(name) == nil {
				return nil
			}
			err := b.DeleteBucket(name)
			if err != nil {
				return fmt.Errorf("DeleteBucket: %v", err)
			}
			return nil
		},
		DatabaseFlagReadWrite,
	)
}

// 从当前存储桶上获取名为 name 的子桶。
//
// 如果名为 name 的子桶不存在，
// 或当前存储桶不具备读权限，
// 或数据库已被关闭，
// 则返回空值。
//
// permission 指示子桶可以得到的权限，
// 只可能是下列之一。
//   - 0 (不具备任何权限)
//   - 1 (只读)
//   - 2 (只写)
//   - 3 (读写)
//   - 4 (可以关闭数据库)
//   - 5 (可以关闭数据库, 只读)
//   - 6 (可以关闭数据库, 只写)
//   - 7 (可以关闭数据库, 读写)
//
// permission 可以指定比现有权限更大的权限，
// 但所得子桶生效的权限将只能是现有权限的子集。
// 例如，如果当前权限为 3，而 permission 为 7，
// 那么最终所获子桶的权限也仍然是 3。
//
// 另外，在数据库根目录处亦可调用此函数
func (b *Bucket) GetSubBucketByName(name []byte, permission int) (result *Bucket) {
	if b.permission&BucketPermissionReadOnly == 0 {
		return nil
	}
	if !b.HasBucket(name) {
		return nil
	}
	return &Bucket{
		locker:     b.locker,
		db:         b.db,
		path:       append(b.path, name),
		permission: permission & b.permission,
	}
}

// 向名为 key 的键处放置数据 data。
// 必须确保当前存储桶具有写权限。
// 试图在根目录上调用此函数不会产生任何效果
func (b *Bucket) PutData(key []byte, data []byte) error {
	b.locker.Lock()
	defer b.locker.Unlock()

	if b.permission&BucketPermissionWriteOnly == 0 {
		return fmt.Errorf("PutData: Permission denial")
	}

	if len(b.path) == 1 {
		return nil
	}

	return b.attachBucket(
		"PutData",
		func(b *bbolt.Bucket) error {
			err := b.Put(key, data)
			if err != nil {
				return fmt.Errorf("PutData: %v", err)
			}
			return nil
		},
		DatabaseFlagReadWrite,
	)
}

// 删除名为 key 的键处所存放的数据。
// 必须确保当前存储桶具有写权限。
// 试图在根目录上调用此函数不会产生任何效果
func (b *Bucket) DeleteData(key []byte) error {
	b.locker.Lock()
	defer b.locker.Unlock()

	if b.permission&BucketPermissionWriteOnly == 0 {
		return fmt.Errorf("PutData: Permission denial")
	}

	if len(b.path) == 1 {
		return nil
	}

	return b.attachBucket(
		"DeleteData",
		func(b *bbolt.Bucket) error {
			err := b.Delete(key)
			if err != nil {
				return fmt.Errorf("DeleteData: %v", err)
			}
			return nil
		},
		DatabaseFlagReadWrite,
	)
}

// 从名为 key 的键处读取数据。
// 必须确保当前存储桶具有读权限。
// 试图在根目录上调用此函数不会产生任何效果
func (b *Bucket) GetData(key []byte) (data []byte, err error) {
	b.locker.RLock()
	defer b.locker.RUnlock()

	if b.permission&BucketPermissionReadOnly == 0 {
		return nil, fmt.Errorf("GetData: Permission denial")
	}

	if len(b.path) == 1 {
		return nil, nil
	}

	err = b.attachBucket(
		"GetData",
		func(b *bbolt.Bucket) error {
			result := b.Get(key)
			data = make([]byte, len(result))
			copy(data, result)
			return nil
		},
		DatabaseFlagReadOnly,
	)
	return
}
