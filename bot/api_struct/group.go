package APIStruct

// 获取群成员信息(get_group_member_info)
type GetGroupMemberInfo struct {
	GroupId int64 `json:"group_id"` // 群号
	UserID  int64 `json:"user_id"`  // QQ 号
	NoCache bool  `json:"no_cache"` // 是否不使用缓存（使用缓存可能更新不及时，但响应更快）
}
