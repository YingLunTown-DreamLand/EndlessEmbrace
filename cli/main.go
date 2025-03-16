package main

import (
	bot "EndlessEmbrace"
	"os/exec"
	"sync"
	"time"
	"yorha_http_service"

	"github.com/pterm/pterm"
)

const (
	BotQQID = "1492551697"
)

func startNapCat() {
	for {
		pterm.Info.Println("尝试启动 NapCat")
		err := exec.Command("napcat", "start", BotQQID).Run()
		if err == nil {
			pterm.Success.Println("NapCat 应该启动成功了，10 秒后将启动 Bot")
			time.Sleep(time.Second * 10)
			break
		} else {
			pterm.Error.Println("启动出现问题，将会重试")
			time.Sleep(time.Second * 1)
		}
	}
}

func runner() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	// prepare
	client, err := bot.NewClient("ws://127.0.0.1:8080")
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	pterm.Success.Printf("%#v\n", client.ConnAns)
	// connect to the go-cqhttp server
	closeDown := sync.Mutex{}
	go func() {
		client.ReadPacketAndProcess(&closeDown)
		waitGroup.Done()
	}()
	go yorha_http_service.NewRouter(client).GinEngine.Run(":2018")
	// process messages
	waitGroup.Wait()
	// set wait groups
}

func main() {
	for {
		startNapCat()
		runner()
	}
}
