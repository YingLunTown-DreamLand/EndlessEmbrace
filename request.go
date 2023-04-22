package main

import (
	"EndlessEmbrace/APIStruct"
	"EndlessEmbrace/BotFunction"
	"EndlessEmbrace/RequestCenter"
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
	sendStruct := APIStruct.SendGroupMsg{
		GroupId: groupId,
		Message: fmt.Sprintf(
			"Output:\n%v\n-----\nError:\n%v",
			resp.StandardOutput,
			resp.StandardOutputError,
		),
		AutoEscape: true,
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
