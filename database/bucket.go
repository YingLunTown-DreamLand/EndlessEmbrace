package database

import (
	"fmt"
	"strings"

	"go.etcd.io/bbolt"
)

// attachBucket 获取 b.path 处的存储桶，
// 并通过回调函数 function 检索此存储桶。
//
// 值得注意的是，attachBucket 是内部实现细节，
// 任何外部调用不应也不得调用此函数。
// 除此外，attachBucket 也不会进行权限检查。
//
// flag 指示操作类型，只可能为下列之一。
//   - 只读 (0)
//   - 只写 (1)
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
func (b *Bucket) attachBucket(callerName string, function func(*bbolt.Bucket) error, flag int) error {
	if b.db == nil || *b.db == nil {
		return fmt.Errorf("%s: attachBucket: use of closed upstream database", callerName)
	}

	f := func(tx *bbolt.Tx) error {
		if len(b.path) == 1 {
			return fmt.Errorf("%s: attachBucket: Invalid operation (Try to access root as bucket)", callerName)
		}
		bucket := tx.Bucket(b.path[1])
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

// Path 返回该存储桶的绝对路径
func (b *Bucket) Path() []string {
	pathString := make([]string, len(b.path))
	for index, value := range b.path {
		pathString[index] = string(value)
	}
	return pathString
}

// PathString 返回该存储桶绝对路径的字符串形式。
// 子桶与子桶之间用 / 连接
func (b *Bucket) PathString() string {
	return strings.Join(b.Path(), "/")
}

// StillAlive 检查当前存储桶是否仍然在数据库中有效
func (b *Bucket) StillAlive() bool {
	if b.db == nil || *b.db == nil {
		return false
	}

	if len(b.path) == 1 {
		return true
	}

	err := b.attachBucket(
		"StillAlive",
		func(b *bbolt.Bucket) error {
			return nil
		},
		DatabaseFlagReadOnly,
	)
	return (err == nil)
}

// 返回当前存储桶中已有键的名称的集合。
// 如果试图在根目录上调用此函数，
// 或当前存储桶已经无效，
// 或当前存储桶不具备读权限，
// 则永远返回非 nil 的空切片
func (b *Bucket) GetKeyMapping() (mapping [][]byte) {
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

// 创建名为 name 的存储桶。
// 必须确保当前存储桶具有写权限。
// 如果名为 name 的存储桶已经存在，
// 则亦不会返回错误
func (b *Bucket) CreateBucket(name []byte) error {
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

	if !b.HasBucket(name) {
		return b.attachBucket(
			"CreateBucket",
			func(b *bbolt.Bucket) error {
				_, err := b.CreateBucket(name)
				if err != nil {
					return fmt.Errorf("CreateBucket: %v", err)
				}
				return nil
			},
			DatabaseFlagReadWrite,
		)
	}
	return nil
}

// 删除名为 name 的存储桶。
// 必须确保当前存储桶具有写权限。
// 如果名为 name 的存储桶不存在，
// 则亦不会返回错误
func (b *Bucket) DeleteBucket(name []byte) error {
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
			err := b.DeleteBucket(name)
			if err != nil {
				return fmt.Errorf("DeleteBucket: %v", err)
			}
			return nil
		},
		DatabaseFlagReadWrite,
	)
}

// 向名为 key 的键处放置数据 data。
// 必须确保当前存储桶具有写权限。
// 试图在根目录上调用此函数不会产生任何效果
func (b *Bucket) PutData(key []byte, data []byte) error {
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
		db:         b.db,
		path:       append(b.path, name),
		permission: permission & b.permission,
	}
}

// 从名为 key 的键处读取数据。
// 必须确保当前存储桶具有读权限。
// 试图在根目录上调用此函数不会产生任何效果
func (b *Bucket) GetData(key []byte) (data []byte, err error) {
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

// CloseDatabase 将数据库关闭。
//
// CloseDatabase 可以从任意具
// 有关闭数据库权限的子桶上调用。
//
// 可以存在多个协程调用此函数，
// 这不会对数据库造成错误的影响
func (b *Bucket) CloseDatabase() error {
	if b.permission&BucketPermissionCloseDatabase == 0 {
		return fmt.Errorf("CloseDatabase: Permission denial")
	}

	if b.db == nil || *b.db == nil {
		return nil
	}

	if !b.StillAlive() {
		return fmt.Errorf("CloseDatabase: Current bucket is dead")
	}

	err := (*b.db).Close()
	if err != nil {
		return fmt.Errorf("CloseDatabase: %v", err)
	}
	*b.db = nil

	return nil
}
