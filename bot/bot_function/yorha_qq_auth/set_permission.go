package yorha_qq_auth

import (
	yorha_defines "EndlessEmbrace/bot_function/yorha_qq_auth/defines"
	ProcessCenter "EndlessEmbrace/process_center"
	"fmt"
	"phoenixbuilder/fastbuilder/string_reader"
	"strconv"
	yorha_qq_auth_key "yorha/qq_auth_key"
)

func ProcessSetPermission(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
	defer func() {
		if r := recover(); r != nil {
			message = `参数错误`
		}
	}()

	reader.JumpSpace()
	qqIDString, _ := reader.ParseNumber(true)
	qqID, err := strconv.ParseInt(qqIDString, 10, 64)
	if err != nil {
		return fmt.Sprintf("无效的 QQ 号：%v", err)
	}

	reader.JumpSpace()
	permissionLevelString, _ := reader.ParseNumber(true)
	permissionLevel, err := strconv.ParseInt(permissionLevelString, 10, 8)
	if err != nil {
		return fmt.Sprintf("无效的权限等级：%v", err)
	}

	resp, err := PostJSON[yorha_defines.ServerResponse](
		fmt.Sprintf("%s/qq_auth/set_permission", YoRHaVerifyServerIP),
		&yorha_qq_auth_key.QQAuthKey.PublicKey,
		yorha_defines.SetPermission{
			AdminGeneralFields: yorha_defines.AdminGeneralFields{
				AdminQQID:     sender.UserId,
				AdminUserName: sender.Card,
			},
			QQID:       qqID,
			Permission: uint8(permissionLevel),
		},
	)
	if err != nil {
		return fmt.Sprintf("操作失败：%v", err)
	}

	switch resp.ResponseType {
	case yorha_defines.ResponseTypeInvalidRequest:
		return "无效请求 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	case yorha_defines.ResponseTypeUserNotBinded:
		return "该 QQ 号没有绑定任何赞颂者账户"
	case yorha_defines.ResponseTypeUnknownPermissionLevel:
		return "无效的权限等级 (仅接受 0~4 之间的值，含边界)"
	case yorha_defines.ResponseTypeSetPermissionSuccess:
		return fmt.Sprintf("已成功将该用户的权限等级设置为 %s", resp.ResultPermissionName)
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}
