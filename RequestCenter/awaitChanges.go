package RequestCenter

import "fmt"

// 等待 go-cqhttp 返回请求 ID 为 key 的请求，
// 并在等到后从 r.request.datas 删除此请求对应的互斥锁
func (r *Resources) AwaitChangesAndDeleteRequest(key uint32) error {
	if !r.TestReqest(key) {
		return fmt.Errorf("AwaitChanges: %v is not existed", key)
	}
	// if key is not existed
	r.request.lockDown.RLock()
	tmp := r.request.datas[RequestId(key)]
	r.request.lockDown.RUnlock()
	// get lock
	tmp.Lock()
	tmp.Unlock()
	// await changes
	err := r.DeleteRequest(key)
	if err != nil {
		return err
	}
	// delete request
	return nil
	// return
}
