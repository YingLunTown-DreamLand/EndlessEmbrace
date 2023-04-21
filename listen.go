package main

import (
	"EndlessEmbrace/ProcessCenter"
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
		case ProcessCenter.ClientStates:
			c.ClientStates = &new
		case map[string]interface{}:
			pterm.Info.Printf("%#v\n", new)
		}
		// do some actions
	}
}
