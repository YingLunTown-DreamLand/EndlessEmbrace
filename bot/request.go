package bot

import (
	APIStruct "EndlessEmbrace/api_struct"
	BotFunction "EndlessEmbrace/bot_function"
	ProcessCenter "EndlessEmbrace/process_center"
	RequestCenter "EndlessEmbrace/request_center"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

func process_uec_requests(commandLine string) (res string) {
	if len(commandLine) < 5 {
		return
	}
	if len(commandLine) > 5 && commandLine[0:5] == "/uec " {
		res = commandLine[5:]
		res = BotFunction.UpgradeExecuteCommands(res)
	}
	return
}

func process_run_codes_requests(sender *ProcessCenter.GroupSender, commandLine string) (res string) {
	var language string = ""
	var content string = ""
	// init values
	switch strings.ToLower(sender.Role) {
	case "owner", "admin":
	default:
		if sender.UserId != 3527679800 && sender.UserId != 862713720 {
			return ""
		}
	}
	// pre check
	if len(commandLine) < 7 {
		return
	}
	switch commandLine[0:7] {
	case "code go":
		content = commandLine[7:]
		language = "go"
	case "code py":
		content = commandLine[7:]
		language = "python"
	default:
		return
	}
	// get content
	resp, err := BotFunction.RunGoAndPythonCodeByWebAPI(language, content)
	if err != nil {
		pterm.Warning.Printf("MasterProcessingCenter: %v\n", err)
		return
	}
	// get responce
	if len(resp.StandardOutputError) > 0 || len(resp.Error) > 0 {
		res = fmt.Sprintf("%v\n%v", resp.StandardOutputError, resp.Error)
	} else {
		res = resp.StandardOutput
	}
	// set result
	return
	// return
}

func (c *Client) SendMessageFastly(groupId int64, message string, allowCQ bool) {
	sendStruct := APIStruct.SendGroupMsg{
		GroupId:    groupId,
		Message:    message,
		AutoEscape: !allowCQ,
	}
	// construct the target struct
	resp, err := c.Resources.SendRequestWithResponce(
		c.Conn,
		RequestCenter.Request{
			Action:    APIStruct.SendGroupMsgAction,
			Params:    sendStruct,
			RequestId: fmt.Sprintf("%d", c.Resources.GetNewRequestId()),
		},
	)
	if err != nil {
		pterm.Warning.Printf("SendMessageFastly: %v\n", err)
		return
	}
	// send request with responce
	pterm.Success.Printf("%#v\n", resp)
	// output success information
}

func (c *Client) MasterProcessingCenter(
	groupId int64,
	messageId int32,
	sender ProcessCenter.GroupSender,
	originMessage []ProcessCenter.Message,
	commandLine string,
) {
	var message string
	// prepare
	if BotFunction.DeleteUnallowMsgIsEnabled {
		if BotFunction.MatchUnallowMessage(groupId, sender, originMessage, commandLine) {
			err := c.Resources.SendRequest(
				c.Conn,
				RequestCenter.Request{
					Action:    APIStruct.DeleteMsgAction,
					Params:    APIStruct.DeleteMessage{MessageId: messageId},
					RequestId: "",
				},
			)
			// delete target message
			if err != nil {
				pterm.Warning.Printf("MasterProcessingCenter: %v\n", err)
				return
			}
			// error check
			pterm.Success.Printf(
				"Match unallow message %#v on QQ group %d which sent from %#v; originMessage = %#v\n",
				commandLine, groupId, sender, originMessage,
			)
			// print success message
			return
			// return
		}
		// match and delete unallow message
	}
	// match unallow message
	if message = process_uec_requests(commandLine); len(message) != 0 {
		c.SendMessageFastly(groupId, message, false)
		return
	}
	if _, message = BotFunction.ProcessPlayMusic(commandLine); len(message) != 0 {
		c.SendMessageFastly(groupId, message, true)
		return
	}
	if message = process_run_codes_requests(&sender, commandLine); len(message) != 0 {
		c.SendMessageFastly(groupId, message, false)
		return
	}
	if _, message = BotFunction.ProcessYoRHaCommand(groupId, &sender, commandLine); len(message) != 0 {
		c.SendMessageFastly(groupId, message, false)
		return
	}
	// bot command running
}
