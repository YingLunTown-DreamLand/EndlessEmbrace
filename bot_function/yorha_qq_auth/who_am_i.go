package yorha_qq_auth

import (
	yorha_defines "EndlessEmbrace/bot_function/yorha_qq_auth/defines"
	yorha_qq_auth_key "EndlessEmbrace/depends/YoRHa/qq_auth_key"
	ProcessCenter "EndlessEmbrace/process_center"
	"fmt"
	"phoenixbuilder/fastbuilder/string_reader"
)

func ProcessWhoAmI(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
	resp, err := PostJSON[yorha_defines.ServerResponse](
		fmt.Sprintf("%s/qq_auth/who_am_i", YoRHaVerifyServerIP),
		&yorha_qq_auth_key.QQAuthKey.PublicKey,
		yorha_defines.WhoAmI{
			QQID: sender.UserId,
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
	case yorha_defines.ResponseTypeWhoAmISuccess:
		return resp.UserData.String()
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}
