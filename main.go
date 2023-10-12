package main

import (
	BotFunction "EndlessEmbrace/bot_function"
	"sync"

	"github.com/pterm/pterm"
)

func main() {
	client, err := NewClient("ws://127.0.0.1:6700")
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	pterm.Success.Printf("%#v\n", client.ConnAns)
	// connect to the go-cqhttp server
	go BotFunction.RepeatTiming(client.Conn, client.Resources)
	// do the routines
	// ^ need reconstruct
	closeDown := sync.Mutex{}
	go client.ReadPacketAndProcess(&closeDown)
	// process messages
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	waitGroup.Wait()
	// set wait groups
}
