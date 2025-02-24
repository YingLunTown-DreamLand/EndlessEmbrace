package BotFunction

import (
	"EndlessEmbrace/bot_function/yorha_qq_auth"
	ProcessCenter "EndlessEmbrace/process_center"
	"phoenixbuilder/fastbuilder/string_reader"
	"strings"
)

const (
	EulogistGroupID = 644154294
)

var ConstSpecificAdmin []int64 = []int64{
	862713720,
}

const HelpCommandConstText = `可用命令如下。
	- 账户绑定 -
		/bind <userName: string> | 将 userName 指代的赞颂者用户绑定到当前 QQ 号
		/unbind | 解绑任何绑定到当前 QQ 号的赞颂者账户

	- 账户信息 -
		/whoami | 查询个人赞颂者账户信息

	- 管理员命令(仅限群管理员可用) -
		- 账户封禁 -
			/ban <qqID: int64> <durationOfSeconds: int32> | 将 qqID 绑定的赞颂者用户封禁 durationOfSeconds 秒
			/unban <qqID: int64> | 解除 qqID 绑定的赞颂者用户的封禁状态

		- 账户绑定 -
			/force_bind <qqID: int64> <userName: string> | 将名为 userName 的赞颂者账户强制绑定到 qqID 上
			/force_unbind <qqID: int64> | 将任何绑定到 qqID 的赞颂者账户从该 QQ 号上解绑

		- 无权限进服 -
			/op <qqID: int64> [rentalServerNumber: string] | 允许 qqID 绑定的赞颂者账户无权限进入 rentalServerNumber 所指代的租赁服。不填租赁服号则默认该用户可无权限进入任意租赁服
			/deop <qqID: int64> [rentalServerNumber: string] | 使 qqID 绑定的赞颂者账户不再能无权限进入 rentalServerNumber 所指代的租赁服。不填租赁服号则默认该用户无法无权限进入除特别指定的任何租赁服

		- 管理用户权限 -
			/set_permission <qqID: int64> <permissionLevel: uint8> | 将 qqID 绑定的赞颂者账户的权限等级设为 permissionLevel (0-受限用户, 1-普通用户, 2-受信用户, 3-管理员, 4-系统)

		- 查询用户信息 -
			/search_user_by_name <userName: string> | 查询名为 userName 的赞颂者用户的信息
			/search_user_by_qqid <qqID: int64> | 查询 qqID 所绑定的赞颂者用户的信息`

func ProcessYoRHaCommand(groupID int64, sender *ProcessCenter.GroupSender, commandLine string) (
	shouldSendMessage bool,
	message string,
) {
	var commandName string = ""
	var requesterIsAdmin bool = false
	reader := string_reader.NewStringReader(&commandLine)

	if groupID != EulogistGroupID {
		return false, ""
	}

	reader.JumpSpace()
	if reader.Next(true) != "/" {
		return false, ""
	}

	for {
		token := reader.Next(true)
		if token == "" || token == " " || token == "\n" || token == "\r" {
			break
		}
		commandName += token
	}
	commandName = strings.ToLower(commandName)

	switch commandName {
	case "bind":
		return true, yorha_qq_auth.ProcessBindEulogistUser(sender, reader)
	case "unbind":
		return true, yorha_qq_auth.ProcessUnbindEulogistUser(sender, reader)
	case "whoami":
		return true, yorha_qq_auth.ProcessWhoAmI(sender, reader)
	case "help":
		return true, HelpCommandConstText
	}

	switch sender.Role {
	case "owner", "admin":
	default:
		for _, value := range ConstSpecificAdmin {
			if value == sender.UserId {
				requesterIsAdmin = true
			}
		}
		if !requesterIsAdmin {
			return true, "未知的命令或您权限不足"
		}
	}

	switch commandName {
	case "ban":
		return true, yorha_qq_auth.ProcessBanEulogistUser(sender, reader)
	case "unban":
		return true, yorha_qq_auth.ProcessUnbanEulogistUser(sender, reader)
	case "force_bind":
		return true, yorha_qq_auth.ProcessForceBindEulogistUser(sender, reader)
	case "force_unbind":
		return true, yorha_qq_auth.ProcessForceUnbindEulogistUser(sender, reader)
	case "op":
		return true, yorha_qq_auth.ProcessSetCouldAccessWithoutOP(sender, reader)
	case "deop":
		return true, yorha_qq_auth.ProcessUnsetCouldAccessWithoutOP(sender, reader)
	case "set_permission":
		return true, yorha_qq_auth.ProcessSetPermission(sender, reader)
	case "search_user_by_name":
		return true, yorha_qq_auth.ProcessSearchUserByName(sender, reader)
	case "search_user_by_qqid":
		return true, yorha_qq_auth.ProcessSearchUserByQQID(sender, reader)
	}

	return true, "未知的命令"
}
