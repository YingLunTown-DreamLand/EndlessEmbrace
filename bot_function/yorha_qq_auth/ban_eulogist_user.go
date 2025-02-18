package yorha_qq_auth

import (
	yorha_defines "EndlessEmbrace/bot_function/yorha_qq_auth/defines"
	yorha_qq_auth_key "EndlessEmbrace/depends/YoRHa/qq_auth_key"
	ProcessCenter "EndlessEmbrace/process_center"
	"fmt"
	"phoenixbuilder/fastbuilder/string_reader"
	"strconv"
	"time"
)

func ProcessBanEulogistUser(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
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
	durationOfSecondsString, _ := reader.ParseNumber(true)
	durationOfSeconds, err := strconv.ParseInt(durationOfSecondsString, 10, 32)
	if err != nil {
		return fmt.Sprintf("无效的封禁时长：%v", err)
	}

	resp, err := PostJSON[yorha_defines.ServerResponse](
		fmt.Sprintf("%s/qq_auth/ban_eulogist_user", YoRHaVerifyServerIP),
		&yorha_qq_auth_key.QQAuthKey.PublicKey,
		yorha_defines.BanEulogistUser{
			AdminGeneralFields: yorha_defines.AdminGeneralFields{
				AdminQQID:     sender.UserId,
				AdminUserName: sender.Card,
			},
			OperationType: yorha_defines.OperationTypeBanEulogistUser,
			QQID:          qqID,
			Duration:      int32(durationOfSeconds),
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
	case yorha_defines.ResponseTypeUserHasBanned:
		return fmt.Sprintf(
			"操作失败，因为对应的赞颂者正处于封禁状态 (将在 %s 时解封)",
			time.Unix(resp.UnbanUnixTime, 0).Format("2006-01-02 15:04:05"),
		)
	case yorha_defines.ResponseTypeBanOperationSuccess:
		return fmt.Sprintf(
			"封禁成功，解封时间为 %s",
			time.Unix(resp.UnbanUnixTime, 0).Format("2006-01-02 15:04:05"),
		)
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}

func ProcessUnbanEulogistUser(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
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
		fmt.Sprintf("%s/qq_auth/ban_eulogist_user", YoRHaVerifyServerIP),
		&yorha_qq_auth_key.QQAuthKey.PublicKey,
		yorha_defines.BanEulogistUser{
			AdminGeneralFields: yorha_defines.AdminGeneralFields{
				AdminQQID:     sender.UserId,
				AdminUserName: sender.Card,
			},
			OperationType: yorha_defines.OperationTypeUnbanEulogistUser,
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
	case yorha_defines.ResponseTypeUserNotBanned:
		return "该用户未处于封禁状态"
	case yorha_defines.ResponseTypeBanOperationSuccess:
		return "已成功解封"
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}
