package RequestCenter

import "sync"

var isAlreadyInited bool = false // 用于标识 Resources 结构体是否已经被初始化

// 初始化 Resources 结构体以供客户端使用。
// 此函数只应该在客户端启动时被调用，因为 Resources 作为资源应当只能存在一个。
// 重复调用此函数将会导致程序惊慌
func (r *Resources) InitStruct() {
	if isAlreadyInited {
		panic("InitStruct: Struct Resources has been initialized")
	}
	// if the struct has been initialized
	r.request = struct {
		lockDown sync.RWMutex
		datas    map[RequestId]*sync.Mutex
	}{
		lockDown: sync.RWMutex{},
		datas:    make(map[RequestId]*sync.Mutex),
	}
	// r.request
	r.responce = struct {
		lockDown sync.RWMutex
		datas    map[RequestId]Responce
	}{
		lockDown: sync.RWMutex{},
		datas:    make(map[RequestId]Responce),
	}
	// r.responce
	r.requestId = 0
	// r.rquestId
}
