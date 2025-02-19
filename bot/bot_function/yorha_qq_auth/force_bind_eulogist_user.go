package yorha_qq_auth

import (
	yorha_defines "EndlessEmbrace/bot_function/yorha_qq_auth/defines"
	ProcessCenter "EndlessEmbrace/process_center"
	"fmt"
	"phoenixbuilder/fastbuilder/string_reader"
	"strconv"
	"strings"
	yorha_qq_auth_key "yorha/qq_auth_key"
)

func ProcessForceBindEulogistUser(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
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
	if reader.Next(true) != `"` {
		return `用户名需要以“"”开头`
	}
	accountNameLower := strings.ToLower(reader.ParseString())

	resp, err := PostJSON[yorha_defines.ServerResponse](
		fmt.Sprintf("%s/qq_auth/force_bind_eulogist_user", YoRHaVerifyServerIP),
		yorha_qq_auth_key.QQAuthKey,
		yorha_defines.ForceBindEulogistUser{
			AdminGeneralFields: yorha_defines.AdminGeneralFields{
				AdminQQID:     sender.UserId,
				AdminUserName: sender.SenderName(),
			},
			BindEulogistUser: yorha_defines.BindEulogistUser{
				OperationType:    yorha_defines.OperationTypeForceBindEulogistUser,
				QQID:             qqID,
				AccountNameLower: accountNameLower,
			},
		},
	)
	if err != nil {
		return fmt.Sprintf("操作失败：%v", err)
	}

	switch resp.ResponseType {
	case yorha_defines.ResponseTypeInvalidRequest:
		return "无效请求 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	case yorha_defines.ResponseTypeUserNotFound:
		return "没有找到指定的赞颂者账户"
	case yorha_defines.ResponseTypeUserHasBinded:
		return "该赞颂者账户已经绑定了 QQ 号，要先强制绑定需要将其强制解绑"
	case yorha_defines.ResponseTypeForceBindOperationSuccess:
		return fmt.Sprintf(
			"已成功将账户 %#v 和 QQ 号 %d 绑定",
			accountNameLower, qqID,
		)
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}

func ProcessForceUnbindEulogistUser(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
	defer func() {
		if r := recover(); r != nil {
			message = `无效的 QQ 号`
		}
	}()

	reader.JumpSpace()
	qqIDString, _ := reader.ParseNumber(true)
	qqID, err := strconv.ParseInt(qqIDString, 10, 64)
	if err != nil {
		return fmt.Sprintf("无效的 QQ 号：%v", err)
	}

	resp, err := PostJSON[yorha_defines.ServerResponse](
		fmt.Sprintf("%s/qq_auth/force_bind_eulogist_user", YoRHaVerifyServerIP),
		yorha_qq_auth_key.QQAuthKey,
		yorha_defines.ForceBindEulogistUser{
			AdminGeneralFields: yorha_defines.AdminGeneralFields{
				AdminQQID:     sender.UserId,
				AdminUserName: sender.SenderName(),
			},
			BindEulogistUser: yorha_defines.BindEulogistUser{
				OperationType: yorha_defines.OperationTypeForceUnbindEulogistUser,
				QQID:          qqID,
			},
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
	case yorha_defines.ResponseTypeForceBindOperationSuccess:
		return fmt.Sprintf(
			"已成功将账户 %#v 与 QQ 号 %d 解绑",
			resp.OriginAccountNameLower, qqID,
		)
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}
