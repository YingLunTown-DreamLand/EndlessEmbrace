package BotFunction

import (
	APIStruct "EndlessEmbrace/api_struct"
	RequestCenter "EndlessEmbrace/request_center"
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pterm/pterm"
)

func GenerateNewTime(layer string, randomSecondsRange int) time.Time {
	ans, _ := time.ParseInLocation(
		"2006-01-02 15:04:05",
		fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05")[0:10], layer),
		time.Local,
	)
	ans = ans.Add(time.Second * time.Duration(rand.Intn(randomSecondsRange)))
	return ans
}

func SendPrivateMsg(
	conn *websocket.Conn,
	resources *RequestCenter.Resources,
	message string,
	userId []int64,
) error {
	for _, value := range userId {
		resp, err := resources.SendRequestWithResponce(
			conn,
			RequestCenter.Request{
				Action: APIStruct.SendPrivateMsgAction,
				Params: APIStruct.SendPrivateMsg{
					UserId:     value,
					Message:    message,
					AutoEscape: false,
				},
				RequestId: fmt.Sprintf("%d", resources.GetNewRequestId()),
			},
		)
		if err != nil {
			return fmt.Errorf("SendPrivateMsg: %v", err)
		}
		// send request with responce
		pterm.Success.Printf("%#v\n", resp)
		// output success information
	}
	// send messages in bulk
	return nil
	// return
}

func TimingMessage(
	conn *websocket.Conn,
	resources *RequestCenter.Resources,
	triggeTime time.Time,
	message string,
	userId []int64,
) error {
	for {
		err := SendPrivateMsg(conn, resources, "S2CHeartBeat", userId)
		if err != nil {
			return fmt.Errorf("TimingMessage: %v", err)
		}
		// send heart beat message
		nextBlock := time.Now().Add(time.Minute * (40 + time.Duration(rand.Intn(20))))
		// get sleep time
		if triggeTime.After(nextBlock) {
			pterm.Info.Printf(
				"The next heartbeat is at %s\n",
				nextBlock.Format("2006-01-02 15:04:05"),
			)
			time.Sleep(time.Until(nextBlock))
		} else {
			pterm.Info.Printf("The routine will be executed in the next block\n")
			time.Sleep(time.Until(triggeTime))
			break
		}
		// waiting for
	}
	// send heart beat message and waiting for
	err := SendPrivateMsg(conn, resources, message, userId)
	if err != nil {
		return fmt.Errorf("TimingMessage: %v", err)
	}
	// send message
	return nil
	// return
}

func RepeatTiming(
	conn *websocket.Conn,
	resources *RequestCenter.Resources,
) {
	for {
		triggerTime := GenerateNewTime("03:45:00", 2400)
		heartBeatStartTime := GenerateNewTime("21:25:00", 900)
		// spawn a time
		if time.Now().After(triggerTime) {
			triggerTime = triggerTime.Add(time.Hour * 24)
		}
		if heartBeatStartTime.After(triggerTime) {
			heartBeatStartTime = heartBeatStartTime.Add(-time.Hour * 24)
		}
		// adjust the time
		pterm.Info.Printf(
			"The plug-in routine is installed:\n1. The first heartbeat will occur at %s\n2. The routine will trigger at %s\n",
			heartBeatStartTime.Format("2006-01-02 15:04:05"),
			triggerTime.Format("2006-01-02 15:04:05"),
		)
		// show when to send the message
		time.Sleep(time.Until(heartBeatStartTime))
		// wait for the heartbeat to start
		err := TimingMessage(
			conn,
			resources,
			triggerTime,
			"[Auto Generated | INFO] Now I'm sleeping, have a good night!",
			[]int64{1279923655, 2503170967},
		)
		if err != nil {
			pterm.Warning.Println(err)
		}
		// plan some routines
		time.Sleep(time.Minute * 40)
		// avoid the routine happen again on the same day
	}
}
