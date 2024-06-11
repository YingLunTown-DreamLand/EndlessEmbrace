package ProcessCenter

// 描述消息内容中的单个元素
type Message struct {
	Data map[string]any `json:"data"`
	Type string         `json:"type"`
}

// 群消息
type GroupMessage struct {
	Anonymous   *Anonymous  `json:"anonymous"`    // 匿名信息, 如果不是匿名消息则为 null
	Font        int         `json:"font"`         // 字体
	GroupId     int64       `json:"group_id"`     // 群号
	Message     []Message   `json:"message"`      // 消息内容
	MessageId   int32       `json:"message_id"`   // 消息 ID
	MessageSeq  int64       `json:"message_seq"`  // 起始消息序号，可通过 get_msg 获得
	MessageType string      `json:"message_type"` // 消息类型
	PostType    string      `json:"post_type"`    // 表示该上报的类型
	RawMessage  string      `json:"raw_message"`  // CQ 码格式的消息
	RealId      int64       `json:"real_id"`      // ... [未知字段，看起来与 MessageSeq 是一致的]
	SelfId      int64       `json:"self_id"`      // 收到事件的机器人 QQ 号
	SubType     string      `json:"sub_type"`     // 消息子类型(normal/anonymous/notice)
	Time        int64       `json:"time"`         // 事件发生的时间戳
	UserId      int64       `json:"user_id"`      // 发送者 QQ 号
	Sender      GroupSender `json:"sender"`       // 发送人信息
}

// 匿名信息
type Anonymous struct {
	Id   int64  `json:"id"`   // 匿名用户 ID
	Name string `json:"name"` // 匿名用户名称
	Flag string `json:"flag"` // 匿名用户 flag ，在调用禁言 API 时需要传入
}

// 发送人信息
type GroupSender struct {
	Age      int32  `json:"age"`      // 年龄
	Area     string `json:"area"`     // 地区
	Card     string `json:"card"`     // 群名片/备注
	Level    string `json:"level"`    // 成员等级
	Nickname string `json:"nickname"` // 昵称
	Role     string `json:"role"`     // 角色(owner/admin/member)
	Sex      string `json:"sex"`      // 性别(male/female/unknown)
	Title    string `json:"title"`    // 专属头衔
	UserId   int64  `json:"user_id"`  // 发送者 QQ 号
}
