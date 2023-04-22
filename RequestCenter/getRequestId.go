package RequestCenter

import "sync/atomic"

func (r *Resources) GetCurrentRequestId() uint32 {
	return atomic.LoadUint32(&r.requestId)
}

func (r *Resources) GetNewRequestId() uint32 {
	return atomic.AddUint32(&r.requestId, 1)
}
