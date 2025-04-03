package database

import "strings"

// Path 返回该存储桶的绝对路径
func (b *bucket) Path() []string {
	pathString := make([]string, len(b.path))
	for index, value := range b.path {
		pathString[index] = string(value)
	}
	return pathString
}

// PathString 返回该存储桶绝对路径的字符串形式。
// 子桶与子桶之间用 / 连接
func (b *bucket) PathString() string {
	return strings.Join(b.Path(), "/")
}

// StillAlive 检查当前存储桶是否仍然在数据库中有效
func (b *bucket) StillAlive() bool {
	b.locker.RLock()
	defer b.locker.RUnlock()
	return b.stillAlive()
}
