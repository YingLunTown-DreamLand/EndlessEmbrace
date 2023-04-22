package APIStruct

// 发送群聊消息(send_group_msg)
type SendGroupMsg struct {
	GroupId    int64  `json:"group_id"`    // 群号
	Message    string `json:"message"`     // 要发送的内容
	AutoEscape bool   `json:"auto_escape"` // 消息内容是否作为纯文本发送
}
