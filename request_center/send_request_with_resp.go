package RequestCenter

import (
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

// 向 go-cqhttp 发送请求但无视返回值
func (r *Resources) SendRequest(
	conn *websocket.Conn,
	request Request,
) error {
	if err := conn.WriteJSON(request); err != nil {
		return fmt.Errorf("SendRequest: %v", err)
	}
	// send packet
	return nil
	// return
}

// 向 go-cqhttp 发送请求并获取返回值
func (r *Resources) SendRequestWithResponce(
	conn *websocket.Conn,
	request Request,
) (Responce, error) {
	requestId, err := strconv.ParseUint(request.RequestId, 10, 32)
	if err != nil {
		return Responce{}, fmt.Errorf("SendRequestWithResponce: %v", err)
	}
	r.WriteRequest(uint32(requestId))
	// write request and lock down r.request.datas[RequestId(requestId)]
	err = conn.WriteJSON(request)
	if err != nil {
		return Responce{}, fmt.Errorf("SendRequestWithResponce: %v", err)
	}
	// send packet
	err = r.AwaitChangesAndDeleteRequest(uint32(requestId))
	if err != nil {
		return Responce{}, fmt.Errorf("SendRequestWithResponce: %v", err)
	}
	// await changes
	resp, err := r.LoadResponceAndDelete(uint32(requestId))
	// load responce
	if err != nil {
		return Responce{}, fmt.Errorf("SendRequestWithResponce: %v", err)
	}
	return resp, nil
	// return
}
