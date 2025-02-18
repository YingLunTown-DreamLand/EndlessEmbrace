package main

import (
	bot "EndlessEmbrace"
	"fmt"
	"os/exec"
	"sync"
	"time"
	"yorha_http_service"

	"github.com/pterm/pterm"
)

const (
	BotQQID = "862713720"
)

func main() {
	client, err := bot.NewClient("ws://127.0.0.1:8080")
	if err != nil {
		pterm.Error.Println(err)
		// print error info
		pterm.Info.Println("尝试自动重启 NapCat")
		err := exec.Command("napcat", BotQQID, "restart").Run()
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("重启出现问题: %v", err))
		}
		pterm.Info.Println("10 秒后终止程序")
		time.Sleep(time.Second * 10)
		// try restart napcat
		return
		// return
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
