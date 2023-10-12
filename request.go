package main

import (
	APIStruct "EndlessEmbrace/api_struct"
	BotFunction "EndlessEmbrace/bot_function"
	RequestCenter "EndlessEmbrace/request_center"
	"fmt"

	"github.com/pterm/pterm"
)

func (c *Client) MasterProcessingCenter(groupId int64, commandLine string) {
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
	var message string
	if len(resp.StandardOutputError) > 0 || len(resp.Error) > 0 {
		message = fmt.Sprintf("%v\n%v", resp.StandardOutputError, resp.Error)
	} else {
		message = resp.StandardOutput
	}
	sendStruct := APIStruct.SendGroupMsg{
		GroupId:    groupId,
		Message:    message,
		AutoEscape: false,
	}
	// construct the target struct
	goCqhttpResp, err := c.Resources.SendRequestWithResponce(
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
	pterm.Success.Printf("%#v\n", goCqhttpResp)
	// output success information
}
