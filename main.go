package main

import (
	"sync"

	"github.com/pterm/pterm"
)

func main() {
	pterm.Info.Println("正在连接ing...")
	client, err := NewClient("ws://127.0.0.1:3001")
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	pterm.Success.Printf("%#v\n", client.ConnAns)
	// connect to the onebot server
	closeDown := sync.Mutex{}
	go func() {
		for {
			client.ReadPacketAndProcess(&closeDown)
			pterm.Warning.Println("连接丢失，正在重新连接...")
			client, err = NewClient("ws://127.0.0.1:3001")
			if err != nil {
				pterm.Error.Println(err)
				return
			}
			pterm.Success.Printf("重连成功: %#v\n", client.ConnAns)
		}
	}()
	// process messages
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	waitGroup.Wait()
	// set wait groups
}
