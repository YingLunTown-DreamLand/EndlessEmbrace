package RequestCenter

import (
	"fmt"
	"sync"
)

// 测定 key 是否在 r.request.datas 中。
// 如果存在，返回真，否则返回假
func (r *Resources) TestReqest(key uint32) bool {
	r.request.lockDown.RLock()
	defer r.request.lockDown.RUnlock()
	// init values
	if r.request.datas[RequestId(key)] == nil {
		return false
	} else {
		return true
	}
	// return
}

// 向 r.request.datas 写入请求 ID 为 key 的请求并锁定对应的互斥锁
func (r *Resources) WriteRequest(key uint32) error {
	if r.TestReqest(key) {
		return fmt.Errorf("WriteRequest: %v is already existed", key)
	}
	// if key is already existed
	r.request.lockDown.Lock()
	defer r.request.lockDown.Unlock()
	// init values
	r.request.datas[RequestId(key)] = &sync.Mutex{}
	r.request.datas[RequestId(key)].Lock()
	// lock down
	return nil
	// return
}

// 从 r.request.datas 删除请求 ID 为 key 的请求
func (r *Resources) DeleteRequest(key uint32) error {
	if !r.TestReqest(key) {
		return fmt.Errorf("WriteRequest: %v is not existed", key)
	}
	// if key is not existed
	r.request.lockDown.Lock()
	defer r.request.lockDown.Unlock()
	// init values
	delete(r.request.datas, RequestId(key))
	newMap := map[RequestId]*sync.Mutex{}
	for key, value := range r.request.datas {
		newMap[key] = value
	}
	r.request.datas = newMap
	// delete values
	return nil
	// return
}
