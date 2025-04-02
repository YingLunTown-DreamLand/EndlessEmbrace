package database

import (
	"fmt"

	"go.etcd.io/bbolt"
)

// ViewBucketWithFunc 以只读模式检索当前存储桶。
// function 指示用于检索该存储桶的函数
func (b *Bucket) ViewBucketWithFunc(function func(b *bbolt.Bucket) error) error {
	if b.permission&BucketPermissionReadOnly == 0 {
		return fmt.Errorf("ViewBucketWithFunc: Permission denial")
	}
	return b.attachBucket("ViewBucketWithFunc", function, DatabaseFlagReadOnly)
}

// WriteBucketWithFunc 以读写模式检索当前存储桶。
// function 指示用于检索该存储桶的函数
func (b *Bucket) WriteBucketWithFunc(function func(b *bbolt.Bucket) error) error {
	if b.permission&BucketPermissionReadOnly == 0 || b.permission&BucketPermissionWriteOnly == 0 {
		return fmt.Errorf("WriteBucketWithFunc: Permission denial")
	}
	return b.attachBucket("WriteBucketWithFunc", function, DatabaseFlagReadWrite)
}

// ViewRootWithFunc 以只读模式检索当前数据库。
// function 指示用于检索这个数据库的函数
func (b *Bucket) ViewRootWithFunc(function func(t *bbolt.Tx) error) error {
	if b.permission&BucketPermissionReadOnly == 0 {
		return fmt.Errorf("ViewRootWithFunc: Permission denial")
	}
	err := (*b.db).View(function)
	if err != nil {
		return fmt.Errorf("ViewRootWithFunc: %v", err)
	}
	return nil
}

// WriteRootWithFunc 以读写模式检索当前数据库。
// function 指示用于检索这个数据库的函数
func (b *Bucket) WriteRootWithFunc(function func(t *bbolt.Tx) error) error {
	if b.permission&BucketPermissionReadOnly == 0 || b.permission&BucketPermissionWriteOnly == 0 {
		return fmt.Errorf("WriteRootWithFunc: Permission denial")
	}
	err := (*b.db).Update(function)
	if err != nil {
		return fmt.Errorf("WriteRootWithFunc: %v", err)
	}
	return nil
}
