package ProcessCenter

// 此结构体用于描述 go-cqhttp 的运行情况。
// 对应 go-cqhttp 中的 心跳包 数据类型
type ClientStates struct {
	Interval      float64             `json:"interval"`        // 距离上一次心跳包的时间(单位是毫秒)
	MetaEventType string              `json:"meta_event_type"` // 元事件类型
	PostType      string              `json:"post_type"`       // 上报类型
	SelfId        float64             `json:"self_id"`         // 收到事件的机器人 QQ 号
	Time          int64               `json:"time"`            // 事件发生的时间戳
	Status        ClientStates_Status `json:"status"`          // 应用程序状态
}

// 应用程序状态
type ClientStates_Status struct {
	AppEnabled     bool                     `json:"app_enabled"`     // 程序是否可用
	AppGood        bool                     `json:"app_good"`        // 程序正常
	AppInitialized bool                     `json:"app_initialized"` // 程序是否初始化完毕
	Good           bool                     `json:"good"`            // ...
	Online         bool                     `json:"online"`          // 是否在线
	PluginsGood    interface{}              `json:"plugins_good"`    // 插件正常(可能为 null)
	Stat           ClientStates_Status_Stat `json:"stat"`            // 统计信息
}

// 统计信息
type ClientStates_Status_Stat struct {
	DisconnectTimes int64 `json:"disconnect_times"`  // 连接断开次数
	LastMessageTime int64 `json:"last_message_time"` // 最后一次消息时间
	LostTimes       int64 `json:"lost_times"`        // 连接丢失次数
	MessageReceived int64 `json:"message_received"`  // 消息接收数
	MessageSent     int64 `json:"message_sent"`      // 消息发送数
	PacketLost      int64 `json:"packet_lost"`       // 丢包数
	PacketReceived  int64 `json:"packet_received"`   // 收包数
	PacketSent      int64 `json:"packet_sent"`       // 发包数
}
