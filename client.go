package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// 创建一个客户端并让它连接到 address 所指代的 go-cqhttp 服务器。
// address 是一个 url 地址，例如 ws://127.0.0.1:8080 。
// 注：go-cqhttp 服务器是一个 websocket 服务器
func NewClient(address string) (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		return &Client{}, fmt.Errorf("NewClient: %v", err)
	}
	// 建立连接
	newStruct := Client{}
	newStruct.Conn = conn
	// 初始化结构体
	err = newStruct.init()
	if err != nil {
		conn.Close()
		return &Client{}, fmt.Errorf("NewClient: %v", err)
	}
	// 初始化客户端
	return &newStruct, nil
	// 返回值
}

// 用于在与 go-cqhttp 建立连接后读取连接状态信息。
// 通常情况下，这个信息是建立连接后由 go-cqhttp 所发送的首个信息
func (c *Client) init() error {
	var new ConnectionReulst
	// init values
	err := c.Conn.ReadJSON(&new)
	if err != nil {
		return fmt.Errorf("init: %v", err)
	}
	// read the first messages
	c.ConnAns = &new
	// set values
	return nil
	// return
}
