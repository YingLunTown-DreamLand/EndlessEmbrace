package database

// Bucket 描述数据库中的存储桶，
// 亦可用于指示数据库的根目录
type Bucket struct {
	b *bucket
}

// StatesHolder 返回可用于获取
// 数据库状态信息的接口
func (b *Bucket) StatesHolder() StatesHolder {
	return b.b
}

// MappingHolder 返回可用于获取
// 数据库映射状态信息的接口
func (b *Bucket) MappingHolder() MappingHolder {
	return b.b
}

// BasicOperationHolder 返回可用于获取
// 数据库基本操作的接口
func (b *Bucket) BasicOperationHolder() BasicOperationHolder {
	return b.b
}

// HighLevelOperationHolder 返回可用于获取
// 数据库高级操作的接口
func (b *Bucket) HighLevelOperationHolder() HighLevelOperationHolder {
	return b.b
}
