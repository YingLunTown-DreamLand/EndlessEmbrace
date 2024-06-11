package BotFunction

import (
	ProcessCenter "EndlessEmbrace/process_center"
	"reflect"
	"strings"
)

// 描述该功能(撤回群聊中的不允许消息)是否需要开启
const DeleteUnallowMsgIsEnabled = true

// 下表描述了需要启用该功能的群聊
var GroupNumber []int64 = []int64{
	194838530,
}

// 下表描述了不被允许的消息
var UnallowMessageList []string = []string{}

// 下表描述了不被允许的 CQ 码
var UnallowCQCodeList []ProcessCenter.Message = []ProcessCenter.Message{
	{
		Data: map[string]any{"qq": "2528622340"},
		Type: "at",
	},
}

// 下表用于放置多名群成员的 QQ 号。
// 任何被列于下表的成员的言论将不受本功能限制
var WhiteList []int64 = []int64{}

func MatchUnallowMessage(
	groupId int64,
	sender ProcessCenter.GroupSender,
	originMessage []ProcessCenter.Message,
	rawMessage string,
) (matched bool) {
	var match bool = false
	// prepare
	switch sender.Role {
	case ProcessCenter.GroupRoleOwner, ProcessCenter.GroupRoleAdmin:
		return false
	}
	// check the role of the target sender
	for _, value := range GroupNumber {
		if groupId == value {
			match = true
			break
		}
	}
	if !match {
		return false
	}
	// check group number
	for _, value := range WhiteList {
		if sender.UserId == value {
			return false
		}
	}
	// check white list
	for _, value := range UnallowMessageList {
		if strings.Contains(rawMessage, value) {
			return true
		}
	}
	// match each message
	for _, value := range originMessage {
		for _, v := range UnallowCQCodeList {
			if value.Type != v.Type {
				continue
			}
			if reflect.DeepEqual(value.Data, v.Data) {
				return true
			}
		}
	}
	// match CQ code
	return false
	// return
}
