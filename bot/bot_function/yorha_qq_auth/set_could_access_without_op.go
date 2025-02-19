package yorha_qq_auth

import (
	yorha_defines "EndlessEmbrace/bot_function/yorha_qq_auth/defines"
	ProcessCenter "EndlessEmbrace/process_center"
	"fmt"
	"phoenixbuilder/fastbuilder/string_reader"
	"strconv"
	yorha_qq_auth_key "yorha/qq_auth_key"
)

func ProcessSetCouldAccessWithoutOP(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
	defer func() {
		if r := recover(); r != nil {
			message = `参数错误`
		}
	}()

	var rentalServerCode string = ""
	var setGeneral bool = false

	reader.JumpSpace()
	qqIDString, _ := reader.ParseNumber(true)
	qqID, err := strconv.ParseInt(qqIDString, 10, 64)
	if err != nil {
		return fmt.Sprintf("无效的 QQ 号：%v", err)
	}

	reader.JumpSpace()

	func() {
		defer func() {
			recover()
		}()
		rentalServerCode, _ = reader.ParseNumber(true)
		_, err := strconv.ParseInt(rentalServerCode, 10, 64)
		if err != nil {
			message = fmt.Sprintf("无效的租赁服号：%v", err)
		}
	}()
	if len(message) > 0 {
		return
	}
	if len(rentalServerCode) == 0 {
		setGeneral = true
	}

	resp, err := PostJSON[yorha_defines.ServerResponse](
		fmt.Sprintf("%s/qq_auth/set_could_access_without_op", YoRHaVerifyServerIP),
		yorha_qq_auth_key.QQAuthKey,
		yorha_defines.SetCouldAccessWithoutOP{
			AdminGeneralFields: yorha_defines.AdminGeneralFields{
				AdminQQID:     sender.UserId,
				AdminUserName: sender.SenderName(),
			},
			OperationType:      yorha_defines.OperationTypeSetServerCouldAccessWithoutOP,
			QQID:               qqID,
			SetGeneral:         setGeneral,
			RentalServerNumber: rentalServerCode,
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
	case yorha_defines.ResponseTypeSetCouldAccessWithoutOPSuccess:
		if setGeneral {
			return "操作成功，该用户现在可以无权限进入任何租赁服"
		}
		return fmt.Sprintf("操作成功，该用户现在可以无权限进入租赁服 %s", rentalServerCode)
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}

func ProcessUnsetCouldAccessWithoutOP(sender *ProcessCenter.GroupSender, reader *string_reader.StringReader) (message string) {
	defer func() {
		if r := recover(); r != nil {
			message = `参数错误`
		}
	}()

	var rentalServerCode string = ""
	var setGeneral bool = false

	reader.JumpSpace()
	qqIDString, _ := reader.ParseNumber(true)
	qqID, err := strconv.ParseInt(qqIDString, 10, 64)
	if err != nil {
		return fmt.Sprintf("无效的 QQ 号：%v", err)
	}

	reader.JumpSpace()

	func() {
		defer func() {
			recover()
		}()
		rentalServerCode, _ = reader.ParseNumber(true)
		_, err := strconv.ParseInt(rentalServerCode, 10, 64)
		if err != nil {
			message = fmt.Sprintf("无效的租赁服号：%v", err)
		}
	}()
	if len(message) > 0 {
		return
	}
	if len(rentalServerCode) == 0 {
		setGeneral = true
	}

	resp, err := PostJSON[yorha_defines.ServerResponse](
		fmt.Sprintf("%s/qq_auth/set_could_access_without_op", YoRHaVerifyServerIP),
		yorha_qq_auth_key.QQAuthKey,
		yorha_defines.SetCouldAccessWithoutOP{
			AdminGeneralFields: yorha_defines.AdminGeneralFields{
				AdminQQID:     sender.UserId,
				AdminUserName: sender.SenderName(),
			},
			OperationType:      yorha_defines.OperationTypeUnsetServerCouldAccessWithoutOP,
			QQID:               qqID,
			SetGeneral:         setGeneral,
			RentalServerNumber: rentalServerCode,
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
	case yorha_defines.ResponseTypeServerNotFound:
		return "该用户在目标租赁服本身就不能无权限进入"
	case yorha_defines.ResponseTypeSetCouldAccessWithoutOPSuccess:
		if setGeneral {
			return "操作成功，该用户现在失去了无权限进入任何租赁服的权能 (原本特别指定的租赁服仍然有效)"
		}
		return fmt.Sprintf("操作成功，该用户现在无法再无权限进入租赁服 %s", rentalServerCode)
	default:
		return "未知错误 (这看起来是个 BUG | 多次重试如果无效请联系群主处理)"
	}
}
