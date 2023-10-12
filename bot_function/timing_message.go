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

func TimingMessage(
	conn *websocket.Conn,
	resources *RequestCenter.Resources,
	triggeTime time.Time,
	message string,
	userId int64,
) error {
	if time.Now().After(triggeTime) {
		triggeTime = triggeTime.Add(time.Hour * 24)
	}
	// adjust the time
	pterm.Info.Printf(
		"The plug-in routine is installed: Will trigger at %s\n",
		triggeTime.Format("2006-01-02 15:04:05"),
	)
	// show when to send the message
	time.Sleep(
		time.Until(
			triggeTime,
		),
	)
	// waiting for
	resp, err := resources.SendRequestWithResponce(
		conn,
		RequestCenter.Request{
			Action: APIStruct.SendPrivateMsgAction,
			Params: APIStruct.SendPrivateMsg{
				UserId:     userId,
				Message:    message,
				AutoEscape: false,
			},
			RequestId: fmt.Sprintf("%d", resources.GetNewRequestId()),
		},
	)
	if err != nil {
		return fmt.Errorf("TimingMessage: %v", err)
	}
	// send request with responce
	pterm.Success.Printf("%#v\n", resp)
	// output success information
	return nil
	// return
}

func RepeatTiming(
	conn *websocket.Conn,
	resources *RequestCenter.Resources,
) {
	for {
		triggerTime, _ := time.ParseInLocation(
			"2006-01-02 15:04:05",
			fmt.Sprintf("%s 03:00:00", time.Now().Format("2006-01-02 15:04:05")[0:10]),
			time.Local,
		)
		randMinutes := rand.Intn(90)
		randSeconds := rand.Intn(59)
		triggerTime = triggerTime.Add(time.Minute * time.Duration(randMinutes))
		triggerTime = triggerTime.Add(time.Second * time.Duration(randSeconds))
		// spawn a time
		err := TimingMessage(
			conn,
			resources,
			triggerTime,
			"[Auto Generated | INFO] Now I'm sleeping, have a good night!",
			2956859612,
		)
		if err != nil {
			pterm.Warning.Println(err)
		}
		// plan some routines
		time.Sleep(time.Minute * 90)
		// avoid the routine happen again on the same day
	}
}
