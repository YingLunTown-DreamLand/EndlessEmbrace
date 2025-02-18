package ProcessCenter

// 下表列出了 QQ 群中群成员可以持有的权限等级，
// 除访客以外，必须为下表其中之一
const (
	GroupRoleMember = "member" // 普通群成员
	GroupRoleAdmin  = "admin"  // 群管理员
	GroupRoleOwner  = "owner"  // 群主
)
