package APIStruct

// 发送群聊消息(send_group_msg)
type SendGroupMsg struct {
	GroupId    int64  `json:"group_id"`    // 群号
	Message    string `json:"message"`     // 要发送的内容
	AutoEscape bool   `json:"auto_escape"` // 消息内容是否作为纯文本发送
}

// 发送私聊消息(send_private_msg)
type SendPrivateMsg struct {
	UserId     int64  `json:"user_id"`     // 对方 QQ 号
	GroupId    int64  `json:"group_id"`    // 主动发起临时会话时的来源群号（可选但机器人本身必须是管理员或群主）
	Message    string `json:"message"`     // 要发送的内容
	AutoEscape bool   `json:"auto_escape"` // 消息内容是否作为纯文本发送
}

// 撤回消息(delete_msg)
type DeleteMessage struct {
	MessageId int32 `json:"message_id"` // 消息 ID
}
