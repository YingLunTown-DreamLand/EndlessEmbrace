package main

import (
	bot "EndlessEmbrace"
	"sync"
	"yorha_http_service"

	"github.com/pterm/pterm"
)

func main() {
	client, err := bot.NewClient("ws://127.0.0.1:8080")
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	pterm.Success.Printf("%#v\n", client.ConnAns)
	// connect to the go-cqhttp server
	closeDown := sync.Mutex{}
	go client.ReadPacketAndProcess(&closeDown)
	go yorha_http_service.NewRouter(client).GinEngine.Run(":2018")
	// process messages
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	waitGroup.Wait()
	// set wait groups
}
