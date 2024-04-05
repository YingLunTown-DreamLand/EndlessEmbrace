package main

import (
	APIStruct "EndlessEmbrace/api_struct"
	BotFunction "EndlessEmbrace/bot_function"
	RequestCenter "EndlessEmbrace/request_center"
	"fmt"

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

func process_run_codes_requests(commandLine string) (res string) {
	var language string = ""
	var content string = ""
	// init values
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

func (c *Client) MasterProcessingCenter(groupId int64, commandLine string) {
	var message string
	// prepare
	if message = process_uec_requests(commandLine); len(message) == 0 {
		message = process_run_codes_requests(commandLine)
		if len(message) == 0 {
			return
		}
	}
	// get message to send
	sendStruct := APIStruct.SendGroupMsg{
		GroupId:    groupId,
		Message:    message,
		AutoEscape: false,
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
		pterm.Warning.Printf("MasterProcessingCenter: %v\n", err)
		return
	}
	// send request with responce
	pterm.Success.Printf("%#v\n", resp)
	// output success information
}
