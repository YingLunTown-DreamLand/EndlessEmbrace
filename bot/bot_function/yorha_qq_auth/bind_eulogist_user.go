package yorha_qq_auth

import (
	yorha_defines "EndlessEmbrace/bot_function/yorha_qq_auth/defines"
	ProcessCenter "EndlessEmbrace/process_center"
	"fmt"
	"phoenixbuilder/fastbuilder/string_reader"
	"strings"
	yorha_qq_auth_key "yorha/qq_auth_key"
)

func ProcessBindEulogistUser(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
	defer func() {
		if r := recover(); r != nil {
			message = `用户名没有以“"”结尾`
		}
	}()

	reader.JumpSpace()
	if reader.Next(true) != `"` {
		return `用户名需要以“"”开头`
	}
	accountNameLower := strings.ToLower(reader.ParseString())

	resp, err := PostJSON[yorha_defines.ServerResponse](
		fmt.Sprintf("%s/qq_auth/bind_eulogist_user", YoRHaVerifyServerIP),
		yorha_qq_auth_key.QQAuthKey,
		yorha_defines.BindEulogistUser{
			OperationType:    yorha_defines.OperationTypeBindEulogistUser,
			QQID:             sender.UserId,
			AccountNameLower: accountNameLower,
		},
	)
	if err != nil {
		return fmt.Sprintf("操作失败：%v", err)
	}

	switch resp.ResponseType {
	case yorha_defines.ResponseTypeInvalidRequest:
		return "无效请求 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	case yorha_defines.ResponseTypeUserNotFound:
		return "指定的赞颂者用户没有找到"
	case yorha_defines.ResponseTypeUserHasBinded:
		return "指定的赞颂者用户已经绑定了 QQ 号"
	case yorha_defines.ResponseTypeUserHasBanned:
		return "目标赞颂者用户已被封禁，任何绑定(解绑)操作都是无效的"
	case yorha_defines.ResponseTypeBindOperationSuccess:
		return fmt.Sprintf("已成功将账号 %#v 与您的 QQ 号进行绑定", accountNameLower)
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}

func ProcessUnbindEulogistUser(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
	resp, err := PostJSON[yorha_defines.ServerResponse](
		fmt.Sprintf("%s/qq_auth/bind_eulogist_user", YoRHaVerifyServerIP),
		yorha_qq_auth_key.QQAuthKey,
		yorha_defines.BindEulogistUser{
			OperationType: yorha_defines.OperationTypeUnbindEulogistUser,
			QQID:          sender.UserId,
		},
	)
	if err != nil {
		return fmt.Sprintf("操作失败：%v", err)
	}

	switch resp.ResponseType {
	case yorha_defines.ResponseTypeInvalidRequest:
		return "无效请求 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	case yorha_defines.ResponseTypeUserNotBinded:
		return "您的 QQ 号没有绑定任何赞颂者账户"
	case yorha_defines.ResponseTypeUserHasBanned:
		return "目标赞颂者用户已被封禁，任何绑定(解绑)操作都是无效的"
	case yorha_defines.ResponseTypeBindOperationSuccess:
		return fmt.Sprintf("您的 QQ 号已与原赞颂者账户 %#v 解绑", resp.OriginAccountNameLower)
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}
