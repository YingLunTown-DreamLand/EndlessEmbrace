package main

import (
	ProcessCenter "EndlessEmbrace/process_center"
	RequestCenter "EndlessEmbrace/request_center"
	"strconv"
	"sync"

	"github.com/pterm/pterm"
)

func (c *Client) ReadPacketAndProcess(
	closeDown *sync.Mutex,
) {
	for {
		if !closeDown.TryLock() {
			return
		} else {
			closeDown.Unlock()
		}
		// if need to close this func
		_, bytes, err := c.Conn.ReadMessage()
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		// read json datas
		ans, err := ProcessCenter.UnMarshal(bytes)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		// unmarshal json datas to the golang struct
		switch new := ans.(type) {
		case RequestCenter.Responce:
			requestID, err := strconv.ParseUint(new.ResponceId, 10, 32)
			if err != nil {
				pterm.Warning.Printf("ReadPacketAndProcess: %v\n", err)
				break
			}
			c.Resources.WriteResponce(uint32(requestID), new)
		case ProcessCenter.ClientStates:
			c.ClientStates = &new
		case ProcessCenter.GroupMessage:
			go c.MasterProcessingCenter(new.GroupId, new.Message)
		case map[string]interface{}:
			pterm.Info.Printf("%#v\n", new)
		}
		// do some actions
	}
}
