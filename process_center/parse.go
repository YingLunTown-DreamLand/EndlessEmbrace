package ProcessCenter

import (
	RequestCenter "EndlessEmbrace/request_center"
	"encoding/json"
	"fmt"
)

type Header struct {
	PostType    string `json:"post_type"`
	ResponceId  string `json:"echo"`
	MessageType string `json:"message_type"`
}

// 将从 go-cqhttp 收到的 JSON 数据解析为 Golang 结构体
func UnMarshal(jsonDatas []byte) (interface{}, error) {
	var header Header
	err := json.Unmarshal(jsonDatas, &header)
	if err != nil {
		return nil, fmt.Errorf("UnMarshal: %v", err)
	}
	// read header
	if len(header.ResponceId) > 0 {
		return decodeToResponce(jsonDatas)
	}
	// if is responce
	switch header.PostType {
	case "message":
		return decodeToMessage(header, jsonDatas)
	case "message_sent":
		return decodeToMessage(header, jsonDatas)
	case "request":
	case "notice":
	case "meta_event":
		var new ClientStates
		err = json.Unmarshal(jsonDatas, &new)
		if err != nil {
			return nil, fmt.Errorf("UnMarshal: %v", err)
		}
		return new, nil
	}
	// decode json datas into golang struct
	return decodeToMap(jsonDatas)
	// default situation
}

// 将从 go-cqhttp 收到的 JSON 数据解析为 RequestCenter.Responce
func decodeToResponce(jsonDatas []byte) (interface{}, error) {
	var new RequestCenter.Responce
	err := json.Unmarshal(jsonDatas, &new)
	if err != nil {
		return nil, fmt.Errorf("UnMarshal: %v", err)
	}
	return new, nil
}

// 将从 go-cqhttp 收到的 JSON 数据解析为 map[string]interface{}
func decodeToMap(jsonDatas []byte) (interface{}, error) {
	var defaultStruct map[string]interface{}
	err := json.Unmarshal(jsonDatas, &defaultStruct)
	if err != nil {
		return nil, fmt.Errorf("UnMarshal: %v", err)
	}
	return defaultStruct, nil
}

// 将从 go-cqhttp 收到的 JSON 数据解析为 message 相关的结构体
func decodeToMessage(header Header, jsonDatas []byte) (interface{}, error) {
	switch header.MessageType {
	case "group":
		var new GroupMessage
		err := json.Unmarshal(jsonDatas, &new)
		if err != nil {
			return nil, fmt.Errorf("UnMarshal: %v", err)
		}
		new.Message = escapeCharacter(new.Message)
		new.RawMessage = escapeCharacter(new.RawMessage)
		return new, nil
	default:
		return decodeToMap(jsonDatas)
	}
	// process message
}
