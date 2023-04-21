package main

import (
	"EndlessEmbrace/ProcessCenter"

	"github.com/gorilla/websocket"
)

// 描述与 go-cqhttp 建立的客户端的结构体
type Client struct {
	Conn         *websocket.Conn
	ConnAns      *ConnectionReulst
	ClientStates *ProcessCenter.ClientStates
}

// 此结构体用于描述客户端与 go-cqhttp 的连接结果。
// 通常情况下，此信息是在与 go-cqhttp 建立连接后由 go-cqhttp 发送的第一个信息。
// 对应 go-cqhttp 中的 生命周期 数据类型
type ConnectionReulst struct {
	PostMethod    uint8   `json:"_post_method"`
	MetaEventType string  `json:"meta_event_type"`
	PostType      string  `json:"post_type"`
	SelfId        float64 `json:"self_id"`
	SubType       string  `json:"sub_type"`
	Time          int64   `json:"time"`
}
