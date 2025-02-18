package RequestCenter

import "sync"

type RequestId uint64 // 请求 ID

// 用于描述客户端的资源
type Resources struct {
	// 放置客户端各个请求对应的互斥锁，
	// 用于阻塞以等待 go-cqhttp 回应请求
	request struct {
		lockDown sync.RWMutex
		datas    map[RequestId]*sync.Mutex
	}
	// 放置 go-cqhttp 对“客户端请求”的返回值
	responce struct {
		lockDown sync.RWMutex
		datas    map[RequestId]Responce
	}
	// 存储当前已累计的请求 ID
	requestId uint32
}

// 用于描述单个请求
type Request struct {
	Action    string      `json:"action"` // 终结点名称
	Params    interface{} `json:"params"` // ...
	RequestId string      `json:"echo"`   // 请求 ID
}

// 用于描述单个请求的返回值
type Responce struct {
	Status     string                 `json:"status"`  // 请求结果
	Retcode    int64                  `json:"retcode"` // ...
	Msg        *string                `json:"msg"`     // 错误消息(仅在 API 调用失败时有该字段)
	Wording    *string                `json:"wording"` // 对错误的详细解释(中文；仅在 API 调用失败时有该字段)
	Data       map[string]interface{} `json:"data"`    // ...
	ResponceId string                 `json:"echo"`    // 对应请求中的请求 ID
}
