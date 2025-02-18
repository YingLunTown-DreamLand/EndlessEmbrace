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

func ProcessSearchUserByName(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
	defer func() {
		if r := recover(); r != nil {
			message = `用户名需要以 “"” 结尾`
		}
	}()

	reader.JumpSpace()
	if reader.Next(true) != `"` {
		return `用户名需要以“"”开头`
	}
	accountNameLower := strings.ToLower(reader.ParseString())

	resp, err := PostJSON[yorha_defines.ServerResponse](
		fmt.Sprintf("%s/qq_auth/search_user", YoRHaVerifyServerIP),
		yorha_qq_auth_key.QQAuthKey,
		yorha_defines.SearchUser{
			AdminGeneralFields: yorha_defines.AdminGeneralFields{
				AdminQQID:     sender.UserId,
				AdminUserName: sender.Card,
			},
			OperationType:    yorha_defines.OperationTypeSearchUserByName,
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
	case yorha_defines.ResponseTypeSearchOperationSuccess:
		return resp.UserData.String()
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}

func ProcessSearchUserByQQID(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
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
		fmt.Sprintf("%s/qq_auth/search_user", YoRHaVerifyServerIP),
		yorha_qq_auth_key.QQAuthKey,
		yorha_defines.SearchUser{
			AdminGeneralFields: yorha_defines.AdminGeneralFields{
				AdminQQID:     sender.UserId,
				AdminUserName: sender.Card,
			},
			OperationType: yorha_defines.OperationTypeSearchUserByQQID,
			QQID:          qqID,
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
	case yorha_defines.ResponseTypeSearchOperationSuccess:
		return resp.UserData.String()
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}
