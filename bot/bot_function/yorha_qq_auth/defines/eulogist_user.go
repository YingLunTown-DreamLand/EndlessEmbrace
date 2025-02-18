package yorha_defines

import (
	"fmt"
	"strings"
	"time"
)

// 描述赞颂者用户权限等级的常量
const (
	UserPermissionDefault uint8 = iota
	UserPermissionDeveloper
	UserPermissionAgent
	UserPermissionAdmin
	UserPermissionSystem
)

// 描述单个赞颂者用户的信息
type User struct {
	AccountNameLower      string `json:"account_name_lower"`
	AccountPasswordSum256 string `json:"account_password_sum_256"`
	BindedQQID            string `json:"binded_qq_id"`

	UnbanUnixTime        int64 `json:"unban_unix_time"`
	UserPermission       uint8 `json:"user_permission"`
	CouldAccessWithoutOP bool  `json:"could_access_without_op"`

	UserEulogistToken   string `json:"user_eulogist_token"`
	UserAuthServerToken string `json:"user_auth_server_token"`
	UserProvidedCookie  string `json:"user_provided_cookie"`

	ServerCouldAccess []ServerCouldAccess `json:"server_could_access"`
}

// 描述单个已绑定的卡槽信息
type ServerCouldAccess struct {
	ServerCode           string `json:"server_could_access"`
	CouldAccessWithoutOP bool   `json:"could_access_without_op"`
	ExpireUnixTime       int64  `json:"expire_unix_time"`
}

// 将 serverCode 转换为安全形式。
// 例如，1234567 将被置为 1*****7。
// 特别地，当 serverCode 的长度为 2 时，
// 将直接返回原值
func (u User) ToSafeRentalServerCode(serverCode string) string {
	if len(serverCode) == 2 {
		return serverCode
	}
	return serverCode[0:1] + strings.Repeat("*", len(serverCode)-2) + serverCode[len(serverCode)-1:]
}

// 返回 boolean 的字符串表达。
// 如果传入真，则返回“是”，
// 否则返回“否”
func (u User) BoolString(boolean bool) string {
	if boolean {
		return "是"
	}
	return "否"
}

// 返回 u.UserPermission 的字符串表达
func (u User) PermissionString() string {
	switch u.UserPermission {
	case UserPermissionDefault:
		return "受限用户"
	case UserPermissionDeveloper:
		return "普通用户"
	case UserPermissionAgent:
		return "受信任用户"
	case UserPermissionAdmin:
		return "系统管理员"
	case UserPermissionSystem:
		return "系统"
	default:
		return "(未知)"
	}
}

// 返回 u 即 User 的字符串表达
func (u User) String() string {
	resultString := "" +
		"用户信息如下: \n" +
		"	用户名: " + u.AccountNameLower + "\n"

	if len(u.BindedQQID) > 0 {
		resultString += "	绑定的 QQ 号: " + u.BindedQQID + "\n"
	} else {
		resultString += "	绑定的 QQ 号: (未绑定)\n"
	}

	if u.UnbanUnixTime > time.Now().Unix() {
		resultString += "	解封时间: " + time.Unix(u.UnbanUnixTime, 0).Format("2006-01-02 15:04:05") + "\n"
	} else {
		resultString += "	解封时间: (未封禁)\n"
	}

	resultString += fmt.Sprintf(
		"	用户权限等级: %s\n",
		u.PermissionString(),
	)

	resultString += fmt.Sprintf(
		"	是否可以无 OP 访问任何租赁服: %s\n",
		u.BoolString(u.CouldAccessWithoutOP || u.UserPermission >= UserPermissionAdmin),
	)

	resultString += fmt.Sprintf(
		"	是否创建了辅助用户: %s\n",
		u.BoolString(len(u.UserAuthServerToken) > 0),
	)

	resultString += fmt.Sprintf(
		"	是否暂存了 Minecraft 账户电子存根: %s\n",
		u.BoolString(len(u.UserProvidedCookie) > 0),
	)

	if len(u.ServerCouldAccess) > 0 {
		resultString += "	用户已绑定的租赁服信息如下: \n"
		for _, value := range u.ServerCouldAccess {
			resultString += fmt.Sprintf(
				"		租赁服号: %s | 卡槽过期时间: %s | 是否可以无 OP 访问: %s\n",
				u.ToSafeRentalServerCode(value.ServerCode),
				time.Unix(value.ExpireUnixTime, 0).Format("2006-01-02 15:04:05"),
				u.BoolString(value.CouldAccessWithoutOP),
			)
		}
	} else {
		resultString += "	用户已绑定的租赁服信息如下: (没有绑定任何租赁服)\n"
	}

	return resultString
}
