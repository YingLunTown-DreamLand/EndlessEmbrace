package ProcessCenter

import (
	"encoding/json"
	"fmt"
)

type Header struct {
	PostType string `json:"post_type"`
}

// 将从 go-cqhttp 收到的 JSON 数据解析为 Golang 结构体
func UnMarshal(jsonDatas []byte) (interface{}, error) {
	var header Header
	err := json.Unmarshal(jsonDatas, &header)
	if err != nil {
		return nil, fmt.Errorf("UnMarshal: %v", err)
	}
	// read header
	switch header.PostType {
	case "message":
	case "message_sent":
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
	var defaultStruct map[string]interface{}
	err = json.Unmarshal(jsonDatas, &defaultStruct)
	if err != nil {
		return nil, fmt.Errorf("UnMarshal: %v", err)
	}
	return defaultStruct, nil
	// default situation
}
