package RequestCenter

import (
	"fmt"
)

// 测定 key 是否在 r.responce.datas 中。
// 如果存在，返回真，否则返回假
func (r *Resources) TestResponce(key uint32) bool {
	r.responce.lockDown.RLock()
	defer r.responce.lockDown.RUnlock()
	// init values
	_, ok := r.responce.datas[RequestId(key)]
	return ok
	// return
}

// 向 r.responce.datas[key] 写入请求的返回值并释放对应的互斥锁
func (r *Resources) WriteResponce(key uint32, resp Responce) error {
	if r.TestResponce(key) {
		return fmt.Errorf("WriteResponce: %v is already existed", key)
	}
	// if key is already existed
	r.responce.lockDown.Lock()
	defer r.responce.lockDown.Unlock()
	// init values
	r.responce.datas[RequestId(key)] = resp
	r.request.datas[RequestId(key)].Unlock()
	// write responce and release the lock
	return nil
	// return
}

// 从 r.responce.datas 提取请求 ID 为 key 的请求的返回值并从 r.responce.datas 中删除它
func (r *Resources) LoadResponceAndDelete(key uint32) (Responce, error) {
	if !r.TestResponce(key) {
		return Responce{}, fmt.Errorf("LoadResponceAndDelete: %v is not existed", key)
	}
	// if key is not existed
	r.responce.lockDown.Lock()
	defer r.responce.lockDown.Unlock()
	// init values
	ans := r.responce.datas[RequestId(key)]
	// get responce
	delete(r.responce.datas, RequestId(key))
	newMap := map[RequestId]Responce{}
	for key, value := range r.responce.datas {
		newMap[key] = value
	}
	r.responce.datas = newMap
	// delete values
	return ans, nil
	// return
}
