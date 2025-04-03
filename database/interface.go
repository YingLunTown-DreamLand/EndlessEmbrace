package database

import "go.etcd.io/bbolt"

// StatesHolder 持有数据库的状态信息
type StatesHolder interface {
	// Path 返回该存储桶的绝对路径
	Path() []string
	// PathString 返回该存储桶绝对路径的字符串形式。
	// 子桶与子桶之间用 / 连接
	PathString() string
	// StillAlive 检查当前存储桶是否仍然在数据库中有效
	StillAlive() bool
}

// MappingHolder 持有数据库的映射信息，
// 用于查找或确定数据库中具有哪些存储桶
// 或键信息
type MappingHolder interface {
	// 返回当前存储桶中已有键的名称的集合。
	// 如果试图在根目录上调用此函数，
	// 或当前存储桶已经无效，
	// 或当前存储桶不具备读权限，
	// 则永远返回非 nil 的空切片
	GetKeyMapping() (mapping [][]byte)
	// 返回当前存储桶中已有存储桶的名称的集合。
	// 如果当前存储桶已经无效，
	// 或数据库已被关闭，
	// 或当前存储桶不具备读权限，
	// 则永远返回非 nil 的空切片
	GetBucketMapping() (mapping [][]byte)
	// 确定数据库中是否存在名为 name 的键。
	// 如果当前存储桶已经无效，
	// 或尝试在根目录上调用此函数，
	// 或当前存储桶不具备读权限，
	// 亦会永远返回假
	HasKey(name []byte) (has bool)
	// 确定数据库中是否存在名为 name 的存储桶。
	// 如果当前存储桶已经无效，
	// 或数据库已被关闭，
	// 或当前存储桶不具备读权限，
	// 则永远返回假
	HasBucket(name []byte) (has bool)
}

// BasicOperationHolder 持有数据库的基本操作实现，
// 例如存储桶或数据的创建、删除和存储
type BasicOperationHolder interface {
	// 创建名为 name 的存储桶。
	// 必须确保当前存储桶具有写权限。
	// 如果名为 name 的存储桶已经存在，
	// 则亦不会返回错误
	CreateBucket(name []byte) error
	// 删除名为 name 的存储桶。
	// 必须确保当前存储桶具有写权限。
	// 如果名为 name 的存储桶不存在，
	// 则亦不会返回错误
	DeleteBucket(name []byte) error
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
	GetSubBucketByName(name []byte, permission int) (result *Bucket)
	// 向名为 key 的键处放置数据 data。
	// 必须确保当前存储桶具有写权限。
	// 试图在根目录上调用此函数不会产生任何效果
	PutData(key []byte, data []byte) error
	// 删除名为 key 的键处所存放的数据。
	// 必须确保当前存储桶具有写权限。
	// 试图在根目录上调用此函数不会产生任何效果
	DeleteData(key []byte) error
	// 从名为 key 的键处读取数据。
	// 必须确保当前存储桶具有读权限。
	// 试图在根目录上调用此函数不会产生任何效果
	GetData(key []byte) (data []byte, err error)
}

type HighLevelOperationHolder interface {
	// ViewBucketWithFunc 以只读模式检索当前存储桶。
	// function 指示用于检索该存储桶的函数
	ViewBucketWithFunc(function func(b *bbolt.Bucket) error) error
	// WriteBucketWithFunc 以读写模式检索当前存储桶。
	// function 指示用于检索该存储桶的函数
	WriteBucketWithFunc(function func(b *bbolt.Bucket) error) error
	// ViewDBWithFunc 以只读模式检索当前数据库。
	// function 指示用于检索这个数据库的函数
	ViewDBWithFunc(function func(t *bbolt.Tx) error) error
	// WriteDBWithFunc 以读写模式检索当前数据库。
	// function 指示用于检索这个数据库的函数
	WriteDBWithFunc(function func(t *bbolt.Tx) error) error
}
