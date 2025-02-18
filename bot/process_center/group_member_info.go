package ProcessCenter

// 群成员信息
type GroupMemberInfo struct {
	GroupID         int64   `json:"group_id"`          // 群号
	UserID          int64   `json:"user_id"`           // QQ 号
	NickName        string  `json:"nickname"`          // 昵称
	Card            string  `json:"card"`              // 群名片/备注
	Sex             string  `json:"sex"`               // 性别 (male/female/unknown)
	Age             int32   `json:"age"`               // 年龄
	Area            string  `json:"area"`              // 地区
	JoinTime        int32   `json:"join_time"`         // 加群时间戳
	LastSentTime    int32   `json:"last_sent_time"`    // 最后发言时间戳
	Level           string  `json:"level"`             // 成员等级
	QQLevel         float64 `json:"qq_level"`          // QQ 等级
	Role            string  `json:"role"`              // 角色 (owner/admin/member)
	Unfriendly      bool    `json:"unfriendly"`        // 是否不良记录成员
	Title           string  `json:"title"`             // 专属头衔
	TitleExpireTime int64   `json:"title_expire_time"` // 专属头衔过期时间戳
	CardChangeable  bool    `json:"card_changeable"`   // 是否允许修改群名片
	ShutUpTimestamp int64   `json:"ShutUpTimestamp"`   // 禁言到期时间
}
